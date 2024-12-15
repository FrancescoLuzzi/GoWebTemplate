package public

import (
	"fmt"
	"mime"
	"net/http"
	"path"
)

func init() {
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".json", "text/json; charset=utf-8")
}

var compressExtensions = map[string]string{
	"gzip": ".gz",
	"br":   ".br",
	"zstd": ".zst",
}

const (
	contentEncodingHeader = "Content-Encoding"
	contentTypeHeader     = "Content-Type"
	cacheControlHeader    = "Cache-Control"
)

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

func searchExtensionEncoding(ext string, encodings map[string]string) string {
	for encName, encExt := range encodings {
		if encExt == ext {
			return encName
		}
	}
	return ""
}

func FixCompressedContentHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		ext := path.Ext(urlPath)
		// get last extension ES: .js.gz -> .gz
		enc := searchExtensionEncoding(ext, compressExtensions)
		if enc != "" {
			w.Header().Set(contentEncodingHeader, enc)
			enc = ""
			// get compressed file extension ES: .js.gz -> .js
			ext = path.Ext(urlPath[:len(urlPath)-len(ext)])
			enc = mime.TypeByExtension(ext)
			if enc != "" {
				w.Header().Set(contentTypeHeader, enc)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func DisableCacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheControlHeader, "no-store")
		next.ServeHTTP(w, r)
	})
}

func CacheHandler(next http.Handler, maxAge uint) http.Handler {
	cacheValue := fmt.Sprintf("max-age=%d, must-revalidate", maxAge)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheControlHeader, cacheValue)
		next.ServeHTTP(w, r)
	})
}
