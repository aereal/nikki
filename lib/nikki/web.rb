require 'digest/sha2'
require 'json-schema'
require 'omniauth'
require 'omniauth-google-oauth2'
require 'rack/common_logger'
require 'sequel'
require 'sinatra/base'
require 'slim'
require 'swagger/blocks'

require 'nikki/service/articles'
require 'nikki/service/user'

module Nikki
  class Web < ::Sinatra::Base
    include Swagger::Blocks

    module Middleware
      class Logger < ::Rack::CommonLogger
        FIELDS = %w(requested_at status path method size host vhost taken content_type origin ua)
        FORMAT = FIELDS.map {|field| "#{field}:%{#{field}}"}.join("\t") + "\n"

        def log(env, status, header, began_at)
          now = Time.now
          length = extract_content_length(header)
          logger = @logger || env['rack.errors']
          logger.write FORMAT % {
            requested_at: now.strftime('%Y-%m-%dT%H:%M:%S%z'),
            status: status.to_s[0..3],
            path: env['PATH_INFO'],
            method: env['REQUEST_METHOD'],
            size: length || '-',
            host: env['HTTP_X_FORWARDED_FOR'] || env['REMOTE_ADDR'] || '-',
            vhost: env['HTTP_X_FORWARDED_HOST'] || env['HTTP_HOST'] || env['SERVER_NAME'] || '-',
            taken: ('%0.6f' % [now - began_at]) || '-',
            content_type: env['HTTP_CONTENT_TYPE'] || '-',
            origin: env['HTTP_ORIGIN'] || '-',
            ua: env['HTTP_USER_AGENT'] || '-',
          }
        end
      end
    end

    use Nikki::Web::Middleware::Logger

    enable :sessions
    enable :logging
    set :views, File.expand_path(File.join(settings.root, '../../templates'))
    set :public_folder, File.expand_path(File.join(settings.root, '../../assets/'))

    configure do
      Slim::Engine.set_options(
        tabsize: 2,
        enable_engines: [],
        format: :html,
        sort_attrs: false,
      )
    end

    configure :development do
      require 'sinatra/reloader'
      register ::Sinatra::Reloader
    end

    use OmniAuth::Builder do
      provider :google_oauth2, ENV['GOOGLE_OAUTH_CLIENT_ID'], ENV['GOOGLE_OAUTH_CLIENT_SECRET']
    end

    helpers do
      def validated_api_request(method, path, &block)
        definitions = self.class.api_schema[:definitions]
        body_schema = self.class.api_schema[:paths][path][method][:parameters].first[:schema] # TODO
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

    get '/' do
      db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
      authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: session[:auth_key])
      initial_props = {
        authedUser: authed_user.nil? ? nil : { name: authed_user.name, slug: authed_user.slug, authKey: authed_user.auth_key, },
      }
      slim :index, locals: { initial_props: JSON.generate(initial_props) }
    end

    get '/auth/:provider/callback' do
      auth = env['omniauth.auth']
      db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
      auth_key = Digest::SHA256.hexdigest("provider:google:uid:#{auth.uid}")
      user = Nikki::Service::User.find_or_register_by(db: db, name: auth[:info][:name], slug: auth[:info][:email], auth_key: auth_key)
      session[:auth_key] = user.auth_key
      redirect '/'
    end

    post '/auth/:provider/callback' do
      auth = env['omniauth.auth']
      db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
      auth_key = Digest::SHA256.hexdigest("provider:google:uid:#{auth.uid}")
      user = Nikki::Service::User.find_or_register_by(db: db, name: auth[:info][:name], slug: auth[:info][:email], auth_key: auth_key)
      session[:auth_key] = user.auth_key
      redirect '/'
    end

    get '/auth/-/logout' do
      session[:visitor_id] = nil
      redirect '/'
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
      headers['allow'] = "HEAD,GET,PUT,POST,DELETE,OPTIONS"
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

    def self.api_schema
      @api_schema ||= Swagger::Blocks.build_root_json([self])
    end
  end
end
