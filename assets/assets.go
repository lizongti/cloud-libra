package assets

import (
	"fmt"
)

var ErrAssetNotFound = fmt.Errorf("asset not found")

type Assets struct {
	providers []Provider
}

// AddProvider adds a new provider to the list of providers.
func (a *Assets) AddProvider(p Provider) {
	a.providers = append(a.providers, p)
}

// Get returns the asset with the given name.
func (a *Assets) Get(name string) ([]byte, error) {
	for _, p := range a.providers {
		if b, err := p.Get(name); err == nil {
			return b, nil
		}
	}

	return nil, fmt.Errorf("%w:%s", ErrAssetNotFound, name)
}

// GetBundle returns the assets bundle with the given name.
// func (a *Assets) GetBundle(name string) (*Bundle, error) {
