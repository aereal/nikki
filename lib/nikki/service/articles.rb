require 'nikki/model/article'

module Nikki
  module Service
    module Articles
      def self.post(db: , title: , body: , author: )
        html_body = body # TODO
        created_at = updated_at = Time.now
        rows = db[:articles].returning.insert(
          title: title,
          html_body: html_body,
          author_id: author.id,
          created_at: created_at,
          updated_at: updated_at,
        )
        row = rows.first
        Nikki::Model::Article.new(**row)
      end
    end
  end
end
