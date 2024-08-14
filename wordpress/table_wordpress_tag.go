package wordpress

import (
	"context"

	"github.com/sogko/go-wordpress"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableWordPressTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wordpress_tag",
		Description: "Represents a tag in WordPress.",
		List: &plugin.ListConfig{
			Hydrate: listTags,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The tag ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The tag name."},
			{Name: "count", Type: proto.ColumnType_INT, Transform: transform.FromValue().Transform(getTagCount), Description: "Count of occurrences of the tag."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromValue().Transform(getTagDescription), Description: "Description of the tag."},
			{Name: "link", Type: proto.ColumnType_STRING, Transform: transform.FromValue().Transform(getTagLink), Description: "Link for the tag."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func listTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	options := &wordpress.TagListOptions{}

	if d.Quals["id"] != nil {
		id := d.EqualsQuals["id"].GetInt64Value()
		options.Include = []int{int(id)}
	}

	err = paginate(ctx, d, func(ctx context.Context, opts interface{}, perPage, offset int) (interface{}, *wordpress.Response, error) {
		options := opts.(*wordpress.TagListOptions)
		options.ListOptions.PerPage = perPage
		options.ListOptions.Offset = offset
		return conn.Tags.List(ctx, options)
	}, options)

	return nil, err
}
