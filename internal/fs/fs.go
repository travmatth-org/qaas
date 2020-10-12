package fs

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

const (
	// ReadAllWriteUser gives specified permissions
	ReadAllWriteUser os.FileMode = 0644
	// ReadExecAllWriteUser gives specified permissions
	ReadExecAllWriteUser os.FileMode = 0755
)

// FS abstracts over the filesystem providing read/write access to
// an in memory filesystem, or a cached os file system
type FS struct {
	client *types.AFS
	files  map[string]types.AFSFile
}

// New returns a new FS
func New() *FS {
	return &FS{nil, make(map[string]types.AFSFile)}
}

// WithMemFs creates an underlying in memory filesystem
func (fs *FS) WithMemFs() *FS {
	fs.client = &afero.Afero{Fs: afero.NewMemMapFs()}
	return fs
}

// WithCachedFs creates a cached os file system
func (fs *FS) WithCachedFs() *FS {
	base := afero.NewOsFs()
	layer := afero.NewMemMapFs()
	// Cache files in the layer for the given
	// time.Duration, a cache duration of 0 means "forever"
	duration := time.Duration(0)
	fs.client = &afero.Afero{Fs: afero.NewCacheOnReadFs(base, layer, duration)}
	return fs
}

// opens specified file in filesystem, locally caches reader under name
func (fs *FS) cache(name, path string) error {
	if _, ok := fs.files[name]; ok {
		return nil
	} else if val, err := fs.client.Open(path); err != nil {
		logger.Error().Err(err).Str("path", path).Msg("Error opening asset")
		return err
	} else {
		fs.files[name] = val
		return nil
	}
}

// LoadAssets walks the given directory, opening and caching assets
func (fs *FS) LoadAssets(dir string) error {
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if filepath.Ext(path) == ".html" {
			name := info.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))
			if err := fs.cache(name, path); err != nil {
				fs.CloseAll()
				return err
			}
			logger.Info().Str("page", path).Msg("Loaded page into memory")
		}
		return nil
	}
	return fs.client.Walk(dir, walk)
}

// Use returns the reader of a file under the specified key
func (fs *FS) Use(key string) types.AFSFile {
	return fs.files[key]
}

// Open a reader for the file at the specified path on the file system
func (fs *FS) Open(path string) (types.AFSFile, error) {
	return fs.client.Open(path)
}

// CloseAll closes all cached files
func (fs *FS) CloseAll() {
	for _, file := range fs.files {
		if err := file.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing asset file")
		} else {
			logger.Info().Str("file", file.Name()).Msg("Asset file closed")
		}
	}
}

func (fs *FS) Write(path string, data []byte, perm os.FileMode) error {
	return fs.client.WriteFile(path, data, perm)
}
