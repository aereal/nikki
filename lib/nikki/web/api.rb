require 'json-schema'
require 'sequel'
require 'sinatra/base'
require 'swagger/blocks'

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

          db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
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

          db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
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

      def self.api_schema
        @api_schema ||= Swagger::Blocks.build_root_json([self])
      end
    end
  end
end