require 'faraday'
require 'json'
require 'jwt'
require 'openssl'

module Nikki
  module Service
    module Auth
      GOOGLE_OAUTH_PUBLIC_KEYS_URL = 'https://www.googleapis.com/oauth2/v1/certs'.freeze

      def self.try_authenticate(id_token: , pub_keys_json_path: )
        encoded_header, * = id_token.split('.')
        header = JSON.parse(Base64.decode64(encoded_header))

        cert = get_cert_for(kid: header['kid'])
        return unless cert

        JWT.decode(id_token, cert.public_key, !!cert.public_key, algorithm: 'RS256')
      end

      # => Maybe[$cert: String]
      def self.get_cert_for(kid: )
        certs = self.fetch_certs_with_cache
        certs[kid]
      end

      def self.fetch_certs_with_cache
        @certs ||=
          begin
            res = Faraday.get(GOOGLE_OAUTH_PUBLIC_KEYS_URL)
            if res.success?
              parsed =
                begin
                  JSON.parse(res.body)
                rescue JSON::ParserError
                  nil
                end
              if parsed
                parsed.
                  map {|(kid, cert_str)| [kid, OpenSSL::X509::Certificate.new(cert_str)] }.
                  to_h
              end
            else
              nil
            end
          end
        @certs
      end
    end
  end
end
