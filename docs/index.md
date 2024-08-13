---
organization: Turbot
category: ["media"]
icon_url: "/images/plugins/turbot/mastodon.svg"
brand_color: "#21759b"
display_name: WordPress
name: wordpress
description: Use SQL to instantly query WordPress posts, authors, and categories.
og_description: Query WordPressMastodon with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/wordpress-social-graphic.png"
---

# WordPress + Steampipe

[WordPress](https://wordpress.com/) is a popular blogging platform.

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL.

For example:

```sql
select
  id,
  title,
  to_char(date, 'YYYY-MM-DD') as publication_date,
  category
from
  wordpress_post
where 
  date > now() - interval '1 week'
```

```
+----------+----------------------------------------------------------------------------+------------------+---------------------+
| id       | title                                                                      | publication_date | category            |
+----------+----------------------------------------------------------------------------+------------------+---------------------+
| 22755480 | Kubernetes 1.31 Arrives with New Support for AI/ML, Networking             | 2024-08-13       | [9983,9985]         |
| 22753818 | With eLxr, Wind River Brings Debian Linux to the Edge                      | 2024-08-12       | [10064,10932]       |
| 22754958 | Internal Developer Portals: 3 Real World Examples                          | 2024-08-12       | [9984,12869]        |
| 22754984 | Google Angular Lead Sees Convergence in JavaScript Frameworks              | 2024-08-09       | [12406,13055]       |
| 22753474 | In the Age of AI, What Should We Teach Student Programmers?                | 2024-08-07       | [12945,13332]       |
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/wordpress/tables)**

## Quick start

### Install

Download and install the latest WordPress plugin:

```bash
steampipe plugin install wordpress
```

### Credentials

| Item        | Description  |
|-------------|--------------|
| Credentials | All API requests require a WordPress username and password. |
| Permissions | API requests have the same permissions as the logged-in user. |
| Radius      | Each connection represents a single WordPress account. |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/wordpress.spc.spc`) |

### Configuration

Installing the latest wordpress plugin will create a config file (`~/.steampipe/config/wordpress`) with a single connection named `wordpress`:

```hcl
connection "wordpress" {
  plugin = "wordpress"

  # The base URL of your WordPress site's REST API
  # endpoint = "https://your-wordpress-site.com/wp-json"

  # Authentication credentials
  # username = "your_username"
  # password = "your_application_password"

}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-wordpress
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)
