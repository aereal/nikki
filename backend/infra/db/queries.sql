-- name: ImportCategories :exec
insert into
  categories (category_id, name)
values
  (?, ?) on conflict (name) do nothing;

-- name: FindCategoriesByNames :many
select
  *
from
  categories
where
  name in (sqlc.slice ('names'));

-- name: CreateArticles :exec
insert into
  articles (article_id, slug)
values
  (?, ?);

-- name: CreateArticleRevisions :exec
insert into
  article_revisions (
    article_revision_id,
    article_id,
    title,
    body,
    authored_at
  )
values
  (?, ?, ?, ?, ?);

-- name: CreateArticlePublications :exec
insert into
  article_publications (article_id, article_revision_id, published_at)
values
  (?, ?, ?);

-- name: MapArticleCategory :exec
insert into
  article_category_mappings (article_id, category_id)
values
  (?, ?);

-- name: FindArticleBySlug :one
select
  articles.article_id,
  articles.slug,
  article_revisions.body,
  article_revisions.title,
  article_publications.published_at
from
  article_revisions
  inner join articles on articles.article_id = article_revisions.article_id
  and articles.slug = ?
  inner join article_publications on article_publications.article_id = articles.article_id
order by
  article_revisions.authored_at desc
limit
  1;

-- name: FindLatestArticles :many
select
  articles.article_id,
  articles.slug,
  article_revisions.body,
  article_revisions.title,
  article_publications.published_at
from
  articles
  inner join article_revisions on article_revisions.article_id = articles.article_id
  and article_revisions.authored_at = (
    select
      max(ar.authored_at)
    from
      article_revisions ar
    where
      ar.article_id = articles.article_id
  )
  inner join article_publications on article_publications.article_id = articles.article_id
order by
  article_publications.published_at desc
limit
  sqlc.arg ('limit');

-- name: FindEarlyArticles :many
select
  articles.article_id,
  articles.slug,
  article_revisions.body,
  article_revisions.title,
  article_publications.published_at
from
  articles
  inner join article_revisions on article_revisions.article_id = articles.article_id
  and article_revisions.authored_at = (
    select
      max(ar.authored_at)
    from
      article_revisions ar
    where
      ar.article_id = articles.article_id
  )
  inner join article_publications on article_publications.article_id = articles.article_id
order by
  article_publications.published_at asc
limit
  sqlc.arg ('limit');

-- name: FindLatestArticlesAfter :many
select
  articles.article_id,
  articles.slug,
  article_revisions.body,
  article_revisions.title,
  article_publications.published_at
from
  articles
  inner join article_revisions on article_revisions.article_id = articles.article_id
  and article_revisions.authored_at = (
    select
      max(ar.authored_at)
    from
      article_revisions ar
    where
      ar.article_id = articles.article_id
  )
  inner join article_publications on article_publications.article_id = articles.article_id
where
  article_publications.published_at > sqlc.arg ('after')
order by
  article_publications.published_at desc
limit
  sqlc.arg ('limit');

-- name: FindEarlyArticlesBefore :many
select
  articles.article_id,
  articles.slug,
  article_revisions.body,
  article_revisions.title,
  article_publications.published_at
from
  articles
  inner join article_revisions on article_revisions.article_id = articles.article_id
  and article_revisions.authored_at = (
    select
      max(ar.authored_at)
    from
      article_revisions ar
    where
      ar.article_id = articles.article_id
  )
  inner join article_publications on article_publications.article_id = articles.article_id
where
  article_publications.published_at < sqlc.arg ('before')
order by
  article_publications.published_at asc
limit
  sqlc.arg ('limit');
