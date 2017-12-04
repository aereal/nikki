#!/rackup

require_relative './lib/nikki/env_validator'
require_relative './lib/nikki/web'

errors = Nikki::EnvValidator.validates
abort "Invalid environment variables; errors: #{errors}" unless errors.empty?

run Nikki::Web
