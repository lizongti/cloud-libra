package file_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudlibraries/libra/internal/repo/file"
	"github.com/mitchellh/go-homedir"
)

func TestPath(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	path := file.NewPath(filepath.Join(home, ".libra"))
	dir, err := path.Directory()
	if err != nil {
		t.Fatal(err)
	}
	path1 := filepath.Join("dir_A", "dir_B", "file_C")
	path2 := filepath.Join("dir_A", "file_B")
	srcDataMap := map[string][]byte{
		path1: []byte(path1),
		path2: []byte(path2),
	}
	if err := dir.Write(srcDataMap); err != nil {
		t.Fatal(err)
	}
	fileMap, err := dir.Files()
	if err != nil {
		t.Fatal(err)
	}
	if len(fileMap) != len(srcDataMap) {
		t.Fatalf("expected fileMap length equals srcDataMap length, fileMap length: %d, srcDataMap length: %d", len(fileMap), len(srcDataMap))
	}
	for path, file := range fileMap {
		data, err := file.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != path {
			t.Fatalf("expected path equals to content, path: %s, content: %s", string(data), path)
		}
	}
	dstDataMap, err := dir.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(dstDataMap) != len(srcDataMap) {
		t.Fatalf("expected dstDataMap length equals srcDataMap length, dstDataMap length: %d, srcDataMap length: %d", len(dstDataMap), len(srcDataMap))
	}
	for path, data := range dstDataMap {
		if string(data) != path {
			t.Fatalf("expected path equals to content, path: %s, content: %s", string(data), path)
		}
	}
}
