create table articles (
  article_id text primary key,
  slug text not null unique
);

create table article_revisions (
  article_revision_id text primary key,
  article_id text not null references articles (article_id),
  title text not null,
  body text not null,
  authored_at text not null
);

create table article_publications (
  article_id text not null references articles primary key,
  article_revision_id text not null references article_revisions (article_revision_id),
  published_at text not null
);

create table categories (
  category_id text primary key,
  name text not null unique
);

create table article_category_mappings (
  article_id text not null references articles (article_id),
  category_id text not null references categories (category_id),
  primary key (article_id, category_id)
);
