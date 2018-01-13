require 'base64'
require 'redcarpet'

require 'nikki/model/article'
require 'nikki/model/pager'

module Nikki
  module Service
    module Articles
      def self.search(db: , limit: , pager: )
        query = db[:articles].
          reverse_order(:created_at).
          limit(limit + 1)
        if pager && pager.from
          query = query.where { created_at <= pager.from }
        end
        articles = query.map {|row| Nikki::Model::Article.new(**row) }

        next_page_token = nil
        if next_article = articles[limit]
          pager = Nikki::Model::Pager.new(from: next_article.created_at)
          next_page_token = pager.to_s
        end

        pager = {
          articles: articles[0, limit],
          next_page_token: next_page_token,
        }
        pager
      end

      def self.find(db: , article_id: )
        row = db[:articles].where(id: article_id).first
        row ? Nikki::Model::Article.new(**row) : nil
      end

      def self.find_by_path(db: , path: )
        row = db[:articles].where(path: path).first
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
          path: created_at.strftime('/%Y/%m/%d/%H%M%S'),
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
