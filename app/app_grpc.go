package app

import (
	"OceanID/app/impls"
	"OceanID/config"
	idService "OceanID/schemes/id_service"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	ctx      context.Context
	listener net.Listener
	*grpc.Server
}

func NewGrpcApp(ctx context.Context) Application {
	return &GrpcServer{
		ctx: ctx,
	}
}

func (g *GrpcServer) GetName() string {
	return "grpc"
}

func (g *GrpcServer) Setup() error {

	args, err := config.AssertArgs(g.ctx)
	if err != nil {
		return errors.Wrap(err, "app.app_grpc")
	}
	if g.listener, err = net.Listen("tcp", args.Get("SERVER_ADDR").(string)); err != nil {
		return errors.Wrap(err, "app.app_grpc start grpc server failed")
	}

	g.Server = grpc.NewServer()

	idService.RegisterOceanIDServer(g.Server, &impls.OceanID{})
	return nil
}

func (g *GrpcServer) Start() error {
	return g.Server.Serve(g.listener)
}

func (g *GrpcServer) Shutdown() error {
	g.Stop()
	return nil
}
