require 'sinatra/base'
require 'slim'

require 'nikki/infra/database'
require 'nikki/model/pager'
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

      not_found do
        slim :not_found, locals: { page_title: 'Not found', site_title: 'Nikki' }
      end

      get '/' do
        pager_token = params[:page]
        pager = Nikki::Model::Pager.new_from_token(pager_token || '')

        db = Nikki::Infra::Database.connection
        pager = Nikki::Service::Articles.search(db: db, limit: 10, pager: pager)
        formatted_articles = pager[:articles].map {|a| Nikki::Service::Articles.format_body(article: a) }
        locals = {
          page_title: nil,
          site_title: 'Nikki',
          articles: formatted_articles,
          next_page_token: pager[:next_page_token],
        }
        slim :index, locals: locals
      end

      get '/*' do
        logger.info("path_info = #{request.path_info}")
        db = Nikki::Infra::Database.connection
        if article = Nikki::Service::Articles.find_by_path(db: db, path: request.path_info)
          formatted_article = Nikki::Service::Articles.format_body(article: article)
          locals = {
            article: formatted_article,
            page_title: formatted_article.title,
            site_title: 'Nikki',
          }
          slim :permalink, locals: locals
        else
          halt 404
        end
      end
    end
  end
end
