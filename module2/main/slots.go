package main

import "github.com/lizongti/libra/module/core/context"

type Slots struct {
	Router
}

func (s *Slots) String() string {
	return "slots"
}

func (*Slots) Login(c context.Context, _ []byte) ([]byte, error) {
	return nil, nil
}
