module Nikki
  module Model
    class Article
      attr_reader :id, :author_id, :created_at, :updated_at, :path
      attr_accessor :title, :body

      def initialize(id: , title: , body: , author_id: , created_at: , updated_at: , path: )
        @id = id
        @title = title
        @body = body
        @author_id = author_id
        @created_at = created_at
        @updated_at = updated_at
        @path = path
      end

      def as_json_hash
        {
          id: self.id,
          title: self.title,
          body: self.body,
          created_at: self.created_at.to_s,
          updated_at: self.updated_at.to_s,
          path: self.path,
        }
      end

      def with_formatted_body(formatted_body)
        @formatted_body = formatted_body
        self.extend(Formatted)
        self
      end

      module Formatted
        attr_reader :formatted_body
      end
    end
  end
end
