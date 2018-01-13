require 'nikki/model/user'

module Nikki
  module Service
    module User
      def self.find_or_register_by(db: , name: , slug: , auth_key: )
        find_by_auth_key(db: db, auth_key: auth_key) || register(db: db, name: name, slug: slug, auth_key: auth_key)
      end

      def self.find_by_auth_key(db: , auth_key: )
        row = db[:users].where(auth_key: auth_key).first
        if row
          Nikki::Model::User.new(**row)
        else
          nil
        end
      end

      def self.find_by_slug(db: , slug: )
        row = db[:users].where(slug: slug).first
        if row
          Nikki::Model::User.new(**row)
        else
          nil
        end
      end

      def self.register(db: , name: , slug: , auth_key: )
        rows = db[:users].returning.insert(
          name: name,
          slug: slug,
          auth_key: auth_key,
        )
        row = rows.first
        Nikki::Model::User.new(**row)
      end
    end
  end
end
