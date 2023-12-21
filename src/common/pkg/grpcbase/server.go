package grpcbase

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"source.local/common/pkg/logger"
)

var loggerUnary = func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	res, err := handler(ctx, req)
	if info.FullMethod != "/grpc.health.v1.Health/Check" {
		(&logger.Logger{NoLoc: true}).Printf("grpc_server [%s] \"%s\"",
			status.Code(err),
			info.FullMethod,
		)
	}
	return res, err
}

func NewServer(registerFn func(server *grpc.Server)) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggerUnary,
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(
				func(p any) error {
					return status.Errorf(codes.Unknown, "panic triggered: %v", p)
				},
			)),
		),
	)
	registerFn(server)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	reflection.Register(server)
	return server
}

func StartServer(server *grpc.Server, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("Serve: net.Listen: %w", err)
	}
	logger.Default.Printf("StartServer: server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("Serve: server.Serve: %w", err)
	}
	return nil
}
