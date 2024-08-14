# Table: wordpress_tag

Represents a category in WordPress

## Examples

### List categories by id and name

```sql
select
  id,
  name
from
  wordpress_tag
```
### Tabulate tags by count

```sql
select
  name, 
  count(*)
from
  wordpress_tag
group by
 name
order by
 count desc
```

### List posts with a tag

```sql
  select
    p.id,
    p.title,
    p.author,
    p.link,
    to_char(p.date, 'YYYY-MM-DD') as day
  from
    wordpress_post p
  where
    '543201891' = any(array(select jsonb_array_elements_text(p.tag)))
```
