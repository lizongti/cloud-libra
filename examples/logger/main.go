package main

import (
	"path/filepath"
	"runtime"

	"github.com/cloudlibraries/libra/assets"
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/cloudlibraries/libra/log"
	"github.com/cloudlibraries/libra/logger"
)

func main() {
	_, file, _, _ := runtime.Caller(0)

	projectDir := filepath.Dir(filepath.Dir(filepath.Dir(file)))

	a := assets.New(assets.NewFileSystemProvider(""))

	assetMap, err := a.GetBundle(filepath.Join(projectDir, "examples", "logger", "config"))
	if err != nil {
		panic(err)
	}

	h := hierarchy.New()
	h.Set("ProjectDir", filepath.ToSlash(projectDir))

	if err := h.LoadAssetMap(assetMap); err != nil {
		panic(err)
	}

	run, err := logger.New(h.Sub("logger"))
	if err != nil {
		panic(err)
	}

	log.SetLogger(run)
	log.Println("Hello, World!")
	log.Println(h)
}
