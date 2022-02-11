package file

import (
	"context"
	"errors"
	"io/ioutil"
	"os"

	"github.com/aceaura/libra/core/device"
)

var (
	ErrIsDir  = errors.New("path is a directory")
	ErrIsFile = errors.New("path is a file")
)

type ReadFileRequest struct {
	Path string
}

type ReadFileResponse struct {
	Data []byte
}

type WriteFileRequest struct {
	Truncate bool
	Path     string
	Data     []byte
}

type WriteFileResponse struct {
}

type ReadDirectoryRequest struct {
	Path string
}

type ReadLocalDirectoryResponse struct {
	DataMap map[string][]byte
}

type Service struct{}

func init() {
	device.Bus().WithService(&Service{})
}

func (s *Service) ReadFile(ctx context.Context, req *ReadFileRequest) (resp *ReadFileResponse, err error) {
	if info, err := os.Stat(req.Path); err != nil {
		return nil, err
	} else if info.IsDir() {
		return nil, ErrIsDir
	}
	f, err := os.Open(req.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	resp = new(ReadFileResponse)
	resp.Data = data
	return resp, nil
}

func (s *Service) WriteFile(ctx context.Context, req *WriteFileRequest) (resp *WriteFileResponse, err error) {
	// if info, err := os.Stat(req.Path); err != nil {
	// 	if os.IsNotExist(err) {
	// 		// os.Create(req.Path)
	// 		os.MkdirAll(req.Path)
	// 	}
	// }
	return nil, nil
}

// func (s *Service) DeleteLocalFile()

func (s *Service) ReadLocalDir(ctx context.Context, req *ReadDirectoryRequest) (resp *ReadLocalDirectoryResponse, err error) {
	if info, err := os.Stat(req.Path); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, ErrIsFile
	}
	return nil, nil
}
