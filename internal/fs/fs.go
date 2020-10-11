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
	ReadAllWriteUser     os.FileMode = 0644
	ReadExecAllWriteUser os.FileMode = 0755
)

type FS struct {
	client *types.AFS
	files  map[string]types.AFSFile
}

func New() *FS {
	return &FS{nil, make(map[string]types.AFSFile)}
}

func (fs *FS) WithMemFs() *FS {
	fs.client = &afero.Afero{Fs: afero.NewMemMapFs()}
	return fs
}

func (fs *FS) WithCachedFs() *FS {
	base := afero.NewOsFs()
	layer := afero.NewMemMapFs()
	// Cache files in the layer for the given
	// time.Duration, a cache duration of 0 means "forever"
	duration := time.Duration(0)
	fs.client = &afero.Afero{Fs: afero.NewCacheOnReadFs(base, layer, duration)}
	return fs
}

func (fs *FS) Open(name, path string) error {
	if _, ok := fs.files[name]; ok {
		return nil
	} else if val, err := fs.client.Open(path); err != nil {
		logger.Error().Str("file", path).Str("name", name).Err(err).Msg("Error opening asset in directory")
		return err
	} else {
		fs.files[name] = val
		return nil
	}
}

func (fs *FS) LoadAssets(dir string) error {
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if filepath.Ext(path) == ".html" {
			name := info.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))
			if err := fs.Open(name, path); err != nil {
				fs.CloseAll()
				return err
			}
			logger.Info().Str("page", path).Msg("Loaded page into memory")
		}
		return nil
	}
	return fs.client.Walk(dir, walk)
}

func (fs *FS) Use(key string) types.AFSFile {
	return fs.files[key]
}

func (fs *FS) Locate(env string) func() (types.AFSFile, error) {
	return func() (types.AFSFile, error) {
		if path := os.Getenv(env); path != "" {
			return fs.client.Open(path)
		}
		path, err := filepath.Abs(filepath.Join("etc", "qaas", "httpd.yml"))
		if err != nil {
			return nil, err
		}
		return fs.client.Open(path)
	}
}

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
