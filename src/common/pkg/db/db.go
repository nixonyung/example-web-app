package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"source.local/common/pkg/env"
	"source.local/common/pkg/secret"
	"source.local/common/pkg/servicebase"
)

var (
	envs struct {
		ContainerName string `env:"POSTGRES_CONTAINER_NAME"`
		DatabaseName  string `env:"POSTGRES_DB"`
		User          string `env:"POSTGRES_USER"`
		Timezone      string `env:"TZ"`
	}
	secrets struct {
		Password string `secret:"postgres-password"`
	}

	DB *gorm.DB
)

func init() {
	if err := env.Parse(&envs); err != nil {
		servicebase.HandleErr(err)
	}
	if err := secret.Parse(&secrets); err != nil {
		servicebase.HandleErr(err)
	}
	if conn, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=5432 sslmode=disable TimeZone=%s",
				envs.ContainerName,
				envs.DatabaseName,
				envs.User,
				secrets.Password,
				envs.Timezone,
			),
		),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	); err != nil {
		log.Fatalln(err)
	} else {
		DB = conn
	}
}
