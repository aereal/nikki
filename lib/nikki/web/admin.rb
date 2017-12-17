require 'digest/sha2'
require 'json-schema'
require 'omniauth'
require 'omniauth-google-oauth2'
require 'sequel'
require 'sinatra/base'
require 'slim'

require 'nikki/service/articles'
require 'nikki/service/user'

module Nikki
  module Web
    class Admin < ::Sinatra::Base
      require 'nikki/web/middleware/logger'

      use Nikki::Web::Middleware::Logger

      enable :sessions
      enable :logging
      set :views, File.expand_path(File.join(settings.root, '../../../templates'))
      set :public_folder, File.expand_path(File.join(settings.root, '../../../assets/'))

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
    end
  end
end
