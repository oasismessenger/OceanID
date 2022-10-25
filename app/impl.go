package app

import (
	oceanID "OceanID/app/ocean_id"
	"context"

	"github.com/RealFax/pkg-ctl"
	"github.com/pkg/errors"
)

type Implement struct {
	ctx      context.Context
	services []Application
}

func NewService(ctx *context.Context) pkgCtl.Handler {
	return &Implement{ctx: *ctx}
}

func (i *Implement) mount(app Application) {
	i.services = append(i.services, app)
}

func (i *Implement) setupAll() (err error) {
	for _, service := range i.services {
		if err = service.Setup(); err != nil {
			return errors.Wrapf(err, "service %s setuped failed", service.GetName())
		}
	}
	return nil
}

func (i *Implement) startAll() (err error) {
	for _, service := range i.services {
		if err = service.Start(); err != nil {
			return errors.Wrapf(err, "service %s started failed", service.GetName())
		}
	}
	return nil
}

func (i *Implement) shutdownAll() (err error) {
	for _, service := range i.services {
		if err = service.Shutdown(); err != nil {
			return errors.Wrapf(err, "service %s shutdowned failed", service.GetName())
		}
	}
	return nil
}

func (i *Implement) Create() error {
	oi, err := oceanID.NewOceanID(i.ctx)
	if err != nil {
		return errors.Wrap(err, "app.impl")
	}
	// mount app
	i.mount(NewGrpcApp(i.ctx, oi))

	return i.setupAll()
}

func (i *Implement) Start() error {
	return i.startAll()
}

func (i *Implement) Destroy() error {
	return i.shutdownAll()
}

func (i *Implement) IsAsync() bool {
	return true
}
