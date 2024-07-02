package usecase

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mqnoy/logistics-app/core/config"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	"github.com/mqnoy/logistics-app/core/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func New(userRepo domain.UserRepository) domain.UserUseCase {

	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) RegisterUser(request dto.RegisterRequest) (resp dto.UserResponse, err error) {
	// Validate email is exist
	existEmail, err := u.GetUserByEmail(request.Email)
	if existEmail != nil && err == nil {
		return resp, cerror.WrapError(http.StatusBadRequest, fmt.Errorf("email already exist"))
	}

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	row := model.User{
		FullName: request.FullName,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	user, err := u.userRepo.InsertUser(row)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return u.ComposeUser(user), nil
}

func (u *userUseCase) ComposeUser(m *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        m.ID,
		FullName:  m.FullName,
		Email:     m.Email,
		Timestamp: dto.ComposeTimestamp(m.TimestampColumn),
	}
}

func (u *userUseCase) GetUserByEmail(email string) (*model.User, error) {
	row, err := u.userRepo.SelectUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("resource not found")
		}

		log.Println(err)
		return nil, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return row, nil
}
