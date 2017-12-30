lib_path = File.expand_path('../lib', __dir__)
$LOAD_PATH.unshift(lib_path) unless $LOAD_PATH.include?(lib_path)

require 'nikki/infra/database'

workers 4
threads 0, 6
preload_app!

on_worker_boot do
  Nikki::Infra::Database.connect!
end

on_worker_shutdown do
  Nikki::Infra::Database.disconnect!
end
