package file

import (
	"context"
	"path/filepath"

	"github.com/aceaura/libra/boost/magic"
	"github.com/aceaura/libra/boost/tree"
	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
)

type Document struct {
	Path string
	Data []byte
}

type Folder struct {
	Path string
	Data map[string][]byte
}

type Tree tree.Tree

type Mapping struct {
	*device.Client
	opts mappingOptions
}

func NewMapping(opt ...ApplyMappingOption) *Mapping {
	opts := defaultMappingOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Mapping{
		Client: device.NewClient(""),
		opts:   opts,
	}
}

func (m *Mapping) String() string {
	return m.opts.name
}

// func (m *Mapping) Replace(t *tree.Tree) {

// }

// func (m *Mapping) Select(v interface{}) (_ interface{}, err error) {
// 	reqWrite := &WriteRequest{
// 		Path:          m.opts.path,
// 		PathState:     PathStateDirectory,
// 		PathRemove:    true,
// 		FileTruncate:  true,
// 		DirectoryData: srcDataMap,
// 	}
// 	data, err := e.Marshal(reqWrite)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func (m *Mapping) selectDocument(document *Document) (*Document, error) {
// 	m.Invoke()
// }

// func (m *Mapping) SelectFolder() {

// }

func (m *Mapping) invokeRead(path string, cmd []string) (result []Document, err error) {
	if err := m.Client.Invoke(m.opts.ctx, &message.Message{
		Route:    route.NewChainRoute(device.Addr(m), magic.GoogleChain("/file-system/read")),
		Encoding: encoding.NewJSON(),
		Data: encoding.Encode(encoding.NewJSON(), &ReadRequest{
			Path: filepath.Join(m.opts.path, path),
		}),
	}, device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(ReadResponse)
		encoding.Decode(msg.Encoding, msg.Data, resp)
		switch resp.PathState {
		case PathStateDirectory:
			for path, data := range resp.DirectoryData {
				result = append(result, Document{
					Path: path,
					Data: data,
				})
			}
		case PathStateFile:
			result = append(result, Document{
				Data: resp.FileData,
			})
		}
		return nil
	})); err != nil {
		return nil, err
	}
	return result, nil
}

type mappingOptions struct {
	path string
	name string
	ctx  context.Context
}

var defaultMappingOptions = mappingOptions{
	ctx: context.Background(),
}

type ApplyMappingOption interface {
	apply(*mappingOptions)
}

type funcMappingOption func(*mappingOptions)

func (f funcMappingOption) apply(opt *mappingOptions) {
	f(opt)
}

func WithMappingPath(path string) funcMappingOption {
	return func(c *mappingOptions) {
		c.path = path
	}
}

func WithMappingName(name string) funcMappingOption {
	return func(c *mappingOptions) {
		c.name = name
	}
}

func WithMappingContext(ctx context.Context) funcMappingOption {
	return func(c *mappingOptions) {
		c.ctx = ctx
	}
}
