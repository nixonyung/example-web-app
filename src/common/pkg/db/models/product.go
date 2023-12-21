package models

import "source.local/common/internal/db/modelbase"

type Product struct {
	modelbase.Entity
	Name         string `gorm:"index"`
	Price        float32
	StockTotal   int
	StockCurrent int
}
