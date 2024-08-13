# Table: wordpress_post

Represents a post in WordPress, including its content, metadata, and associated information.

## Examples

### List recent posts with their titles and publication dates

```sql
select
  id,
  title,
  to_char(date, 'YYYY-MM-DD') as publication_date
from
  wordpress_post
where
  date > now() - interval '1 week'
```

### Find recent posts in a category

```sql
with posts as (
  select
    title,
    link,
    jsonb_array_elements(category)::int as category_id
  from
    wordpress_post
  where date > now() - interval '1 week'
)
select
  *
from
  posts
where
  category_id = 12945
```


### List posts for an author

```sql
select
  id,
  title,
  to_char(date, 'YYYY-MM-DD') as publication_date,
  category
from
  wordpress_post
where 
  author = 2488
```

### List recent posts in a category

```sql
  with posts as (
    select
      p.id,
      p.title,
      p.author,
      p.link,
      jsonb_array_elements_text(p.category) as category_id,
      to_char(p.date, 'YYYY-MM-DD') as day
    from
      wordpress_post p
    where
      '12406' = any(array(select jsonb_array_elements_text(p.category)))
      and p.date > now() - interval '1 month'
  )
  select
    p.title,
    p.author,
    p.day,
    string_agg(distinct p.category_id, ', ' order by p.category_id) as category_ids,
    p.link
  from
    posts p
  group by
    p.id, p.title, p.author, p.day, p.link
  order by
    p.day desc, p.author;
```

