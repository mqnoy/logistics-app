package usecase

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	"gorm.io/gorm"
)

type goodUseCase struct {
	goodRepo domain.GoodRepository
}

func New(goodRepo domain.GoodRepository) domain.GoodUseCase {
	return &goodUseCase{
		goodRepo: goodRepo,
	}
}

func (u *goodUseCase) CreateGood(param dto.CreateParam[dto.GoodCreateRequest]) (resp dto.GoodResponse, err error) {

	createValue := param.CreateValue

	// validate good is exist on database
	goodExist, err := u.ValidateExistGood(createValue.Code)
	if err != nil {
		return resp, err
	}

	if goodExist != nil {
		return resp, cerror.WrapError(http.StatusBadRequest, fmt.Errorf("duplicate resource"))
	}

	// Persist insert goods
	good := model.Good{
		Code:        createValue.Code,
		Name:        createValue.Name,
		Description: createValue.Description,
	}
	goodRow, err := u.goodRepo.InsertGood(good)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// insert default stock
	stock := model.GoodStock{
		Total:  0,
		GoodID: goodRow.ID,
	}
	if _, err := u.goodRepo.InsertGoodStock(stock); err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return u.ComposeGood(goodRow), nil
}

func (u *goodUseCase) ValidateExistGood(code string) (row *model.Good, err error) {
	row, err = u.goodRepo.SelectGoodByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Println(err)
		return nil, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return row, nil
}

func (u *goodUseCase) ComposeGood(m *model.Good) dto.GoodResponse {
	return dto.GoodResponse{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		GoodStockResponse: dto.GoodStockResponse{
			Total: m.GoodStock.Total,
		},
		Timestamp: dto.ComposeTimestamp(m.TimestampColumn),
	}
}

func (u *goodUseCase) DetailGoodById(id string) (row *model.Good, err error) {
	row, err = u.goodRepo.SelectGoodById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerror.WrapError(http.StatusNotFound, fmt.Errorf("resource not found"))
		}

		log.Println(err)
		return nil, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return row, nil
}

func (u *goodUseCase) DetailGood(param dto.DetailParam) (resp dto.GoodResponse, err error) {
	row, err := u.DetailGoodById(param.ID)
	if err != nil {
		return resp, err
	}

	return u.ComposeGood(row), nil
}
