package osutil

import (
	"log"
	"os"
	"path/filepath"

	"github.com/zbiljic/go-filelock"
)

func ReadFile(path string) ([]byte, error) {
	stat := Stat(path)
	switch stat.Type() {
	case StatFile:
		return os.ReadFile(path)
	case StatNotExists:
		return nil, nil
	default:
		return nil, stat.Err()
	}
}

func WriteFile(path string, data []byte) error {
	if err := StatIsNotException(path); err != nil {
		return err
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}

type Lock struct {
	filelock.TryLockerSafe
}

func NewLock(path string) *Lock {
	tryLockerSafe, err := filelock.New(path)
	if err != nil {
		log.Panic(err)
	}
	return &Lock{
		TryLockerSafe: tryLockerSafe,
	}
}

func (l *Lock) ReadFile(path string) ([]byte, error) {
	if err := l.Lock(); err != nil {
		return nil, err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			log.Panic(err)
		}
	}()

	return ReadFile(path)
}

func (l *Lock) WriteFile(path string, data []byte) error {
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			log.Panic(err)
		}
	}()

	return WriteFile(path, data)
}
