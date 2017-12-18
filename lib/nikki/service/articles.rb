require 'nikki/model/article'

module Nikki
  module Service
    module Articles
      def self.find(db: , article_id: )
        row = db[:articles].where(id: article_id).first
        row ? Nikki::Model::Article.new(**row) : nil
      end

      def self.search_by_author(db: , author: )
        rows = db[:articles].where(author_id: author.id)
        rows.map {|r| Nikki::Model::Article.new(**r) }
      end

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

      def self.update(db: , article: )
        updated_at = Time.now
        rows = db[:articles].returning.where(id: article.id).update(
          title: article.title,
          html_body: article.html_body,
          updated_at: updated_at,
        )
        row = rows.first
        Nikki::Model::Article.new(**row)
      end
    end
  end
end
