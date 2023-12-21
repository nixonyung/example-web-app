package models

import "source.local/common/internal/db/modelbase"

type User struct {
	modelbase.Entity
	Username     string `gorm:"unique;index"`
	PasswordHash string
	Role         UserRole
	Balance      float64
	Inventory    []Product `gorm:"many2many:users_products"`
}
