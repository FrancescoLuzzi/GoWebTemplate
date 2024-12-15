package cache

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/a-h/templ"
	"github.com/hashicorp/golang-lru/v2"
)

// golang-lru is thread safe
var cacheKeyToContent, _ = lru.New[string, cacheContent](16)

type cacheContent struct {
	expireAt time.Time
	content  []byte
}

type cacheComponent struct {
	cacheDuration time.Duration
	key           string
}

func (c cacheComponent) Render(ctx context.Context, w io.Writer) error {
	// Get children.
	children := templ.GetChildren(ctx)
	ctx = templ.ClearChildren(ctx)
	if children == nil {
		return nil
	}
	cc, isCached := cacheKeyToContent.Get(c.key)
	if isCached {
		if cc.expireAt.After(time.Now()) {
			_, err := w.Write(cc.content)
			return err
		} else {
			cacheKeyToContent.Remove(c.key)
		}
	}
	// Render children to a buffer.
	var buf bytes.Buffer
	err := children.Render(ctx, &buf)
	if err != nil {
		return err
	}
	// Cache the result.
	cacheKeyToContent.Add(c.key, cacheContent{
		expireAt: time.Now().Add(c.cacheDuration),
		content:  buf.Bytes(),
	})
	// Write the result to the output.
	_, err = w.Write(buf.Bytes())
	return err
}

// Cache a component by key for a specific duration
func Cache(key string, duration time.Duration) templ.Component {
	return cacheComponent{
		cacheDuration: duration,
		key:           key,
	}
}
