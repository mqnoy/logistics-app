package domain

import (
	"context"

	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
)

type GoodUseCase interface {
	CreateGood(param dto.CreateParam[dto.GoodCreateRequest]) (resp dto.GoodResponse, err error)
	DetailGood(param dto.DetailParam) (resp dto.GoodResponse, err error)
	ListGoods(param dto.ListParam[dto.FilterCommonParams]) (resp dto.ListResponse[dto.GoodResponse], err error)
	UpdateGood(param dto.UpdateParam[dto.GoodUpdateRequest]) (resp dto.GoodResponse, err error)
	DeleteGood(param dto.DetailParam) error
	SnapshotGood(code string) (result dto.EntitySnapshot, row *model.Good, err error)
	IncreaseStock(ctx context.Context, param dto.UpdateParam[dto.GoodStockRequest]) error
}

type GoodRepository interface {
	WithTrx(trxHandle *gorm.DB) GoodRepository

	InsertGood(row model.Good) (*model.Good, error)
	SelectGoodByCode(code string) (row *model.Good, err error)
	SelectGoodById(id string) (row *model.Good, err error)
	SelectAndCountGood(param dto.ListParam[dto.FilterCommonParams]) (dto.SelectAndCount[model.Good], error)
	UpdateGoodById(id string, values interface{}) error
	DeleteGoodById(id string) error

	InsertGoodStock(row model.GoodStock) (*model.GoodStock, error)
	UpdateGoodStockByGoodId(goodId string, values interface{}) error
}
