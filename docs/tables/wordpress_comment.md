# Table: wordpress_comment

Represents a comment on a WordPress post.

## Examples

### Count comments by author

```sql
select 
  author_name, 
  count(*) 
from 
  wordpress_comment
group by
  author_name
order by count desc
```

### List recent comments

```sql
with posts as (
  select
    id,
    title,
    link
  from
    wordpress_post
  where
    date > now() - interval '1 month'  
)
select
  c.date,
  c.author_name,
  p.title,
  p.link
from
  wordpress_comment c
join
  posts p
on
  c.post = p.id
order by
  c.date desc
limit
  10
```  