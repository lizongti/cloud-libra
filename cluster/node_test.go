package cluster_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/cluster"
	"github.com/aceaura/libra/cluster/component"
)

type TestService struct {
	component.ComponentBase
}

type TestHandlerRequest struct {
	Index int
}
type TestHandlerResponse struct {
	Index int
}

func (*TestService) TestHandler(_ context.Context, req *TestHandlerRequest) (resp *TestHandlerResponse, err error) {
	resp = &TestHandlerResponse{Index: req.Index}
	return
}

func TestBoot(t *testing.T) {
	cluster.NewNode().WithComponent(&TestService{}).Boot()
}
