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

type ReadLocalFileRequest struct {
	Path string
}

type ReadLocalFileResponse struct {
	Data []byte
}

type ReadLocalDirectoryRequest struct {
	Path string
}

type ReadLocalDirectoryResponse struct {
	DataMap map[string][]byte
}

type Service struct{}

func init() {
	device.Bus().WithService(&Service{})
}

func (s *Service) ReadLocalFile(ctx context.Context, req *ReadLocalFileRequest) (resp *ReadLocalFileResponse, err error) {
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
	resp = new(ReadLocalFileResponse)
	resp.Data = data
	return resp, nil
}

func (s *Service) ReadLocalDir(ctx context.Context, req *ReadLocalDirectoryRequest) (resp *ReadLocalDirectoryResponse, err error) {
	if info, err := os.Stat(req.Path); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, ErrIsFile
	}
	return nil, nil
}
