-- name: ReviseArticle :exec
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
