package main

import (
	"fmt"
	_ "time/tzdata"

	"source.local/common/pkg/env"
	"source.local/common/pkg/httpbase"
	"source.local/common/pkg/logger"
	"source.local/users/pkg/users"
)

var (
	envs struct {
		Port int `env:"API_GATEWAY_CONTAINER_PORT"`
	}
)

func init() {
	if err := env.Parse(&envs); err != nil {
		logger.Default.Fatal(err)
	}
}

func main() {
	server := httpbase.NewServer()
	{
		server.Post("/users", users.CreateUserHTTP)
		server.Get("/users", users.FindUsersHTTP)
		server.Post("/users/add-balance", users.AddBalance)
		server.Delete("/users", users.DeleteUser)
	}
	if err := server.Listen(fmt.Sprintf(":%d", envs.Port)); err != nil {
		logger.Default.Fatal(err)
	}
}
