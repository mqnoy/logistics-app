package dto

import "net/http"

type GoodResponse struct {
	ID                string            `json:"id"`
	Code              string            `json:"code"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	IsActive          bool              `json:"is_active"`
	GoodStockResponse GoodStockResponse `json:"stock"`
	Timestamp
}

type GoodCreateRequest struct {
	Code        string `json:"code" validate:"min=1,max=10,required"`
	Name        string `json:"name" validate:"min=1,max=20,required"`
	Description string `json:"description" validate:"required"`
}

func (t *GoodCreateRequest) Bind(r *http.Request) error {
	return nil
}

type GoodStockResponse struct {
	Total int `json:"total"`
}

type GoodUpdateRequest struct {
	Code        string `json:"code" validate:"min=1,max=10,required"`
	Name        string `json:"name" validate:"min=1,max=20,required"`
	Description string `json:"description" validate:"required"`
	IsActive    bool   `json:"is_active" validate:"boolean"`
}

func (t *GoodUpdateRequest) Bind(r *http.Request) error {
	return nil
}

type GoodSnapShot struct {
	ID                string            `json:"id"`
	Code              string            `json:"code"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	IsActive          bool              `json:"is_active"`
	GoodStockSnapshot GoodStockSnapshot `json:"stock"`
}

type GoodStockSnapshot struct {
	Total int `json:"total"`
}

type GoodStockRequest struct {
	Total int
}
