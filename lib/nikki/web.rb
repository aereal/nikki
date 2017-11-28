require 'omniauth'
require 'omniauth-google-oauth2'
require 'sinatra/base'
require 'slim'

module Nikki
  module Model
    class User
      attr_reader :id

      def initialize(id)
        @id = id
      end
    end
  end

  class Web < ::Sinatra::Base
    enable :sessions
    enable :logging
    set :views, File.expand_path(File.join(settings.root, '../../templates'))

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
      visitor, * = [session[:visitor_id]].compact.map {|id| Nikki::Model::User.new(id) }
      slim :index, locals: { visitor: visitor }
    end

    get '/auth/:provider/callback' do
      auth = env['omniauth.auth']
      session[:visitor_id] = "provider:google:id:#{auth.uid}"
      redirect '/'
    end

    post '/auth/:provider/callback' do
      auth = env['omniauth.auth']
      session[:visitor_id] = "provider:google:id:#{auth.uid}"
      redirect '/'
    end

    get '/auth/-/logout' do
      session[:visitor_id] = nil
      redirect '/'
    end
  end
end
