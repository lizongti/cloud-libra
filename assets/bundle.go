package assets

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"path/filepath"
)

var ErrBundleNotFound = errors.New("bundle not found")

func IsBundleNotFound(err error) bool {
	return errors.Is(err, ErrBundleNotFound)
}

// Bundle get assets from multiple providers.
// If any bundle is in specific provider, this bundle should be complected.
// If not, unexpected error would occur.
type Bundle struct {
	name      string
	providers []Provider
}

// NewBundle returns a new bundle.
func NewBundle(name string, provider ...Provider) *Bundle {
	return &Bundle{
		providers: provider,
	}
}

func (b *Bundle) Get() (assetMap map[string][]byte, err error) {
	for _, provider := range b.providers {
		assetMap, err = NewProviderBundle(b.name, provider).Get()
		if err != nil {
			return nil, err
		}

		if len(assetMap) > 0 {
			return assetMap, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrBundleNotFound, b.name)
}

type ProviderBundle struct {
	name     string
	provider Provider
	assetMap map[string][]byte
}

func NewProviderBundle(name string, provider Provider) *ProviderBundle {
	return &ProviderBundle{
		name:     name,
		provider: provider,
		assetMap: make(map[string][]byte),
	}
}

// Get returns the asset with the given name.
func (pb *ProviderBundle) Get() (map[string][]byte, error) {
	if err := pb.getRecursive(pb.name); err != nil {
		return nil, err
	}

	return pb.assetMap, nil
}

func (pb *ProviderBundle) getRecursive(name string) error {
	// Get names by listing name
	names, err := pb.provider.List(name)
	if err != nil {
		return err
	}

	// Name has no sub assets
	if len(names) == 0 {
		return pb.getFile(name)
	}

	// Name has sub assets
	for _, name := range names {
		if err := pb.getRecursive(Join(name, name)); err != nil {
			return err
		}
	}

	// Add self to asset map
	asset, err := pb.provider.Get(name)
	if IsAssetNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	cryptName := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	pb.assetMap[Join(name, cryptName)] = asset

	return nil
}

func (pb *ProviderBundle) getFile(name string) error {
	switch filepath.Ext(name) {
	case ".zip":
		return pb.getZip(name)
	default:
		return pb.getRaw(name)
	}
}

func (pb *ProviderBundle) getRaw(name string) error {
	data, err := pb.provider.Get(name)
	if IsAssetNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	pb.assetMap[name] = data

	return nil
}

func (pb *ProviderBundle) getZip(name string) error {
	data, err := pb.provider.Get(name)
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

		pb.assetMap[Join(name, file.Name)] = buf.Bytes()
	}

	return nil
}
