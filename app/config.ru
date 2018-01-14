#!/rackup

lib_path = File.expand_path('./lib', __dir__)
$LOAD_PATH.unshift(lib_path) unless $LOAD_PATH.include?(lib_path)

require 'nikki/env/validator'

errors = Nikki::Env::Validator.validates
abort "Invalid environment variables; errors: #{errors}" unless errors.empty?

require 'nikki/web/admin'
require 'nikki/web/api'
require 'nikki/web/public'

public_origin = ENV['PUBLIC_ORIGIN']
admin_origin = ENV['ADMIN_ORIGIN']
api_origin = ENV['API_ORIGIN']

map "#{public_origin}/" do
  run Nikki::Web::Public
end

map "#{admin_origin}/" do
  run Nikki::Web::Admin
end

map "#{api_origin}/" do
  run Nikki::Web::Api
end
