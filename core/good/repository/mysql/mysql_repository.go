package mysql

import (
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
)

type mysqlGoodRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.GoodRepository {
	return &mysqlGoodRepository{
		db: db,
	}
}

func (m mysqlGoodRepository) InsertGood(row model.Good) (*model.Good, error) {
	err := m.db.Create(&row).Error
	return &row, err
}

func (m mysqlGoodRepository) SelectGoodByCode(code string) (row *model.Good, err error) {
	if err := m.db.
		Where("code=?", code).First(&row).
		Error; err != nil {
		return nil, err
	}

	return row, nil
}

func (m mysqlGoodRepository) InsertGoodStock(row model.GoodStock) (*model.GoodStock, error) {
	err := m.db.Create(&row).Error
	return &row, err
}
