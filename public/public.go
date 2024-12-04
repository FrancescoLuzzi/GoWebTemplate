package public

import (
	"fmt"
	"mime"
	"path"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

var (
	compressExtensions map[string]string
	onceExtensions     sync.Once
)

func init() {
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".json", "text/json; charset=utf-8")
}

func getCompressExtensions(ctx fiber.Ctx) map[string]string {
	onceExtensions.Do(func() {
		compressExtensions = ctx.App().Config().CompressedFileSuffixes
	})
	return compressExtensions
}

func searchExtensionEncoding(ext string, encodings map[string]string) string {
	for encName, encExt := range encodings {
		if encExt == ext {
			return encName
		}
	}
	return ""
}

const (
	contentEncodingHeader = "Content-Encoding"
	contentTypeHeader     = "Content-Type"
	cacheControlHeader    = "Cache-Control"
)

func RegisterAssets(app *fiber.App) {
	app.Use("/public/assets/*", static.New("", static.Config{
		FS:       AssetFs(),
		Compress: true,
		ModifyResponse: func(c fiber.Ctx) error {
			urlPath := c.Request().URI().String()
			ext := path.Ext(urlPath)
			// get last extension ES: .js.gz -> .gz
			enc := searchExtensionEncoding(ext, getCompressExtensions(c))
			if enc != "" {
				c.Set(contentEncodingHeader, enc)
				enc = ""
				// get compressed file extension ES: .js.gz -> .js
				ext = path.Ext(urlPath[:len(urlPath)-len(ext)])
				enc = mime.TypeByExtension(ext)
				if enc != "" {
					c.Set(contentTypeHeader, enc)
				}
			}
			return nil
		},
	}))
}
