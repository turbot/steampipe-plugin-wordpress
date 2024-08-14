# Table: wordpress_author

Represents a category in WordPress

## Examples

### List authors matching "Richard"

```sql
select
  *
from
  wordpress_author
where
  name ~* 'richard'
```
### Tabulate authors by domain

```sql
with domains as (
  select (regexp_match(lower(url), '^(?:https?://)?(?:www\.)?([^/:]+)'))[1] as domain
  from wordpress_author
  where url is not null
)
select 
  domain,
  count(*) as count
from domains
group by domain
order by count desc;
```

