package apigateway

import (
	"fmt"

	"source.local/common/pkg/env"
	"source.local/common/pkg/servicebase"
)

var (
	envs struct {
		Port int `env:"APIGATEWAY_CONTAINER_PORT"`
	}
)

func init() {
	if err := env.Parse(&envs); err != nil {
		servicebase.HandleErr(err)
	}
}

func URL(path string) string {
	return fmt.Sprintf("http://api-gateway:%d%s", envs.Port, path)
}
