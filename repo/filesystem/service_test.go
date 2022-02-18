package filesystem_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/aceaura/libra/boost/magic"
	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/repo/filesystem"
	"github.com/mitchellh/go-homedir"
)

func TestService(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}

	var (
		ctx    = context.Background()
		e      = encoding.NewChainEncoding(magic.UnixChain("json.base64.lazy"), magic.UnixChain("lazy.base64.json"))
		rWrite = route.NewChainRoute(magic.GoogleChain("/client"), magic.GoogleChain("/file-system/write"))
		rRead  = route.NewChainRoute(magic.GoogleChain("/client"), magic.GoogleChain("/file-system/read"))
	)

	path1 := filepath.Join("dir_A", "dir_B", "file_C")
	path2 := filepath.Join("dir_A", "file_B")
	srcDataMap := map[string][]byte{
		path1: []byte(path1),
		path2: []byte(path2),
	}

	client := device.NewClient().WithName("Client")
	service := &filesystem.Service{}
	fileSystemRouter := device.NewRouter().WithName("FileSystem").WithService(service)
	bus := device.NewRouter().WithBus().WithName("Bus").WithDevice(fileSystemRouter).WithDevice(client)
	t.Logf("\n%s", device.Tree(bus))

	reqWrite := &filesystem.WriteRequest{
		Path:          filepath.Join(home, ".libra", "service_test"),
		PathState:     filesystem.PathStateDirectory,
		PathRemove:    true,
		FileTruncate:  true,
		DirectoryData: srcDataMap,
	}
	data, err := e.Marshal(reqWrite)
	if err != nil {
		t.Fatal(err)
	}
	msg := &message.Message{
		Route:    rWrite,
		Encoding: e,
		Data:     data,
	}
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(filesystem.WriteResponse)
		if err := msg.Encoding.Unmarshal(msg.Data, resp); err != nil {
			return err
		}
		return nil
	})
	if err = client.Invoke(ctx, msg, processor); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}

	reqRead := &filesystem.ReadRequest{
		Path: filepath.Join(home, ".libra", "service_test"),
	}
	data, err = e.Marshal(reqRead)
	if err != nil {
		t.Fatal(err)
	}
	msg = &message.Message{
		Route:    rRead,
		Encoding: e,
		Data:     data,
	}
	processor = device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(filesystem.ReadResponse)
		if err := msg.Encoding.Unmarshal(msg.Data, resp); err != nil {
			return err
		}
		if resp.PathState != filesystem.PathStateDirectory {
			t.Fatal("expected path state directory")
		}
		dstDataMap := resp.DirectoryData
		if len(dstDataMap) != len(srcDataMap) {
			t.Fatalf("expected dstDataMap length equals srcDataMap length, dstDataMap length: %d, srcDataMap length: %d", len(dstDataMap), len(srcDataMap))
		}
		for path, data := range dstDataMap {
			if string(data) != path {
				t.Fatalf("expected path equals to content, path: %s, content: %s", string(data), path)
			}
		}
		return nil
	})
	if err = client.Invoke(ctx, msg, processor); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}
}
