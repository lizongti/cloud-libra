package assets

import (
	"crypto/md5"
	"errors"
	"fmt"
	"path/filepath"
)

var ErrorBundleNotFound = errors.New("bundle not found")

// Bundle get assets from multiple providers.
// If any bundle is in specific provider, this bundle should be complected.
// If not, unexpected error would occur.
type Bundle struct {
	providers []Provider
}

// NewBundle returns a new bundle.
func NewBundle(provider ...Provider) *Bundle {
	return &Bundle{
		providers: provider,
	}
}

func (b *Bundle) Get(path string) (map[string][]byte, error) {
	for _, provider := range b.providers {
		assetMap, err := b.getFromProvider(path, provider)
		if err != nil {
			return nil, err
		}

		if len(assetMap) > 0 {
			return assetMap, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrorBundleNotFound, path)
}

// Get returns the asset with the given name.
func (b *Bundle) getFromProvider(path string, provider Provider) (map[string][]byte, error) {
	assetMap := make(map[string][]byte)

	if err := b.getFromProviderRecursive(path, provider, assetMap); err != nil {
		return nil, err
	}

	return assetMap, nil
}

func (b *Bundle) getFromProviderRecursive(path string, provider Provider, assetMap map[string][]byte) error {
	// Get names in path
	names, err := provider.List(path)
	if err != nil {
		return err
	}

	// Name has no sub assets
	if len(names) == 0 {
		data, err := provider.Get(path)
		if err != nil {
			return err
		}

		assetMap[path] = data

		return nil
	}

	// Name has sub assets
	for _, name := range names {
		if err := b.getFromProviderRecursive(filepath.Join(path, name), provider, assetMap); err != nil {
			return err
		}
	}

	// Add self to dataMap
	data, err := provider.Get(path)
	if err != nil {
		return err
	}

	cryptName := fmt.Sprintf("%x", md5.Sum([]byte(path)))
	assetMap[filepath.Join(path, cryptName)] = data

	return nil
}
