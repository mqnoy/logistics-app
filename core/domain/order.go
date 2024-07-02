package domain

import (
	"context"

	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
)

type OrderUseCase interface {
	OrderIn(ctx context.Context, param dto.CreateParam[dto.OrderInRequest]) (resp dto.OrderResponse, err error)
}

type OrderRepository interface {
	WithTrx(trxHandle *gorm.DB) OrderRepository
	InsertOrder(row model.Order) (*model.Order, error)
}
