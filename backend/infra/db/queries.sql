-- name: ImportCategories :exec
insert into
  categories (category_id, name)
values
  (?, ?) on conflict (name) do nothing;
