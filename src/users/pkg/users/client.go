package users

import (
	"source.local/common/pkg/env"
	"source.local/common/pkg/grpcbase"
	"source.local/common/pkg/logger"
	"source.local/users/internal/proto"
)

var (
	envs struct {
		Port int `env:"USERS_CONTAINER_PORT"`
	}

	Client proto.UsersServiceClient
)

func init() {
	if err := env.Parse(&envs); err != nil {
		logger.Default.Fatal(err)
	}
	if conn, err := grpcbase.Conn("users", envs.Port); err != nil {
		logger.Default.Fatal(err)
	} else {
		Client = proto.NewUsersServiceClient(conn)
	}
}
