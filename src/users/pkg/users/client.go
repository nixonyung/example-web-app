package users

import (
	"log"

	"source.local/common/pkg/env"
	"source.local/common/pkg/grpcbase"
	"source.local/common/pkg/servicebase"
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
		servicebase.HandleErr(err)
	}
	if conn, err := grpcbase.Conn("users", envs.Port); err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		Client = proto.NewUsersServiceClient(conn)
	}
}
