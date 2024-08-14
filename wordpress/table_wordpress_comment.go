package wordpress

import (
	"context"

	"github.com/sogko/go-wordpress"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableWordPressComment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wordpress_comment",
		Description: "Represents a comment in WordPress.",
		List: &plugin.ListConfig{
			Hydrate: listComments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
				{Name: "date", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
				{Name: "post", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The comment ID."},
			{Name: "post", Type: proto.ColumnType_INT, Description: "The ID of the post the comment is associated with."},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The name of the comment author."},
			{Name: "content", Type: proto.ColumnType_STRING, Transform: transform.FromValue().Transform(getCommentContent), Description: "The content of the comment."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().Transform(getCommentDate), Description: "The comment publication date."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the comment (approved, pending, spam, etc.)."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

func listComments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("WordPress listComments id", "id", d.Quals["id"])
	plugin.Logger(ctx).Debug("WordPress listComments date", "author", d.Quals["date"])
	plugin.Logger(ctx).Debug("WordPress listComments post", "post", d.Quals["post"])

	options := &wordpress.CommentListOptions{}

	if d.Quals["id"] != nil {
		id := d.EqualsQuals["id"].GetInt64Value()
		options.Post = []int{int(id)}
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

	plugin.Logger(ctx).Debug("WordPress listComments API request options", "options", options)

	err = paginate(ctx, d, func(ctx context.Context, opts interface{}, perPage, offset int) (interface{}, *wordpress.Response, error) {
		options := opts.(*wordpress.CommentListOptions)
		options.ListOptions.PerPage = perPage
		options.ListOptions.Offset = offset
		return conn.Comments.List(ctx, options)
	}, options)

	return nil, err
}
