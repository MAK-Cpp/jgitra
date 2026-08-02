//go:build clean

package main

import (
	"os"
	"path/filepath"

	"jgitra/internal/config"
)

func main() {
	_ = os.Remove(filepath.Join("bin", config.App))
}
