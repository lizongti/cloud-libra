package assets

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
)

type ProviderBundle struct {
	provider Provider
	path     string
}

func NewProviderBundle(provider Provider, path string) *ProviderBundle {
	return &ProviderBundle{
		provider: provider,
		path:     path,
	}
}

func (pb *ProviderBundle) Get(path string) ([]byte, error) {
	dataMap := make(map[string][]byte)

	if err := pb.getRecursive(pb.path, dataMap); err != nil {
		return nil, err
	}

	return dataMap[path], nil
}

func (pb *ProviderBundle) getRecursive(path string, dataMap map[string][]byte) error {
	// Get names in path
	names, err := pb.provider.List(path)
	if err != nil {
		return err
	}

	// Name has no sub assets
	if len(names) == 0 {
		data, err := pb.provider.Get(path)
		if err != nil {
			return err
		}

		dataMap[path] = data

		return nil
	}

	// Name has sub assets
	for _, name := range names {
		if err := pb.getRecursive(filepath.Join(path, name), dataMap); err != nil {
			return err
		}
	}

	// Add self to dataMap
	data, err := pb.provider.Get(path)
	if err != nil {
		return err
	}

	cryptName := fmt.Sprintf("%x", md5.Sum([]byte(path)))
	dataMap[filepath.Join(path, cryptName)] = data

	return nil
}
