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
        db = Nikki::Infra::Database.connection
        authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: session[:auth_key])
        initial_props = {
          authedUser: authed_user.nil? ? nil : { name: authed_user.name, slug: authed_user.slug, authKey: authed_user.auth_key, },
        }
        slim :graphiql, locals: { initial_props: JSON.generate(initial_props) }
      end

      get '/' do
        db = Nikki::Infra::Database.connection
        authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: session[:auth_key])
        initial_props = {
          authedUser: authed_user.nil? ? nil : { name: authed_user.name, slug: authed_user.slug, authKey: authed_user.auth_key, },
        }
        slim :index, locals: { initial_props: JSON.generate(initial_props) }
      end

      get '/articles/:id' do
        db = Nikki::Infra::Database.connection
        authed_user = Nikki::Service::User.find_by_auth_key(db: db, auth_key: session[:auth_key])
        article = Nikki::Service::Articles.find(db: db, article_id: params[:id])
        initial_props = {
          authedUser: authed_user.nil? ? nil : { name: authed_user.name, slug: authed_user.slug, authKey: authed_user.auth_key, },
          article: article ? article.as_json_hash : nil,
        }
        slim :index, locals: { initial_props: JSON.generate(initial_props) }
      end
    end
  end
end
