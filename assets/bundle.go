package assets

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"path/filepath"
)

type Bundle struct {
	name     string
	provider Provider
	assetMap map[string][]byte
}

func NewBundle(name string, provider Provider) *Bundle {
	return &Bundle{
		name:     slash(name),
		provider: provider,
		assetMap: make(map[string][]byte),
	}
}

// Get returns the asset with the given name.
func (b *Bundle) Get() (map[string][]byte, error) {
	if err := b.getRecursive(b.name); err != nil {
		return nil, err
	}

	return b.assetMap, nil
}

func (b *Bundle) getRecursive(name string) error {
	// Get names by listing name
	names, err := b.provider.List(name)
	if err != nil {
		return err
	}

	// Name has no sub assets
	if len(names) == 0 {
		return b.getFile(name)
	}

	// Name has sub assets
	for _, subName := range names {
		if err := b.getRecursive(slash(name, subName)); err != nil {
			return err
		}
	}

	// Add self to asset map
	asset, err := b.provider.Get(name)
	if IsAssetNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	cryptName := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	b.assetMap[slash(name, cryptName)] = asset

	return nil
}

func (b *Bundle) getFile(name string) error {
	switch filepath.Ext(name) {
	case ".zip":
		return b.getZip(name)
	default:
		return b.getRaw(name)
	}
}

func (b *Bundle) getRaw(name string) error {
	data, err := b.provider.Get(name)
	if IsAssetNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	b.assetMap[name] = data

	return nil
}

func (b *Bundle) getZip(name string) error {
	data, err := b.provider.Get(name)
	if IsAssetNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(fileReader); err != nil {
			return err
		}

		b.assetMap[slash(filepath.Dir(name), file.Name)] = buf.Bytes()
	}

	return nil
}
