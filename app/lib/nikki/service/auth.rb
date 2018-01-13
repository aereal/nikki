require 'jwt'
require 'openssl'

module Nikki
  module Service
    module Auth
      def self.try_authenticate(id_token: , pub_keys_json_path: )
        encoded_header, * = id_token.split('.')
        header = JSON.parse(Base64.decode64(encoded_header))

        @pub_keys ||= JSON.load(open(pub_keys_json_path))
        certs = @pub_keys.map {|(kid, cert)| [kid, OpenSSL::X509::Certificate.new(cert)] }.to_h
        cert = certs[header['kid']]
        return unless cert

        decoded_token = JWT.decode(id_token, cert.public_key, !!cert.public_key, algorithm: 'RS256')
      end
    end
  end
end
