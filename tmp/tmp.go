package main

import (
	"path/filepath"

	"github.com/cloudlibraries/libra/assets"
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/cloudlibraries/libra/log"
	"github.com/cloudlibraries/libra/logger"
)

func main() {
	projectDir := "D:\\github.com\\cloudlibraries\\libra"

	a := assets.New(assets.NewFileSystemProvider(""))
	assetMap, err := a.GetBundle(filepath.Join(projectDir, "config"))
	if err != nil {
		panic(err)
	}

	h := hierarchy.New()
	h.Set("ProjectDir", filepath.ToSlash(projectDir))
	if err := h.LoadAssetMap(assetMap); err != nil {
		panic(err)
	}

	run, err := logger.New(h.Child("logger"))
	if err != nil {
		panic(err)
	}

	log.SetLogger(run)

	log.Println("test")
	log.Println(hierarchy.New())
}
