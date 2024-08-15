package wordpress

import (
	"context"

	"github.com/sogko/go-wordpress"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableWordPressPost(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wordpress_post",
		Description: "Represents a post in WordPress.",
		List: &plugin.ListConfig{
			Hydrate: listPosts,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "author", Require: plugin.Optional},
				{Name: "date", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The post ID."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromValue().Transform(getPostTitle), Description: "The post title."},
			{Name: "link", Type: proto.ColumnType_STRING, Description: "The post link."},
			{Name: "content", Type: proto.ColumnType_STRING, Transform: transform.FromValue().Transform(getPostContent), Description: "The post content."},
			{Name: "author", Type: proto.ColumnType_INT, Description: "The post author ID."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().Transform(getPostDate), Description: "The post publication date."},
			{Name: "category", Type: proto.ColumnType_JSON, Transform: transform.FromValue().Transform(getPostCategory), Description: "The categories assigned to the post."},
			{Name: "tag", Type: proto.ColumnType_JSON, Transform: transform.FromValue().Transform(getPostTag), Description: "The tags assigned to the post."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func listPosts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("WordPress listPosts author", "author", d.Quals["author"])
	plugin.Logger(ctx).Debug("WordPress listPosts date", "date", d.Quals["date"])

	options := &wordpress.PostListOptions{}

	if d.Quals["author"] != nil {
		id := d.EqualsQuals["author"].GetInt64Value()
		options.Author = []int{int(id)}
	}

	if d.Quals["date"] != nil {
		for _, q := range d.Quals["date"].Quals {
			switch q.Operator {
			case ">=", ">":
				t := q.Value.GetTimestampValue().AsTime()
				options.After = &t
			case "<=", "<":
				t := q.Value.GetTimestampValue().AsTime()
				options.Before = &t
			}
		}
	}

	plugin.Logger(ctx).Debug("WordPress listPosts API request options", "options", options)

	err = paginate(
		ctx, // Context for the operation
		d,   // QueryData containing information about the current query

		// This is the ListFunc - an anonymous function that wraps the actual API call
		func(ctx context.Context, opts interface{}, perPage, offset int) (interface{}, *wordpress.Response, error) {
			// Type assert the generic opts to the specific PostListOptions
			// This is necessary because opts is passed as interface{} for flexibility
			options := opts.(*wordpress.PostListOptions)

			// Set up pagination parameters
			options.ListOptions.PerPage = perPage // Number of items per page
			options.ListOptions.Offset = offset   // Starting position for this page

			// Set up sorting parameters
			options.ListOptions.OrderBy = "id" // Sort by post ID
			options.ListOptions.Order = "asc"  // Sort in ascending order

			// Make the actual API call using the WordPress SDK
			// conn.Posts.List returns ([]wordpress.Post, *wordpress.Response, error)
			// But we return it as (interface{}, *wordpress.Response, error) to match ListFunc signature
			return conn.Posts.List(ctx, options)
		},

		// Pass in the initial options (filters, etc.) set up earlier in the listPosts function
		options)

	return nil, err
}
