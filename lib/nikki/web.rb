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

    get '/' do
      visitor, * = [session[:visitor_id]].compact.map {|id| Nikki::Model::User.new(id) }
      slim :index, locals: { visitor: visitor }
    end

    get '/auth/-/logout' do
      session[:visitor_id] = nil
      redirect '/'
    end
  end
end
