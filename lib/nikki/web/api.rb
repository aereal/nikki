require 'graphql'
require 'json-schema'
require 'sequel'
require 'set'
require 'sinatra/base'
require 'swagger/blocks'

require 'nikki/infra/database'
require 'nikki/model/article'
require 'nikki/service/articles'
require 'nikki/service/user'

module Nikki
  module Web
    class Api < ::Sinatra::Base
      require 'nikki/web/middleware/logger'

      include Swagger::Blocks

      use Nikki::Web::Middleware::Logger

      enable :logging

      configure :development do
        require 'sinatra/reloader'
        register ::Sinatra::Reloader
        root = File.expand_path(File.join(settings.root, '../../..'))
        also_reload "#{root}/lib/**/*.rb"
      end

      helpers do
        def validated_api_request(method, path, &block)
          definitions = self.class.api_schema[:definitions]
          body_schema = self.class.api_schema[:paths][path][method][:parameters].first[:schema] || {} # TODO
          req_schema = body_schema.merge(definitions: definitions)
          parsed_body = JSON.load(request.body)
          errors = JSON::Validator.fully_validate(req_schema, parsed_body)
          if errors.empty?
            block.call(parsed_body)
          else
            content_type 'application/json'
            halt(422, JSON.generate(errors: errors.map(&:to_s)))
          end
        end
      end

      class LazySearchUser
        def initialize(query_ctx, user_id)
          @user_id = user_id
          @db_connection = query_ctx[:db_connection]
          @state = query_ctx[:lazy_search_user] ||= {
            pending_ids: Set.new,
            loaded_models: {},
          }
          @state[:pending_ids] << user_id
        end

        def user
          loaded = @state[:loaded_models][@user_id]
          if !loaded
            pending_ids = @state[:pending_ids].to_a
            users = @db_connection[:users].where(id: pending_ids).to_a
            users.each do |row|
              @state[:loaded_models][row[:id]] = Nikki::Model::User.new(**row)
            end
            @state[:pending_ids].clear
            loaded = @state[:loaded_models][@user_id]
          end
          loaded
        end
      end

      UserType = GraphQL::ObjectType.define do
        name 'User'
        description 'blog user'
        field :id, types.ID
        field :name, types.String
      end

      ArticleType = GraphQL::ObjectType.define do
        name 'Article'
        description 'blog post'
        field :id, types.ID
        field :title, !types.String
        field :body, !types.String do
          resolve ->(obj, args, ctx) do
            obj.html_body
          end
        end
        field :created_at, types.String
        field :updated_at, types.String
        field :author, UserType do
          resolve ->(obj, args, ctx) do
            LazySearchUser.new(ctx, obj.author_id)
          end
        end
      end

      ArticleInputType = GraphQL::InputObjectType.define do
        name 'ArticleInputType'
        description 'properties for creating an article'

        argument :title, !types.String do
          description 'Title of the article'
        end
        argument :body, !types.String do
          description 'Body of the article'
        end
      end

      ArticleUpdateInputType = GraphQL::InputObjectType.define do
        name 'ArticleUpdateInputType'
        description 'properties for updating an article'

        argument :title, types.String do
          description 'Title of the article'
        end
        argument :body, types.String do
          description 'Body of the article'
        end
      end

      QueryType = GraphQL::ObjectType.define do
        name 'Query'
        description 'root query'

        field :articles do
          type types[ArticleType]
          description 'search articles'
          argument :limit, types.Int
          resolve ->(obj, args, ctx) {
            db = Nikki::Infra::Database.connection
            rows = db[:articles].limit(args['limit'])
            rows.map {|r| Nikki::Model::Article.new(**r) }
          }
        end
      end

      MutationType = GraphQL::ObjectType.define do
        name 'Mutation'

        field :postArticle, ArticleType do
          description 'post new article'
          argument :article, ArticleInputType
          resolve ->(t, args, ctx) do
            if visitor = ctx[:visitor]
              Nikki::Service::Articles.post(
                db: ctx[:db_connection],
                title: args[:article][:title],
                body: args[:article][:body],
                author: visitor,
              )
            else
              GraphQL::ExecutionError.new("Authentication required")
            end
          end
        end

        field :updateArticle, ArticleType do
          description 'update the article'
          argument :articleId, !types.ID
          argument :article, ArticleUpdateInputType
          resolve ->(t, args, ctx) do
            if visitor = ctx[:visitor]
              if article = Nikki::Service::Articles.find(db: ctx[:db_connection], article_id: args[:articleId])
                article.title = args[:article][:title] if args[:article][:title]
                article.html_body = args[:article][:body] if args[:article][:body]

                Nikki::Service::Articles.update(
                  db: ctx[:db_connection],
                  article: article,
                )
              else
                GraphQL::ExecutionError.new("Article not found")
              end
            else
              GraphQL::ExecutionError.new("Authentication required")
            end
          end
        end
      end

      Schema = GraphQL::Schema.define do
        query QueryType
        mutation MutationType
        lazy_resolve(LazySearchUser, :user)
      end

      options '/graphql' do
        headers['Access-Control-Allow-Methods'] = "HEAD,GET,PUT,POST,DELETE,OPTIONS"
        headers["Access-Control-Allow-Headers"] = "X-Requested-With, X-HTTP-Method-Override, Content-Type, Cache-Control, Accept, Visitor-Key"
        headers 'Access-Control-Allow-Origin' => '*'
        200
      end
      [:get, :post].each do |method|
        send(method, '/graphql') do
          headers 'Access-Control-Allow-Origin' => '*'
          content_type :json

          db = Nikki::Infra::Database.connection

          visitor =
            begin
              if visitor_key = request.get_header('HTTP_VISITOR_KEY')
                Nikki::Service::User.find_by_auth_key(db: db, auth_key: visitor_key)
              else
                nil
              end
            end

          params_hash = %r{\Aapplication/json(?:\b|\z)} === request.media_type ?
            JSON.parse(request.body.read) :
            params
          query = params_hash['query']
          variables = params_hash['variables'] || {}
          context = {
            db_connection: db,
            visitor: visitor,
          }
          result = Schema.execute(query, context: context, variables: variables)
          JSON.generate(result.to_h)
        end
      end

      swagger_root do
        key :swagger, '2.0'
        key :host, 'localhost:9292'
        info do
          key :title, 'Nikki API'
          key :version, '0.0.1'
        end
        security_definition :visitor_key do
          key :type, :apiKey
          key :name, :'visitor-key'
          key :in, :header
        end
      end

      get '/schema' do
        content_type 'application/json'
        headers 'Access-Control-Allow-Origin' => '*'
        JSON.generate(self.class.api_schema)
      end

      swagger_schema :NewArticle do
        property :title do
          key :type, 'string'
        end
        property :body do
          key :type, 'string'
        end
        key :required, [:title, :body]
      end

      swagger_schema :Article do
        property :id do
          key :type, 'integer'
        end
        property :title do
          key :type, 'string'
        end
        property :html_body do
          key :type, 'string'
        end
        key :required, [:id, :title, :html_body]
      end

      swagger_path '/articles' do
        operation :post do
          key :summary, 'Create a new article'
          parameter do
            key :required, true
            key :description, 'new article'
            key :name, :new_article
            key :in, :body
            schema do
              key :'$ref', :NewArticle
            end
          end
          response 200 do
            key :description, 'OK'
            schema do
              key :'$ref', :Article
            end
          end
          response 401 do
            key :description, 'Authentication failed'
          end
          response 422 do
            key :description, 'Invalid parameters'
          end
          security do
            key :visitor_key, []
          end
        end
      end

      options '/articles' do
        headers['Access-Control-Allow-Methods'] = "HEAD,GET,PUT,POST,DELETE,OPTIONS"
        headers["Access-Control-Allow-Headers"] = "X-Requested-With, X-HTTP-Method-Override, Content-Type, Cache-Control, Accept, Visitor-Key"
        headers 'Access-Control-Allow-Origin' => '*'
        200
      end
      post '/articles' do
        content_type 'application/json'
        headers 'Access-Control-Allow-Origin' => '*'

        validated_api_request(:post, :'/articles') do |parsed_body|
          visitor_key = request.get_header('HTTP_VISITOR_KEY')

          unless visitor_key
            halt 401, JSON.generate(errors: ['Unauthorized; no visitor-key header given'])
          end

          db = Nikki::Infra::Database.connection
          authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: visitor_key)

          unless authed_user
            halt 401, JSON.generate(errors: ['Unauthorized; visitor-key is invalid'])
          end

          posted_article = Nikki::Service::Articles.post(db: db, title: parsed_body['title'], body: parsed_body['body'], author: authed_user)
          JSON.generate(posted_article.as_json_hash)
        end
      end

      swagger_path '/articles/{articleId}' do
        operation :get do
          key :summary, 'Get a article'
          parameter do
            key :required, true
            key :description, 'id of the article'
            key :name, :articleId
            key :in, :path
          end
          response 200 do
            key :description, 'OK'
            schema do
              key :'$ref', :Article
            end
          end
          response 401 do
            key :description, 'Authentication failed'
          end
          security do
            key :visitor_key, []
          end
        end

        operation :put do
          key :summary, 'Update an article'
          parameter do
            key :required, true
            key :description, 'id of the article'
            key :name, :articleId
            key :in, :path
          end
          response 200 do
            key :description, 'OK'
            schema do
              key :'$ref', :Article
            end
          end
          response 401 do
            key :description, 'Authentication failed'
          end
          security do
            key :visitor_key, []
          end
        end
      end

      options '/articles/:id' do
        headers['Access-Control-Allow-Methods'] = "HEAD,GET,PUT,POST,DELETE,OPTIONS"
        headers["Access-Control-Allow-Headers"] = "X-Requested-With, X-HTTP-Method-Override, Content-Type, Cache-Control, Accept, Visitor-Key"
        headers 'Access-Control-Allow-Origin' => '*'
        200
      end
      get '/articles/:id' do
        content_type 'application/json'
        headers 'Access-Control-Allow-Origin' => '*'

        validated_api_request(:get, :'/articles/{articleId}') do |parsed_body|
          visitor_key = request.get_header('HTTP_VISITOR_KEY')

          unless visitor_key
            halt 401, JSON.generate(errors: ['Unauthorized; no visitor-key header given'])
          end

          db = Nikki::Infra::Database.connection
          authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: visitor_key)

          unless authed_user
            halt 401, JSON.generate(errors: ['Unauthorized; visitor-key is invalid'])
          end

          article = Nikki::Service::Articles.find(db: db, article_id: params[:id])
          unless article
            halt 404, JSON.generate(errors: ['Article not found'])
          end

          JSON.generate(article.as_json_hash)
        end
      end

      put '/articles/:id' do
        content_type 'application/json'
        headers 'Access-Control-Allow-Origin' => '*'

        validated_api_request(:put, :'/articles/{articleId}') do |parsed_body|
          visitor_key = request.get_header('HTTP_VISITOR_KEY')

          unless visitor_key
            halt 401, JSON.generate(errors: ['Unauthorized; no visitor-key header given'])
          end

          db = Nikki::Infra::Database.connection
          authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: visitor_key)

          unless authed_user
            halt 401, JSON.generate(errors: ['Unauthorized; visitor-key is invalid'])
          end

          article = Nikki::Service::Articles.find(db: db, article_id: params[:id])
          unless article
            halt 404, JSON.generate(errors: ['Article not found'])
          end

          article.title = parsed_body['title']
          article.html_body = parsed_body['body']

          updated_article = Nikki::Service::Articles.update(db: db, article: article)

          JSON.generate(updated_article.as_json_hash)
        end
      end

      def self.api_schema
        @api_schema ||= Swagger::Blocks.build_root_json([self])
      end
    end
  end
end
