package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrPathIsDirectory = errors.New("path is directory")
	ErrPathIsFile      = errors.New("path is file")
)

type PathStateType int

const (
	PathStateError PathStateType = iota
	PathStateEmpty
	PathStateFile
	PathStateDirectory
)

var pathStateName = map[PathStateType]string{
	PathStateError:     "error",
	PathStateEmpty:     "empty",
	PathStateFile:      "file",
	PathStateDirectory: "directory",
}

func (t PathStateType) String() string {
	if s, ok := pathStateName[t]; ok {
		return s
	}
	return fmt.Sprintf("pathStateName=%d?", int(t))
}

type Path struct {
	path string
}

func NewPath(path string) *Path {
	return &Path{
		path: path,
	}
}

func (p *Path) String() string {
	return p.path
}

func (p *Path) State() (PathStateType, error) {
	if info, err := os.Stat(p.path); err != nil {
		if os.IsNotExist(err) {
			return PathStateEmpty, nil
		} else {
			return PathStateError, err
		}
	} else if info.IsDir() {
		return PathStateDirectory, nil
	} else {
		return PathStateFile, nil
	}
}

func (p *Path) Remove() error {
	state, err := p.State()
	if err != nil {
		return err
	}
	if state == PathStateEmpty {
		return nil
	}
	return os.RemoveAll(p.path)
}

func (p *Path) File() (*File, error) {
	state, err := p.State()
	if err != nil {
		return nil, err
	}
	if state == PathStateDirectory {
		return nil, ErrPathIsDirectory
	}
	if state == PathStateEmpty {
		if err := os.MkdirAll(filepath.Dir(p.path), os.ModePerm); err != nil {
			return nil, err
		}
		f, err := os.Create(p.path)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return &File{
		path: p.path,
	}, nil
}

func (p *Path) Dir() (*Directory, error) {
	state, err := p.State()
	if err != nil {
		return nil, err
	}
	if state == PathStateFile {
		return nil, ErrPathIsFile
	}
	if state == PathStateEmpty {
		if err := os.MkdirAll(p.path, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &Directory{
		path: p.path,
	}, nil
}

type File struct {
	path string
}

func (f *File) String() string {
	return fmt.Sprintf("File[%s]", f.path)
}

func (f *File) Path() *Path {
	return NewPath(f.path)
}

func (f *File) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *File) Write(data []byte) error {
	return ioutil.WriteFile(f.path, data, os.ModePerm)
}

func (f *File) Append(data []byte) error {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

type Directory struct {
	path string
}

func (d *Directory) String() string {
	return fmt.Sprintf("Directory[%s]", d.path)
}

func (d *Directory) Path() *Path {
	return NewPath(d.path)
}

func (d *Directory) FilesList() {
	var paths []string
	filepath.Walk(d.path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths)
		}
		return nil
	})
}

func (d *Directory) Files() (map[string]*File, error) {
	fileMap := make(map[string]*File)
	filepath.Walk(d.path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f, err := NewPath(path).File()
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(d.path, path)
			if err != nil {
				return err
			}
			fileMap[relPath] = f
		}
		return nil
	})
	return fileMap, nil
}

func (d *Directory) Write(dataMap map[string][]byte) error {
	for dataPath, data := range dataMap {
		f, err := NewPath(filepath.Join(d.path, dataPath)).File()
		if err != nil {
			return err
		}
		if err := f.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func (d *Directory) Read() (map[string][]byte, error) {
	files, err := d.Files()
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string][]byte)
	for path, file := range files {
		data, err := file.Read()
		if err != nil {
			return nil, err
		}
		dataMap[path] = data
	}
	return dataMap, nil
}
