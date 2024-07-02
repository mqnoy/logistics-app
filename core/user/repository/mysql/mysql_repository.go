package mysql

import (
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) domain.UserRepository {
	return &mysqlUserRepository{
		DB: db,
	}
}

func (m mysqlUserRepository) InsertUser(row model.User) (*model.User, error) {
	err := m.DB.Create(&row).Error
	return &row, err
}

func (m mysqlUserRepository) SelectUserByEmail(email string) (row *model.User, err error) {
	if err := m.DB.First(&row, "email=?", email).Error; err != nil {
		return nil, err
	}
	return row, nil
}
