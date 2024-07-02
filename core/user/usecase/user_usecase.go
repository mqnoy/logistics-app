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

func (u *userUseCase) LoginUser(payload dto.LoginRequest) (resp dto.LoginResponse, err error) {
	userRow, err := u.GetUserByEmail(payload.Email)
	if err != nil {
		return resp, err
	}

	// Compare password
	if err := u.ComparePassword(userRow.Password, payload.Password); err != nil {
		return resp, err
	}

	// Generate accessToken
	accessTknExpiry := jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.AccessTokenExpiry) * time.Second))
	accessTkn, err := u.GenerateToken(accessTknExpiry, userRow.ID)
	if err != nil {
		return resp, err
	}

	// Generate refreshToken
	refreshTknExpiry := jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.RefreshTokenExpiry) * time.Second))
	refreshTkn, err := u.GenerateToken(refreshTknExpiry, userRow.ID)
	if err != nil {
		return resp, err
	}

	return dto.LoginResponse{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
		UserResponse: u.ComposeUser(userRow),
	}, nil
}

func (u *userUseCase) GenerateToken(expiredIn *jwt.NumericDate, subjectId string) (string, error) {
	key := []byte(config.AppConfig.JWT.Key)
	mapClaims := token.GenerateMapClaims(token.CustomClaimOptions{
		ExpiredTime: expiredIn,
		SubjectId:   subjectId,
	})

	token, err := token.Generate(mapClaims, key)
	if err != nil {
		return "", cerror.WrapError(http.StatusInternalServerError, err)
	}

	return token, nil
}

func (u *userUseCase) ComparePassword(password string, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("password doesn't match"))
		}

		log.Println(err)
		return cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return nil
}
