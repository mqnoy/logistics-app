package dto

import "net/http"

type UserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Timestamp
}

type RegisterRequest struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"email,normalize"`
	Password string `json:"password" validate:"required"`
}

func (rr *RegisterRequest) Bind(r *http.Request) error {
	return nil
}
