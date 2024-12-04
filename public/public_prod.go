//go:build prod
// +build prod

package public

import (
	"embed"
	"io/fs"
	"sync"
)

//go:embed assets
var staticAssetsFS embed.FS

var (
	onceFS   sync.Once
	assetsFS fs.FS
)

func AssetFs() fs.FS {
	onceFS.Do(func() {
		assetsFS, _ = fs.Sub(staticAssetsFS, "assets")
	})
	return assetsFS
}
