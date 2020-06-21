package config

import (
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	old := os.Args
	os.Args = []string{}
	defer func{ os.Args = old }()
	t.Skip("Not Implemented")
}
