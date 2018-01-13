require 'sequel'
require 'logger'

module Nikki
  module Infra
    module Database
      def self.connection
        @connection ||= self.connect!
      end

      def self.connect!
        db = Sequel.connect("postgres://postgres:postgres@#{ENV['DB_HOST']}/nikki", log_connection_info: true)
        logger = Logger.new($stdout)
        logger.progname = 'sequel'
        db.logger = logger
        logger.debug("Connected")
        db
      end

      def self.disconnect!
        if @connection
          @connection.disconnect
        end
      end
    end
  end
end
