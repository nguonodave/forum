package handlers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// During tests, the templates dir is in the parent directory's parent hierarchy
	templatesDir = filepath.Join("../", templatesDir)
	os.Exit(m.Run())
}
