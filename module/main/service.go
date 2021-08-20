package main

import (
	"github.com/lizongti/libra/module/core/component"
	"github.com/lizongti/libra/module/core/service"
)

type Service struct {
	components []component.Component
}

var _ service.Service = (*Service)(nil)

func (s *Service) RegisterComponent(c component.Component) {
	s.components = append(s.components, c)
}
