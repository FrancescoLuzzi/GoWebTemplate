package cache

import (
	"bytes"
	"context"
	"encoding"
	"io"
	"log/slog"
	"time"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/app_ctx"
	"github.com/a-h/templ"
)

type cachedPage struct {
	key string
	ttl time.Duration
}

type pageCacheValue struct {
	buf bytes.Buffer
}

var _ encoding.BinaryUnmarshaler = &pageCacheValue{}
var _ encoding.BinaryMarshaler = &pageCacheValue{}

// MarshalBinary implements encoding.BinaryMarshaler.
func (c *pageCacheValue) MarshalBinary() (data []byte, err error) {
	return c.buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (c *pageCacheValue) UnmarshalBinary(data []byte) error {
	fullDataLen := len(data)
	var err error
	for tmpWrite, err := c.buf.Write(data); err != nil && fullDataLen-tmpWrite > 0; {
		fullDataLen -= tmpWrite
		data = data[tmpWrite:]
	}
	return err
}

func (c *cachedPage) Render(ctx context.Context, w io.Writer) error {
	cache := app_ctx.Cache(ctx)
	slog.Info("Rendering cached page", "key", c.key, "ttl", c.ttl, "cache", cache != nil)
	children := templ.GetChildren(ctx)
	ctx = templ.ClearChildren(ctx)
	if children == nil {
		return nil
	}
	if cache == nil {
		return children.Render(ctx, w)
	}
	page := pageCacheValue{}
	err := cache.Get(ctx, c.key, &page)
	if err == nil {
		slog.Info("Cache hit", "key", c.key, "ttl", c.ttl, "page_length", page.buf.Len())
		_, err := w.Write(page.buf.Bytes())
		return err
	}
	// Render children to a buffer.
	err = children.Render(ctx, &page.buf)
	if err != nil {
		return err
	}
	// Cache the result.
	cache.Set(ctx, c.key, &page, c.ttl)
	slog.Info("Cache miss", "key", c.key, "ttl", c.ttl, "page_length", page.buf.Len())
	// Write the result to the output.
	_, err = w.Write(page.buf.Bytes())
	return err
}

// Cache a component by key for a specific duration
//
// Example usage:
//
//	templ.Cache("my-key", 5*time.Minute){
//			<p> Hello World </p>
//	}
func Cache(key string, duration time.Duration) templ.Component {
	return &cachedPage{
		key: key,
		ttl: duration,
	}
}
