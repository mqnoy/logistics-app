package dto

import "net/http"

type UserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Timestamp
}

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"email,normalize"`
	Password string `json:"password" validate:"required"`
}

func (rr *RegisterRequest) Bind(r *http.Request) error {
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,normalize"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginRequest) Bind(r *http.Request) error {
	return nil
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	UserResponse UserResponse `json:"user"`
}
