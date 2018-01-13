module Nikki
  module Model
    class User
      attr_reader :id, :name, :slug, :auth_key

      def initialize(id: , name: , slug: , auth_key: )
        @id = id
        @name = name
        @slug = slug
        @auth_key = auth_key
      end
    end
  end
end
