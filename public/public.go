package public

import (
	"fmt"
	"mime"
	"net/http"
	"path"
)

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

var compressExtensions = map[string]string{
	"gzip": ".gz",
	"br":   ".br",
	"zstd": ".zst",
}

func init() {
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".json", "text/json; charset=utf-8")
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
)

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
