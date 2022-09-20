package main

import (
	"fmt"

	"github.com/cloudlibraries/libra/assets"
	"github.com/cloudlibraries/libra/osutil"
)

func main() {
	// fmt.Println(assets.NewAssets(
	// 	assets.NewFileSystemProvider(""),
	// ).GetAsset("assets/assets.go"))
	bundleAssets, _ := assets.NewAssets(
		assets.NewFileSystemProvider(osutil.GetProjectPath("libra")),
	).GetBundle("assets\\assets.go")
	for name, asset := range bundleAssets {
		fmt.Println(name, len(asset))
	}
}
