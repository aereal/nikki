require 'graphql'
require 'set'
require 'sinatra/base'

require 'nikki/infra/database'
require 'nikki/model/article'
require 'nikki/service/articles'
require 'nikki/service/user'

module Nikki
  module Web
    class Api < ::Sinatra::Base
      require 'nikki/web/middleware/logger'

      use Nikki::Web::Middleware::Logger

      enable :logging

      configure :development do
        require 'sinatra/reloader'
        register ::Sinatra::Reloader
        root = File.expand_path(File.join(settings.root, '../../..'))
        also_reload "#{root}/lib/**/*.rb"
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
        field :body, !types.String
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
          argument :limit, !types.Int
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
                article.body = args[:article][:body] if args[:article][:body]

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

      module MutationRejector
        def self.before_query(query)
          if query.mutation?
            raise "Mutation not allowed"
          end
        end

        def self.after_query(query)
          # no-op
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
          schema = Schema
          if request.get?
            schema.instrument(:query, MutationRejector)
          end
          result = schema.execute(query, context: context, variables: variables)
          JSON.generate(result.to_h)
        end
      end
    end
  end
end
