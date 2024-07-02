package mysql

import (
	"log"

	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
)

type mysqlOrderRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.OrderRepository {
	return &mysqlOrderRepository{
		db: db,
	}
}

func (m mysqlOrderRepository) WithTrx(trxHandle *gorm.DB) domain.OrderRepository {
	if trxHandle == nil {
		log.Println("transaction not found")
		return m
	}
	m.db = trxHandle

	return m
}

func (m mysqlOrderRepository) InsertOrder(row model.Order) (*model.Order, error) {
	err := m.db.Create(&row).Error
	return &row, err
}
