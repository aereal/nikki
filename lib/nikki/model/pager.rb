require 'base64'
require 'time'

module Nikki
  module Model
    class Pager
      FIELDS_DELIMITER = ':'.freeze
      PAIRS_DELIMITER = ';'.freeze

      attr_reader :from

      def self.new_from_token(token)
        raw_token =
          begin
            Base64.urlsafe_decode64(token)
          rescue ArgumentError
            ''
          end
        params = Hash[
          raw_token.split(PAIRS_DELIMITER).map {|pair| pair.split(FIELDS_DELIMITER, 2) }
        ]
        new(from: params['from'] ? Time.parse(params['from']) : nil)
      end

      def initialize(from: )
        @from = from
      end

      def to_s
        params = {
          from: self.from.iso8601(6),
        }
        raw_token = params.
          map {|pair| pair.join(FIELDS_DELIMITER) }.
          join(PAIRS_DELIMITER)
        Base64.urlsafe_encode64(raw_token)
      end
    end
  end
end
