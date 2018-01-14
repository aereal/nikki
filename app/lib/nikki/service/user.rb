require 'nikki/model/user'

module Nikki
  module Service
    module User
      def self.find_or_register_by(db: , name: , slug: )
        find_by_slug(db: db, slug: slug) || register(db: db, name: name, slug: slug)
      end

      def self.find_by_slug(db: , slug: )
        row = db[:users].where(slug: slug).first
        if row
          Nikki::Model::User.new(**row)
        else
          nil
        end
      end

      def self.register(db: , name: , slug: )
        rows = db[:users].returning.insert(
          name: name,
          slug: slug,
        )
        row = rows.first
        Nikki::Model::User.new(**row)
      end
    end
  end
end
