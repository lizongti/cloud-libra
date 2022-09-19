package assets

import (
	"errors"
	"fmt"
	"path/filepath"
)

var ErrAssetNotFound = errors.New("asset not found")

func IsAssetNotFound(err error) bool {
	return errors.Is(err, ErrAssetNotFound)
}

type Assets struct {
	providers []Provider
}

// AddProvider adds a new provider to the list of providers.
func (a *Assets) AddProvider(p Provider) {
	a.providers = append(a.providers, p)
}

// Get returns the asset with the given name.
func (a *Assets) GetAsset(name string) ([]byte, error) {
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
	return NewBundle(name, a.providers...).Get()
}

func Join(names ...string) string {
	return filepath.ToSlash(filepath.Join(names...))
}
