require 'redcarpet'

require 'nikki/model/article'

module Nikki
  module Service
    module Articles
      def self.search(db: , limit: )
        db[:articles].
          reverse_order(:created_at).
          limit(limit).
          map {|row| Nikki::Model::Article.new(**row) }
      end

      def self.find(db: , article_id: )
        row = db[:articles].where(id: article_id).first
        row ? Nikki::Model::Article.new(**row) : nil
      end

      def self.post(db: , title: , body: , author: )
        created_at = updated_at = Time.now
        rows = db[:articles].returning.insert(
          title: title,
          body: body,
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
          body: article.body,
          updated_at: updated_at,
        )
        row = rows.first
        Nikki::Model::Article.new(**row)
      end

      def self.format_body(article: )
        @formatter ||=
          begin
            renderer = Redcarpet::Render::HTML.new(
              no_styles: true,
              safe_links_only: true,
            )
            exts = {
              no_intra_emphasis: true,
              tables: true,
              autolink: true,
              disable_indented_code_blocks: true,
              footnotes: true,
            }
            Redcarpet::Markdown.new(renderer, exts)
          end
        formatted_body = @formatter.render(article.body)
        article.with_formatted_body(formatted_body)
      end
    end
  end
end
