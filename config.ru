#!/rackup

lib_path = File.expand_path('./lib', __dir__)
$LOAD_PATH.unshift(lib_path) unless $LOAD_PATH.include?(lib_path)

require 'nikki/env/validator'
require 'nikki/web'

errors = Nikki::Env::Validator.validates
abort "Invalid environment variables; errors: #{errors}" unless errors.empty?

run Nikki::Web
