module Nikki
  module Model
    class User
      attr_reader :id, :name, :slug

      def initialize(id: , name: , slug: )
        @id = id
        @name = name
        @slug = slug
      end
    end
  end
end
