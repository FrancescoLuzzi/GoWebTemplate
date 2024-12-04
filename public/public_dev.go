//go:build !prod
// +build !prod

package public

import (
	"io/fs"
	"os"
	"sync"
)

var (
	onceFS   sync.Once
	assetsFS fs.FS
)

func AssetFs() fs.FS {
	onceFS.Do(func() {
		assetsFS = os.DirFS("public/assets")
	})
	return assetsFS
}
