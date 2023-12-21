package main

import (
	"log"
	_ "time/tzdata"

	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/common/pkg/servicebase"
)

func main() {
	if err := db.DB.AutoMigrate(
		&models.Product{},
		&models.User{},
	); err != nil {
		servicebase.HandleErr(err)
	} else {
		log.Println("AutoMigrate success")
	}
}
