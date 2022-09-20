package assets

import (
	"errors"
	"fmt"
	"path/filepath"
)

var (
	ErrAssetNotFound  = errors.New("asset not found")
	ErrBundleNotFound = errors.New("bundle not found")
)

func IsBundleNotFound(err error) bool {
	return errors.Is(err, ErrBundleNotFound)
}

func IsAssetNotFound(err error) bool {
	return errors.Is(err, ErrAssetNotFound)
}

// Assets get assets from multiple providers.
type Assets struct {
	providers []Provider
}

// NewAssets returns a new Assets.
func NewAssets(providers ...Provider) *Assets {
	return &Assets{
		providers: providers,
	}
}

// AddProvider adds a new provider to the list of providers.
func (a *Assets) AddProvider(p Provider) {
	a.providers = append(a.providers, p)
}

// Get returns the asset with the given name.
func (a *Assets) GetAsset(name string) ([]byte, error) {
	name = slash(name)
	for _, provider := range a.providers {
		data, err := provider.Get(name)

		switch {
		case IsAssetNotFound(err):
			continue
		case err != nil:
			return nil, err
		default:
			return data, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrAssetNotFound, name)
}

// GetBundle returns the bundle assets with the given name.
func (a *Assets) GetBundle(name string) (map[string][]byte, error) {
	name = slash(name)
	for _, provider := range a.providers {
		assetMap, err := NewBundle(name, provider).Get()
		if err != nil {
			return nil, err
		}

		if len(assetMap) > 0 {
			return assetMap, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrBundleNotFound, name)
}

func slash(names ...string) string {
	return filepath.ToSlash(filepath.Join(names...))
}
