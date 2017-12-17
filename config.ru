#!/rackup

lib_path = File.expand_path('./lib', __dir__)
$LOAD_PATH.unshift(lib_path) unless $LOAD_PATH.include?(lib_path)

require 'nikki/env/validator'

errors = Nikki::Env::Validator.validates
abort "Invalid environment variables; errors: #{errors}" unless errors.empty?

require 'nikki/web/admin'
require 'nikki/web/api'

map 'https://admin.nikki.dev/' do
  run Nikki::Web::Admin
end

map 'https://api.nikki.dev/' do
  run Nikki::Web::Api
end
