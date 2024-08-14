package wordpress

import (
	"context"
  "reflect"
	//"sync"
	//"sync/atomic"
	//"time"

	"github.com/sogko/go-wordpress"	
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"	
)

func getPostDate(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  post := d.Value.(*wordpress.Post)
	date := post.Date.Time
	return date, nil
}

func getPostTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  post := d.Value.(*wordpress.Post)
	title := post.Title.Rendered
	return title, nil
}

func getPostCategory(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  post := d.Value.(*wordpress.Post)
	categories := post.Categories
	return categories, nil
}

func getPostTag(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  post := d.Value.(*wordpress.Post)
	tags := post.Tags
	return tags, nil
}

func getPostContent(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  post := d.Value.(*wordpress.Post)
	date := post.Content.Rendered
	return date, nil
}

func getCommentDate(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  comment := d.Value.(*wordpress.Comment)
	date := comment.Date.Time
	return date, nil
}

func getCommentContent(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  comment := d.Value.(*wordpress.Comment)
	content := comment.Content.Rendered
	return content, nil
}

func getTagCount(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  tag := d.Value.(*wordpress.Tag)
	count := tag.Count
	return count, nil
}

func getTagDescription(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  tag := d.Value.(*wordpress.Tag)
	description := tag.Description
	return description, nil
}

func getTagLink(ctx context.Context, d *transform.TransformData) (interface{}, error) {
  tag := d.Value.(*wordpress.Tag)
	link := tag.Link
	return link, nil
}



type ListFunc func(context.Context, interface{}, int, int) (interface{}, *wordpress.Response, error)

func paginate(ctx context.Context, d *plugin.QueryData, listFunc ListFunc, options interface{}) error {
	perPage := 100
	offset := 0

	for {
			plugin.Logger(ctx).Debug("WordPress paginate", "offset", offset)

			items, _, err := listFunc(ctx, options, perPage, offset)
			if err != nil {
					plugin.Logger(ctx).Debug("wordpress.paginate", "query_error", err)
					return err
			}

			itemsSlice := reflect.ValueOf(items)
			for i := 0; i < itemsSlice.Len(); i++ {
					d.StreamListItem(ctx, itemsSlice.Index(i).Interface())
			}

			// If fewer items than perPage were returned, it's the last page
			if itemsSlice.Len() < perPage {
					break
			}

			// Update the offset for the next page
			offset += perPage
	}

	return nil
}


/*
func paginate(ctx context.Context, d *plugin.QueryData, listFunc ListFunc, options interface{}) error {
	perPage := 100
	concurrencyLimit := 10
	ch := make(chan int, concurrencyLimit)
	var wg sync.WaitGroup
	var done int32 = 0 // Use atomic for done flag as well
	delay := 100 * time.Millisecond

	// Launch workers
	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for offset := range ch {
				if atomic.LoadInt32(&done) == 1 || ctx.Err() != nil {
					break
				}

				// Introduce a delay before making the API call
				time.Sleep(delay)

				plugin.Logger(ctx).Debug("WordPress paginate", "offset", offset)

				items, _, err := listFunc(ctx, options, perPage, offset)
				if err != nil {
					plugin.Logger(ctx).Error("wordpress.paginate", "query_error", err)
					return
				}

				itemsSlice := reflect.ValueOf(items)
				plugin.Logger(ctx).Debug("WordPress paginate", "offset", offset, "items", itemsSlice.Len())

				for i := 0; i < itemsSlice.Len(); i++ {
					d.StreamListItem(ctx, itemsSlice.Index(i).Interface())
				}

				if itemsSlice.Len() < perPage {
					atomic.StoreInt32(&done, 1)
				}
			}
		}()
	}

	// Feed the worker pool with ordered offsets
	offset := 0
	for {
		if atomic.LoadInt32(&done) == 1 || ctx.Err() != nil {
			break
		}

		ch <- offset

		offset += perPage
	}

	// Close the channel to signal workers to exit
	close(ch)

	// Wait for all workers to complete
	wg.Wait()

	return nil
}

func paginate(ctx context.Context, d *plugin.QueryData, listFunc ListFunc, options interface{}) error {
	perPage := 100
	var offset int64 = 0 // Use int64 to ensure compatibility with atomic operations
	concurrencyLimit := 10
	ch := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup
	var done int32 = 0 // Use atomic for done flag as well
	delay := 100 * time.Millisecond 
	for {
			if atomic.LoadInt32(&done) == 1 || ctx.Err() != nil {
					break
			}

			ch <- struct{}{}
			wg.Add(1)

			go func() {
					defer func() {
							<-ch
							wg.Done()
					}()

					currentOffset := atomic.AddInt64(&offset, int64(perPage)) - int64(perPage)

					// Introduce a delay before making the API call
					time.Sleep(delay)

					plugin.Logger(ctx).Debug("WordPress paginate", "offset", currentOffset)

					items, _, err := listFunc(ctx, options, perPage, int(currentOffset))
					if err != nil {
							plugin.Logger(ctx).Error("wordpress.paginate", "query_error", err)
							return
					}

					itemsSlice := reflect.ValueOf(items)
					plugin.Logger(ctx).Debug("WordPress paginate", "offset", currentOffset, "items", itemsSlice.Len())

					for i := 0; i < itemsSlice.Len(); i++ {
						d.StreamListItem(ctx, itemsSlice.Index(i).Interface())
					}

					if itemsSlice.Len() < perPage {
							atomic.StoreInt32(&done, 1)
					}
			}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return nil
}

func paginate(ctx context.Context, d *plugin.QueryData, listFunc ListFunc, options interface{}) error {
	perPage := 100
	offset := 0

	// Use a buffered channel to limit the number of concurrent requests
	concurrencyLimit := 10
	ch := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to synchronize access to shared variables
	done := false

	for {
			mu.Lock()
			if done {
					mu.Unlock()
					break
			}
			currentOffset := offset
			offset += perPage
			mu.Unlock()

			ch <- struct{}{}
			wg.Add(1)

			go func(offset int) {
					defer func() {
							<-ch
							wg.Done()
					}()

					plugin.Logger(ctx).Debug("WordPress paginate", "offset", offset)

					items, _, err := listFunc(ctx, options, perPage, offset)
					if err != nil {
							plugin.Logger(ctx).Error("wordpress.paginate", "query_error", err)
							return
					}

					itemsSlice := reflect.ValueOf(items)
					for i := 0; i < itemsSlice.Len(); i++ {
							d.StreamListItem(ctx, itemsSlice.Index(i).Interface())
					}

					// If fewer items than perPage were returned, it's the last page
					if itemsSlice.Len() < perPage {
							mu.Lock()
							done = true
							mu.Unlock()
					}
			}(currentOffset)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return nil
}
*/