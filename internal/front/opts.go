package front

import (
	"net/http"

	"github.com/codemodus/formlark/internal/spa/assets"
	assetfs "github.com/elazarl/go-bindata-assetfs"
)

type options struct {
	fs http.Handler
}

// Option ...
type Option func(*options) error

// WithFileSystemAssets ...
func WithFileSystemAssets() Option {
	return func(opts *options) error {
		opts.fs = http.FileServer(http.Dir("internal/spa/src/srv"))
		return nil
	}
}

// WithEmbeddedAssets ...
func WithEmbeddedAssets() Option {
	return func(opts *options) error {
		opts.fs = defaultFileServer()

		return nil
	}
}

func defaultFileServer() http.Handler {
	return http.FileServer(
		&assetfs.AssetFS{
			Asset:     assets.Asset,
			AssetDir:  assets.AssetDir,
			AssetInfo: assets.AssetInfo,
			Prefix:    "",
		},
	)
}
