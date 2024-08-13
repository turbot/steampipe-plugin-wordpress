# Table: wordpress_category

Represents a category in WordPress

## Examples

### List categories by id and name

```sql
select
  id,
  name
from
  wordpress_category
```

### List categories matching `data`

```sql
select
  id,
  name
from
  wordpress_category
where
  name ~* 'data'
```
