package main

import (
	_ "time/tzdata"

	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/common/pkg/logger"
)

func main() {
	if err := db.DB.AutoMigrate(
		&models.Product{},
		&models.User{},
	); err != nil {
		logger.Default.Fatal(err)
	} else {
		logger.Default.Printf("db_init: AutoMigrate success")
	}
}
