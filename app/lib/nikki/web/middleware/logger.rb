require 'rack/common_logger'

module Nikki::Web::Middleware
  class Logger < ::Rack::CommonLogger
    FIELDS = %w(requested_at status path method size host vhost taken content_type origin ua)
    FORMAT = FIELDS.map {|field| "#{field}:%{#{field}}"}.join("\t") + "\n"

    def log(env, status, header, began_at)
      now = Time.now
      length = extract_content_length(header)
      logger = @logger || env['rack.errors']
      logger.write FORMAT % {
        requested_at: now.strftime('%Y-%m-%dT%H:%M:%S%z'),
        status: status.to_s[0..3],
        path: env['PATH_INFO'],
        method: env['REQUEST_METHOD'],
        size: length || '-',
        host: env['HTTP_X_FORWARDED_FOR'] || env['REMOTE_ADDR'] || '-',
        vhost: env['HTTP_X_FORWARDED_HOST'] || env['HTTP_HOST'] || env['SERVER_NAME'] || '-',
        taken: ('%0.6f' % [now - began_at]) || '-',
        content_type: env['HTTP_CONTENT_TYPE'] || '-',
        origin: env['HTTP_ORIGIN'] || '-',
        ua: env['HTTP_USER_AGENT'] || '-',
      }
    end
  end
end
