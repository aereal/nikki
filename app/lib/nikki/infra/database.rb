require 'sequel'
require 'logger'

module Nikki
  module Infra
    module Database
      def self.connection
        @connection ||= self.connect!
      end

      def self.connect!
        dsn_url = ENV['DB_DSN_URL']
        db = Sequel.connect(dsn_url, log_connection_info: true)
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
