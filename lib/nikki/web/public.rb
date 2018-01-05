require 'sinatra/base'
require 'slim'

module Nikki
  module Web
    class Public < ::Sinatra::Base
      require 'nikki/web/middleware/logger'

      use Nikki::Web::Middleware::Logger

      enable :sessions
      enable :logging
      templates_root = File.expand_path(File.join(settings.root, '../../../templates'))
      set :views, File.join(templates_root, 'public')
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
        root = File.expand_path(File.join(settings.root, '../../..'))
        also_reload "#{root}/lib/**/*.rb"
      end

      get '/' do
        slim :index
      end
    end
  end
end
