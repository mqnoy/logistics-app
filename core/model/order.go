package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mqnoy/logistics-app/core/enum"
	"gorm.io/gorm"
)

type Order struct {
	UUIDColumn
	RequestAt time.Time `gorm:"column:request_at"`
	Total     int       `gorm:"column:total"`
	Type      int       `gorm:"column:type"`
	GoodSnapShotColumn
	TimestampColumn
}

func (m Order) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("id", uuid)

	return nil
}

func (m *Order) ParseGoodSnapshot() (b []byte, err error) {
	b, err = m.GoodSnapShot.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (m Order) GetOrderTypeName() string {
	orderTypeMap := map[int]string{
		int(enum.ORDER_IN):  "ORDER_IN",
		int(enum.ORDER_OUT): "ORDER_OUT",
	}

	return orderTypeMap[m.Type]
}
