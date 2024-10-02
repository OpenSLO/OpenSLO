package internal

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	moduleRoot string
	once       sync.Once
)

// FindModuleRoot finds the root of the current module.
// It does so by looking for a go.mod file in the current working directory.
func FindModuleRoot() string {
	once.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dir = filepath.Clean(dir)
		for {
			if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
				moduleRoot = dir
				return
			}
			d := filepath.Dir(dir)
			if d == dir {
				break
			}
			dir = d
		}
	})
	return moduleRoot
}
