package osutil

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

var (
	ErrPathIsDir  = errors.New("path is dir")
	ErrPathIsFile = errors.New("path is file")
)

type stat struct {
	info fs.FileInfo
	err  error
}

type statType string

var (
	StatDir       statType = "StatDir"
	StatFile      statType = "StatFile"
	StatNotExists statType = "StatNotExists"
	StatException statType = "StatException"
)

func Stat(path string) stat {
	info, err := os.Stat(path)
	return stat{
		info: info,
		err:  err,
	}
}

func StatInfo(path string) fs.FileInfo {
	return Stat(path).Info()
}

func (s stat) Info() fs.FileInfo {
	return s.info
}

func StatErr(path string) error {
	return Stat(path).Err()
}

func (s stat) Err() error {
	if s.err == nil {
		if s.info.IsDir() {
			return ErrPathIsDir
		} else {
			return ErrPathIsFile
		}
	}
	return s.err
}

func (s stat) Type() statType {
	if s.err != nil {
		if os.IsNotExist(s.err) {
			return StatNotExists
		} else {
			return StatException
		}
	} else if s.info.IsDir() {
		return StatDir
	}
	return StatFile
}

func (s stat) IsDir() bool {
	return s.Type() == StatDir
}

func (s stat) StatIsDir() error {
	if s.IsDir() {
		return nil
	}
	return s.Err()
}

func (s stat) MustDir() {
	if err := s.StatIsDir(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsNotDir() bool {
	return s.Type() != StatDir
}

func (s stat) StatIsNotDir() error {
	if s.IsNotDir() {
		return nil
	}
	return s.Err()
}

func (s stat) MustNotDir() {
	if err := s.StatIsNotDir(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsFile() bool {
	return s.Type() == StatFile
}

func (s stat) StatIsFile() error {
	if s.IsFile() {
		return nil
	}
	return s.Err()
}

func (s stat) MustFile() {
	if err := s.StatIsFile(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsNotFile() bool {
	return s.Type() != StatFile
}

func (s stat) StatIsNotFile() error {
	if s.IsNotFile() {
		return nil
	}
	return s.Err()
}

func (s stat) MustNotFile() {
	if err := s.StatIsNotFile(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsExist() bool {
	return s.Type() == StatFile || s.Type() == StatDir
}

func (s stat) StatIsExist() error {
	if s.IsExist() {
		return nil
	}
	return s.Err()
}

func (s stat) MustExist() {
	if err := s.StatIsExist(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsNotExist() bool {
	return s.Type() == StatNotExists
}

func (s stat) StatIsNotExist() error {
	if s.IsNotExist() {
		return nil
	}
	return s.Err()
}

func (s stat) MustNotExist() {
	if err := s.StatIsNotExist(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsException() bool {
	return s.Type() == StatException
}

func (s stat) StatIsException() error {
	if s.IsException() {
		return nil
	}
	return s.Err()
}

func (s stat) MustException() {
	if err := s.StatIsException(); err != nil {
		log.Panic(err)
	}
}

func (s stat) IsNotException() bool {
	return s.Type() != StatException
}

func (s stat) StatIsNotException() error {
	if s.IsNotException() {
		return nil
	}
	return s.Err()
}

func (s stat) MustNotExeption() {
	if err := s.StatIsNotException(); err != nil {
		log.Panic(err)
	}
}

func StatType(path string) statType {
	return Stat(path).Type()
}

func IsDir(path string) bool {
	return Stat(path).IsDir()
}

func StatIsDir(path string) error {
	return Stat(path).StatIsDir()
}

func MustDir(path string) {
	Stat(path).MustDir()
}

func IsNotDir(path string) bool {
	return Stat(path).IsNotDir()
}

func StatIsNotDir(path string) error {
	return Stat(path).StatIsNotDir()
}

func MustNotDir(path string) {
	Stat(path).MustNotDir()
}

func IsFile(path string) bool {
	return Stat(path).IsFile()
}

func StatIsFile(path string) error {
	return Stat(path).StatIsFile()
}

func MustFile(path string) {
	Stat(path).MustFile()
}

func IsNotFile(path string) bool {
	return Stat(path).IsNotFile()
}

func StatIsNotFile(path string) error {
	return Stat(path).StatIsNotFile()
}

func MustNotFile(path string) {
	Stat(path).MustNotFile()
}

func IsExist(path string) bool {
	return Stat(path).IsExist()
}

func StatIsExist(path string) error {
	return Stat(path).StatIsExist()
}

func MustExist(path string) {
	Stat(path).MustExist()
}

func IsNotExist(path string) bool {
	return Stat(path).IsNotExist()
}

func StatIsNotExist(path string) error {
	return Stat(path).StatIsNotExist()
}

func MustNotExist(path string) {
	Stat(path).MustNotExist()
}

func IsException(path string) bool {
	return Stat(path).IsException()
}

func StatIsException(path string) error {
	return Stat(path).StatIsException()
}

func MustException(path string) {
	Stat(path).MustException()
}

func IsNotException(path string) bool {
	return Stat(path).IsNotException()
}

func StatIsNotException(path string) error {
	return Stat(path).StatIsNotException()
}

func MustNotExeption(path string) {
	Stat(path).MustNotExeption()
}
