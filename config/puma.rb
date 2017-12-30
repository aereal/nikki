lib_path = File.expand_path('../lib', __dir__)
$LOAD_PATH.unshift(lib_path) unless $LOAD_PATH.include?(lib_path)

require 'nikki/infra/database'

workers 4
threads 0, 6
preload_app!

on_worker_boot do
  tried = 0
  max_tries = 3
  wait_time = 1
  connected = false
  while !connected && tried < max_tries do
    begin
      Nikki::Infra::Database.connect!
      connected = true
    rescue Sequel::DatabaseConnectionError
      sleep wait_time
      wait_time = wait_time * 2
    end
  end
end

on_worker_shutdown do
  Nikki::Infra::Database.disconnect!
end
