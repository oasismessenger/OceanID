package config

import (
	"context"
	pkgCtl "github.com/RealFax/pkg-ctl"
)

type Implement struct {
	ctx *context.Context
}

func NewLoader(ctx *context.Context) pkgCtl.Handler {
	return &Implement{ctx: ctx}
}

func (i *Implement) Create() error {
	parseArgs(i.ctx)
	return nil
}

func (i *Implement) Start() error {
	return nil
}

func (i *Implement) Destroy() error {
	return nil
}

func (i *Implement) IsAsync() bool {
	return false
}
