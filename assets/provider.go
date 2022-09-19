package assets

import (
	"fmt"
	"os"
	"path/filepath"
)

// Provider is the interface that wraps the Get method.
//
// Examples use the following hierarchy:
//
//		|- data/
//			|- com.zip
//			|- foo.txt
//			|- img/
//	  			|- a.png
//			|- js/
//
// com.zip file use the following hierarchy:
//
//		|- com.txt
//		|- press.txt
//	 	|- tense/
//			|- ed.txt
type Provider interface {
	// Get returns the asset with the given name.
	// Get("foo.txt") would return ErrAssetNotFound
	// Get("data/foo.txt") would return the content of foo.txt
	// Get("data/img/a.png") would return the content of img/a.png
	// Get("data/com.zip") would return the content of com.zip
	// Get("data/img") would return ErrAssetNotFound
	// Get("data/img/") would return ErrAssetNotFound
	// Get("data/js") would return ErrAssetNotFound
	// Exception on Get would return unexpected error.
	Get(string) ([]byte, error)

	// List return sub assets in the given name.
	// List("js") would return []string{}
	// List("data") would return []string{"foo.txt", "img", "com.zip"}
	// List("data/img") would return []string{"a.png"}
	// List("data/img/") would return []string{"a.png"}
	// List("data/img/a.png") would return []string{}
	// List("data/foo.txt") would return []string{}
	// List("data/com.zip") would return []string{"com.txt", "press.txt", "tense/ed.txt"}
	// List("data/js") would return []string{}
	// Exception on List would return unexpected error.
	List(string) ([]string, error)
}

// FileSystemProvider is a Provider that uses the file system.
type FileSystemProvider struct {
	root string
}

// NewFileSystemProvider returns a new FileSystemProvider.
func NewFileSystemProvider(root string) *FileSystemProvider {
	return &FileSystemProvider{
		root: root,
	}
}

// Get return the file content from os in the given name.
func (p *FileSystemProvider) Get(name string) ([]byte, error) {
	path := filepath.Join(p.root, name)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrAssetNotFound, name)
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

// List return the sub assets in the given name.
func (p *FileSystemProvider) List(name string) ([]string, error) {
	path := filepath.Join(p.root, name)

	infos, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(infos))
	for _, info := range infos {
		names = append(names, info.Name())
	}

	return names, nil
}

type BinDataProvider struct {
	Asset    func(string) ([]byte, error)
	AssetDir func(string) ([]string, error)
}

func NewBinDataProviderProvider() *BinDataProvider {
	return &BinDataProvider{}
}

func (p *BinDataProvider) Get(name string) ([]byte, error) {
	data, err := p.Asset(name)
	if err.Error() == fmt.Sprintf("Asset %s not found", name) {
		return nil, fmt.Errorf("%w: %s", ErrAssetNotFound, name)
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *BinDataProvider) List(name string) ([]string, error) {
	names, err := p.AssetDir(name)
	if err.Error() == fmt.Sprintf("Asset %s not found", name) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return names, nil
}

type ETCD3Provider struct{}

func NewETCD3Provider() *ETCD3Provider {
	return &ETCD3Provider{}
}

func (p *ETCD3Provider) Get(name string) ([]byte, error) {
	panic("implement me")
}

func (p *ETCD3Provider) List(name string) ([]string, error) {
	panic("implement me")
}
