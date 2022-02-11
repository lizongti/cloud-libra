package filesystem

import (
	"context"
	"errors"

	"github.com/aceaura/libra/core/device"
)

var (
	ErrIsDir  = errors.New("path is a directory")
	ErrIsFile = errors.New("path is a file")
)

type ReadRequest struct {
	Path string
}

type ReadResponse struct {
	PathState     PathStateType
	FileData      []byte
	DirectoryData map[string][]byte
}

type WriteRequest struct {
	Path          string
	PathState     PathStateType
	PathRemove    bool
	FileTruncate  bool
	FileData      []byte
	DirectoryData map[string][]byte
}

type WriteResponse struct {
}

type Service struct{}

func init() {
	device.Bus().WithDevice(device.NewRouter().WithName("FileSystem").WithService(&Service{}))
}

func (s *Service) Read(ctx context.Context, req *ReadRequest) (resp *ReadResponse, err error) {
	resp = new(ReadResponse)
	path := NewPath(req.Path)
	state, err := path.State()
	if err != nil {
		return nil, err
	}
	resp.PathState = state
	switch resp.PathState {
	case PathStateEmpty:
		return resp, nil
	case PathStateFile:
		file, err := path.File()
		if err != nil {
			return nil, err
		}
		data, err := file.Read()
		if err != nil {
			return nil, err
		}
		resp.FileData = data
		return resp, nil
	case PathStateDirectory:
		directory, err := path.Directory()
		if err != nil {
			return nil, err
		}
		dataMap, err := directory.Read()
		if err != nil {
			return nil, err
		}
		resp.DirectoryData = dataMap
		return resp, nil
	}

	return resp, nil
}

func (s *Service) Write(ctx context.Context, req *WriteRequest) (resp *WriteResponse, err error) {
	resp = new(WriteResponse)
	path := NewPath(req.Path)
	state, err := path.State()
	if err != nil {
		return nil, err
	}
	if req.PathRemove || req.PathState != state {
		if err := path.Remove(); err != nil {
			return nil, err
		}
	}

	if req.PathState == PathStateFile {
		file, err := path.File()
		if err != nil {
			return nil, err
		}
		if req.FileTruncate {
			if err := file.Write(req.FileData); err != nil {
				return nil, err
			}
		} else {
			if err := file.Append(req.FileData); err != nil {
				return nil, err
			}
		}

		return resp, nil
	}

	if req.PathState == PathStateDirectory {
		directory, err := path.Directory()
		if err != nil {
			return nil, err
		}
		if req.FileTruncate {
			if err := directory.Write(req.DirectoryData); err != nil {
				return nil, err
			}
		} else {
			if err := directory.Append(req.DirectoryData); err != nil {
				return nil, err
			}
		}
		return resp, nil
	}

	return resp, nil
}
