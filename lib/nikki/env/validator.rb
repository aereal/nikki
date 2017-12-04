require 'json'
require 'json-schema'
require 'yaml'

require 'nikki/env'

module Nikki
  module Env
    module Validator
      def self.validates(input = ENV.to_h)
        schema = Nikki::Env.schema
        JSON::Validator.fully_validate(schema, input, strict: true)
      end
    end
  end
end
