package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Good struct {
	UUIDColumn
	Code        string `gorm:"column:code;type:varchar(10);not null;unique"`
	Name        string `gorm:"column:name;type:varchar(20);not null"`
	Description string `gorm:"column:description"`
	IsActive    bool   `gorm:"column:is_active;default:1"`
	TimestampColumn
}

func (m Good) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("id", uuid)

	return nil
}
