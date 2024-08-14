package wordpress

import (
	"context"

	"github.com/sogko/go-wordpress"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableWordPressAuthor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wordpress_author",
		Description: "Represents an author in WordPress.",
		List: &plugin.ListConfig{
			Hydrate: listAuthors,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The author ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The author's name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The author's description."},
			{Name: "link", Type: proto.ColumnType_STRING, Description: "The author's link."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The author's url."},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "The author's slug."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func listAuthors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("WordPress listAuthors id", "id", d.Quals["id"])

	options := &wordpress.UserListOptions{}

	if d.Quals["id"] != nil {
		id := d.EqualsQuals["id"].GetInt64Value()
		options.Include = []int{int(id)}
	}

	plugin.Logger(ctx).Debug("WordPress listAuthors API request options", "options", options)

	err = paginate(ctx, d, func(ctx context.Context, opts interface{}, perPage, offset int) (interface{}, *wordpress.Response, error) {
		options := opts.(*wordpress.UserListOptions)
		options.ListOptions.PerPage = perPage
		options.ListOptions.Offset = offset
		return conn.Users.List(ctx, options)
	}, options)

	return nil, err
}
