package main

import "github.com/lizongti/libra/module/core/context"

type Context struct {
}

var _ context.Context = (*Context)(nil)
