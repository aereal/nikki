require 'sinatra/base'

module Nikki
  class Web < ::Sinatra::Base
    configure :development do
      require 'sinatra/reloader'
      register ::Sinatra::Reloader
    end

    get '/' do
      'hi'
    end
  end
end
