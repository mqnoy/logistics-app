package domain

import (
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
)

type UserUseCase interface {
	RegisterUser(request dto.RegisterRequest) (resp dto.UserResponse, err error)
}

type UserRepository interface {
	InsertUser(row model.User) (*model.User, error)
	SelectUserByEmail(email string) (*model.User, error)
}
