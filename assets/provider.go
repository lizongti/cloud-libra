package assets

type Provider interface {
	// Get returns the asset with the given name.
	Get(name string) ([]byte, error)
	// Keys return all keys in the provider
	Keys() []string
	// List return sub assets in a directory
	List(dir string) ([]string, error)
}
