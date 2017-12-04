require 'json'
require 'json-schema'
require 'yaml'

module Nikki
  module EnvValidator
    def self.validates(input = ENV.to_h)
      root = File.expand_path('../../etc', __dir__)
      schema = JSON.dump(YAML.load_file(File.join(root, 'env-schema.yml')))
      JSON::Validator.fully_validate(schema, input, strict: true)
    end
  end
end
