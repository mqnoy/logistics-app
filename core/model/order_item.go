package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	UUIDColumn
	OrderID string `gorm:"column:order_id"`
	Order   Order  `gorm:"foreignKey:OrderID"`
	Total   int    `gorm:"column:total"`
	GoodID  string `gorm:"column:good_id"`
	GoodSnapShotColumn
	TimestampColumn
}

func (m OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("id", uuid)

	return nil
}
