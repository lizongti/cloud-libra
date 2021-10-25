package service

import (
	"reflect"

	"github.com/aceaura/libra/cluster/router"
	"github.com/aceaura/libra/codec"
)

type ServiceInterface interface {
	Bind(*router.Router)
}

type Service struct {
	codec      codec.Codec
	renameFunc func(reflect.Method) string
	codeFunc   func(reflect.Method) uint64
}

func (s *Service) Bind(router router.Router) {
}

func (s *Service) OnInit() {

}
