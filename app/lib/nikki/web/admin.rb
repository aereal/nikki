require 'sequel'
require 'sinatra/base'
require 'slim'

require 'nikki/infra/database'
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
        root = File.expand_path(File.join(settings.root, '../../..'))
        also_reload "#{root}/lib/**/*.rb"
      end

      get '/graphql' do
        slim :graphiql, locals: { initial_props: '{}' }
      end

      get '/' do
        slim :index, locals: { initial_props: '{}' }
      end

      get '/articles/:id' do
        db = Nikki::Infra::Database.connection
        article = Nikki::Service::Articles.find(db: db, article_id: params[:id])
        initial_props = {
          article: article ? article.as_json_hash : nil,
        }
        slim :index, locals: { initial_props: JSON.generate(initial_props) }
      end
    end
  end
end
