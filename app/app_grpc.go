package app

import (
	"OceanID/app/ocean_id"
	"net"

	"OceanID/app/impls"
	"OceanID/config"
	"OceanID/schemes/id_service"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	enable   bool
	ctx      context.Context
	listener net.Listener
	oceanID  oceanID.IDPool
	*grpc.Server
}

func (g *GrpcServer) GetName() string {
	return "grpc"
}

func (g *GrpcServer) Setup() error {

	args, err := config.AssertArgs(g.ctx)
	if err != nil {
		return errors.Wrap(err, "app.app_grpc")
	}
	serverAddr := args.Get("GRPC_SERVER_ADDR").(string)
	if serverAddr == "" {
		return nil
	}

	if g.listener, err = net.Listen("tcp", serverAddr); err != nil {
		return errors.Wrap(err, "app.app_grpc start grpc server failed")
	}

	g.Server = grpc.NewServer()

	idService.RegisterOceanIDServer(g.Server, oceanID.Mount[*impls.OceanIDGrpc](
		g.oceanID,
		&impls.OceanIDGrpc{},
	))

	g.enable = true

	return nil
}

func (g *GrpcServer) Start() error {
	if !g.enable {
		return nil
	}
	return g.Server.Serve(g.listener)
}

func (g *GrpcServer) Shutdown() error {
	if !g.enable {
		return nil
	}
	g.Stop()
	return nil
}

func NewGrpcApp(ctx context.Context, oi oceanID.IDPool) Application {
	return &GrpcServer{
		ctx:     ctx,
		oceanID: oi,
	}
}
