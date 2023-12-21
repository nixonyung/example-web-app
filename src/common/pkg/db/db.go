package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"source.local/common/pkg/env"
	"source.local/common/pkg/logger"
	"source.local/common/pkg/secret"
)

var (
	envs struct {
		DatabaseName string `env:"POSTGRES_DB"`
		User         string `env:"POSTGRES_USER"`
		Timezone     string `env:"TZ"`
	}
	secrets struct {
		Password string `secret:"postgres-password"`
	}

	DB *gorm.DB
)

func init() {
	if err := env.Parse(&envs); err != nil {
		logger.Default.Fatal(err)
	}
	if err := secret.Parse(&secrets); err != nil {
		logger.Default.Fatal(err)
	}
	if conn, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=postgres dbname=%s user=%s password=%s port=5432 sslmode=disable TimeZone=%s",
				envs.DatabaseName,
				envs.User,
				secrets.Password,
				envs.Timezone,
			),
		),
		&gorm.Config{
			Logger: &dbLogger{LogLevel: gorm_logger.Info},
		},
	); err != nil {
		logger.Default.Fatal(err)
	} else {
		DB = conn
	}
}
