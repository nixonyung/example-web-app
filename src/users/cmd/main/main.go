package main

import (
	"google.golang.org/grpc"
	"source.local/common/pkg/env"
	"source.local/common/pkg/grpcbase"
	"source.local/common/pkg/logger"
	"source.local/users/internal/proto"
	"source.local/users/internal/server"
)

var (
	envs struct {
		Port int `env:"USERS_CONTAINER_PORT"`
	}
)

func init() {
	if err := env.Parse(&envs); err != nil {
		logger.Default.Fatal(err)
	}
}

func main() {
	server := grpcbase.NewServer(func(s *grpc.Server) {
		proto.RegisterUsersServiceServer(s, &server.Server{})
	})
	if err := grpcbase.StartServer(server, envs.Port); err != nil {
		logger.Default.Fatal(err)
	}
}
