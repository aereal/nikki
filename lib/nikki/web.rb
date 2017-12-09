require 'digest/sha2'
require 'omniauth'
require 'omniauth-google-oauth2'
require 'sequel'
require 'sinatra/base'
require 'slim'
require 'swagger/blocks'

require 'nikki/service/user'

module Nikki
  class Web < ::Sinatra::Base
    include Swagger::Blocks

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

    swagger_root do
      key :swagger, '2.0'
      key :host, 'localhost:9292'
      info do
        key :title, 'Nikki API'
        key :version, '0.0.1'
      end
    end

    get '/schema' do
      content_type 'application/json'
      schema_json = Swagger::Blocks.build_root_json([self.class])
      headers 'Access-Allow-Allow-Origin' => '*'
      JSON.generate(schema_json)
    end

    get '/' do
      db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki")
      authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: session[:auth_key])
      initial_props = {
        authedUser: authed_user.nil? ? nil : { name: authed_user.name, slug: authed_user.slug },
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
  end
end
