package afs

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

// AFS abstracts over the filesystem providing read/write access to
// an in memory filesystem, or a cached os file system
type AFS struct {
	client *types.AFS
	files  map[string]types.AFSFile
}

// New returns a new FS
func New() *AFS {
	return &AFS{nil, make(map[string]types.AFSFile)}
}

// WithMemFs creates an underlying in memory filesystem
func (afs *AFS) WithMemFs() *AFS {
	afs.client = &afero.Afero{Fs: afero.NewMemMapFs()}
	return afs
}

// WithCachedFs creates a cached os file system
func (afs *AFS) WithCachedFs() *AFS {
	var (
		base  = afero.NewOsFs()
		layer = afero.NewMemMapFs()
		// Cache files in the layer for the given
		// time.Duration, a cache duration of 0 means "forever"
		duration = time.Duration(0)
	)
	afs.client = &afero.Afero{Fs: afero.NewCacheOnReadFs(base, layer, duration)}
	return afs
}

// opens specified file in filesystem, locally caches reader under name
func (afs *AFS) cache(name, path string) error {
	if _, ok := afs.files[name]; ok {
		return nil
	} else if val, err := afs.client.Open(path); err != nil {
		logger.Error().Err(err).Str("path", path).Msg("Error opening asset")
		return err
	} else {
		afs.files[name] = val
		return nil
	}
}

// LoadAssets walks the given directory, opening and caching assets
func (afs *AFS) LoadAssets(dir string) error {
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if filepath.Ext(path) == ".html" {
			name := info.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))
			if err := afs.cache(name, path); err != nil {
				afs.CloseAll()
				return err
			}
			logger.Info().Str("page", path).Msg("Loaded page into memory")
		}
		return nil
	}
	return afs.client.Walk(dir, walk)
}

// Use returns the reader of a file under the specified key
func (afs *AFS) Use(key string) types.AFSFile {
	return afs.files[key]
}

// Open a reader for the file at the specified path on the file system
func (afs *AFS) Open(path string) (types.AFSFile, error) {
	return afs.client.Open(path)
}

// CloseAll closes all cached files
func (afs *AFS) CloseAll() {
	for _, file := range afs.files {
		if err := file.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing asset file")
		} else {
			logger.Info().Str("file", file.Name()).Msg("Asset file closed")
		}
	}
}

// ReadFile opens file on the filesytem into a []byte array
func (afs *AFS) ReadFile(path string) ([]byte, error) {
	return afs.client.ReadFile(path)
}

// Write saves file to file system
func (afs *AFS) Write(path string, data []byte, perm os.FileMode) error {
	return afs.client.WriteFile(path, data, perm)
}
