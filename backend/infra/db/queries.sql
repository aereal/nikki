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
