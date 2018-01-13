require 'sinatra/base'
require 'slim'

require 'nikki/infra/database'
require 'nikki/service/articles'

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
        db = Nikki::Infra::Database.connection
        articles = Nikki::Service::Articles.search(db: db, limit: 10)
        formatted_articles = articles.map {|a| Nikki::Service::Articles.format_body(article: a) }
        locals = {
          page_title: 'Nikki',
          articles: formatted_articles,
        }
        slim :index, locals: locals
      end
    end
  end
end
