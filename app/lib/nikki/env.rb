require 'json'
require 'yaml'

module Nikki
  module Env
    def self.schema
      @schema ||=
        begin
          root = File.expand_path('../../etc', __dir__)
          JSON.dump(YAML.load_file(File.join(root, 'env-schema.yml')))
        end
    end
  end
end
