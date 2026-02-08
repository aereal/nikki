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
  article_revisions.title
from
  articles
  inner join article_revisions on article_revisions.article_id = articles.article_id
where
  articles.slug = ?;
