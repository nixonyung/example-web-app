package apigateway

import (
	"fmt"

	"source.local/common/pkg/env"
	"source.local/common/pkg/logger"
)

var (
	envs struct {
		Port int `env:"API_GATEWAY_CONTAINER_PORT"`
	}

	Host string
)

func init() {
	if err := env.Parse(&envs); err != nil {
		logger.Default.Fatal(err)
	}
	Host = fmt.Sprintf("api-gateway:%d", envs.Port)
}
