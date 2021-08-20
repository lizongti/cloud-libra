package node

import (
	"github.com/lizongti/libra/router"
	"github.com/lizongti/libra/service"
)

type Node interface {
	WithRouter(router.Router)
	RegisterService(service.Service)
}
