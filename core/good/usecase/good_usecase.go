package usecase

import (
	"errors"
	"fmt"
	"log"
	"math"
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

func (u *goodUseCase) ListGoods(param dto.ListParam[dto.FilterCommonParams]) (resp dto.ListResponse[dto.GoodResponse], err error) {
	pagination := param.Pagination
	param.Pagination.Offset = (pagination.Page - 1) * pagination.Limit

	rows, err := u.goodRepo.SelectAndCountGood(param)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// Create pagination metadata
	totalItems := rows.Count
	totalPages := int(math.Ceil(float64(totalItems) / float64(pagination.Limit)))

	return dto.ListResponse[dto.GoodResponse]{
		Rows: u.ComposeListGood(rows.Rows),
		MetaData: dto.Pagination{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
	}, nil
}

func (u *goodUseCase) ComposeListGood(m []*model.Good) []dto.GoodResponse {
	results := make([]dto.GoodResponse, len(m))
	for idx, el := range m {
		results[idx] = u.ComposeGood(el)
	}

	return results
}

func (u *goodUseCase) UpdateGood(param dto.UpdateParam[dto.GoodUpdateRequest]) (resp dto.GoodResponse, err error) {
	row, err := u.DetailGoodById(param.ID)
	if err != nil {
		return resp, err
	}

	updateValue := param.UpdateValue

	// validate good is exist when code is not same with existing row
	if row.Code != updateValue.Code {
		goodExist, err := u.ValidateExistGood(updateValue.Code)
		if err != nil {
			return resp, err
		}

		if goodExist != nil {
			return resp, cerror.WrapError(http.StatusBadRequest, fmt.Errorf("duplicate resource"))
		}

	}

	// persist update data
	values := map[string]interface{}{
		"code":        updateValue.Code,
		"name":        updateValue.Name,
		"description": updateValue.Description,
		"is_active":   updateValue.IsActive,
	}
	if err := u.goodRepo.UpdateGoodById(row.ID, values); err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// persist select updated data
	return u.DetailGood(dto.DetailParam{
		ID: param.ID,
	})
}

func (u *goodUseCase) DeleteGood(param dto.DetailParam) error {
	row, err := u.DetailGoodById(param.ID)
	if err != nil {
		return err
	}

	// persist delete data
	if err := u.goodRepo.DeleteGoodById(row.ID); err != nil {
		log.Println(err)
		return cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return nil
}

func (u *goodUseCase) SnapshotGood(code string) (result dto.EntitySnapshot, row *model.Good, err error) {
	row, err = u.goodRepo.SelectGoodByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil, cerror.WrapError(http.StatusNotFound, fmt.Errorf("resource not found"))
		}

		log.Println(err)
		return result, nil, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// perform convert to json
	snapshot := dto.GoodSnapShot{
		ID:          row.ID,
		Code:        row.Code,
		Name:        row.Name,
		Description: row.Description,
		IsActive:    row.IsActive,
		GoodStockSnapshot: dto.GoodStockSnapshot{
			Total: row.GoodStock.Total,
		},
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		log.Println(err)
		return result, nil, err
	}

	return dto.EntitySnapshot{
		Snapshot: datatypes.JSON(snapshotJSON),
	}, row, nil
}
