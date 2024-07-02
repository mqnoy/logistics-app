package domain

import (
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
)

type GoodUseCase interface {
	CreateTask(param dto.CreateParam[dto.GoodCreateRequest]) (resp dto.GoodResponse, err error)
}

type GoodRepository interface {
	InsertGood(row model.Good) (*model.Good, error)
	SelectGoodByCode(code string) (row *model.Good, err error)

	InsertGoodStock(row model.GoodStock) (*model.GoodStock, error)
}
