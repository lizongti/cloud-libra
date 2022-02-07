package file

import (
	"errors"
	"os"
)

var (
	ErrIsDir  = errors.New("path is a directory")
	ErrIsFile = errors.New("path is a file")
)

func Exist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func FileExist(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else if info.IsDir() {
		return false, ErrIsDir
	} else {
		return true, nil
	}
}

func DirExist(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else if info.IsDir() {
		return true, nil
	} else {
		return false, ErrIsFile
	}
}
