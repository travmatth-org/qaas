package types

import (
	"github.com/spf13/afero"
)

// internal/config/config.go
// internal/afs/afs.go

// AFS ...
type AFS = afero.Afero

// AFSFile ...
type AFSFile = afero.File
