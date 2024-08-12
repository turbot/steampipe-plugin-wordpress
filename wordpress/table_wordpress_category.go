package wordpress

import (
	"context"

	"github.com/sogko/go-wordpress"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableWordPressCategory(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wordpress_category",
		Description: "Represents a category in WordPress.",
		List: &plugin.ListConfig{
			Hydrate: listCategories,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The category ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The category name."},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "The category slug."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The category description."},
			{Name: "count", Type: proto.ColumnType_INT, Description: "The number of posts in the category."},
			{Name: "link", Type: proto.ColumnType_STRING, Description: "The URL of the category."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func listCategories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	options := &wordpress.CategoryListOptions{}

	if d.Quals["id"] != nil {
		id := d.EqualsQuals["id"].GetInt64Value()
		options.Include = []int{int(id)}
	}

	err = paginate(ctx, d, func(ctx context.Context, opts interface{}, perPage, offset int) (interface{}, *wordpress.Response, error) {
		options := opts.(*wordpress.CategoryListOptions)
		options.ListOptions.PerPage = perPage
		options.ListOptions.Offset = offset
		return conn.Categories.List(ctx, options)
	}, options)

	return nil, err
}
