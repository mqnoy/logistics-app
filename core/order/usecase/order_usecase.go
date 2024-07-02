package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/mqnoy/logistics-app/core/constant"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/enum"
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	transaction "github.com/mqnoy/logistics-app/core/transaction_manager/repository"
	"github.com/mqnoy/logistics-app/core/util"
	"gorm.io/gorm"
)

type orderUseCase struct {
	txManager   transaction.TransactionManager
	orderRepo   domain.OrderRepository
	goodUseCase domain.GoodUseCase
}

func New(txManager transaction.TransactionManager, orderRepo domain.OrderRepository, goodUseCase domain.GoodUseCase) domain.OrderUseCase {
	return &orderUseCase{
		txManager:   txManager,
		orderRepo:   orderRepo,
		goodUseCase: goodUseCase,
	}
}

func (u *orderUseCase) OrderIn(ctx context.Context, param dto.CreateParam[dto.OrderInRequest]) (resp dto.OrderResponse, err error) {
	createValue := param.CreateValue

	// perform snapshot good
	goodSnapshot, goodRow, err := u.goodUseCase.SnapshotGood(createValue.Good.Code)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	// Persist insert order
	order := model.Order{
		RequestAt: util.GetCurrentTime(),
		Total:     createValue.Total,
		GoodSnapShotColumn: model.GoodSnapShotColumn{
			GoodSnapShot: goodSnapshot.Snapshot,
		},
		Type: int(enum.ORDER_IN),
	}

	// perform with transaction
	trx := u.txManager.AcquireTx(ctx)
	ctx = context.WithValue(ctx, constant.TrxKey, trx)

	// persist insert order
	orderRow, err := u.orderRepo.WithTrx(trx).InsertOrder(order)
	if err != nil {
		// rollback transaction
		u.txManager.CommitOrRollback(ctx, trx, true)

		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// call goodUseCase for increasing stock
	if err := u.goodUseCase.IncreaseStock(ctx, dto.UpdateParam[dto.GoodStockRequest]{
		ID: goodRow.ID,
		UpdateValue: dto.GoodStockRequest{
			Total: createValue.Total,
		},
	}); err != nil {
		// rollback transaction on IncreaseStock
		return resp, err
	}

	// commit transaction
	u.txManager.CommitOrRollback(ctx, trx, false)

	// compose order
	resp, err = u.ComposeOrder(orderRow)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return resp, nil
}

func (u *orderUseCase) ComposeOrder(m *model.Order) (resp dto.OrderResponse, err error) {
	goodSnapshotCol, err := m.GoodSnapShotColumn.ParseGoodSnapshot()
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	var goodResponse dto.GoodResponse
	if err := json.Unmarshal(goodSnapshotCol, &goodResponse); err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return dto.OrderResponse{
		ID:        m.ID,
		RequestAt: util.DateToEpoch(m.RequestAt),
		Type: dto.OrderTypeResponse{
			ID:   m.Type,
			Name: m.GetOrderTypeName(),
		},
		GoodSnapshotResponse: goodResponse,
		Total:                m.Total,
		Timestamp:            dto.ComposeTimestamp(m.TimestampColumn),
	}, nil
}

func (u *orderUseCase) OrderOut(ctx context.Context, param dto.CreateParam[dto.OrderInRequest]) (resp dto.OrderResponse, err error) {
	createValue := param.CreateValue

	// perform snapshot good
	goodSnapshot, goodRow, err := u.goodUseCase.SnapshotGood(createValue.Good.Code)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	// Persist insert order
	order := model.Order{
		RequestAt: util.GetCurrentTime(),
		Total:     createValue.Total,
		GoodSnapShotColumn: model.GoodSnapShotColumn{
			GoodSnapShot: goodSnapshot.Snapshot,
		},
		Type: int(enum.ORDER_OUT),
	}

	// perform with transaction
	trx := u.txManager.AcquireTx(ctx)
	ctx = context.WithValue(ctx, constant.TrxKey, trx)

	// persist insert order
	orderRow, err := u.orderRepo.WithTrx(trx).InsertOrder(order)
	if err != nil {
		// rollback transaction
		u.txManager.CommitOrRollback(ctx, trx, true)

		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// call goodUseCase for decreasing stock
	if err := u.goodUseCase.DecreaseStock(ctx, dto.UpdateParam[dto.GoodStockRequest]{
		ID: goodRow.ID,
		UpdateValue: dto.GoodStockRequest{
			Total: createValue.Total,
		},
	}); err != nil {
		// rollback transaction on DecreaseStock
		return resp, err
	}

	// commit transaction
	u.txManager.CommitOrRollback(ctx, trx, false)

	// compose order
	resp, err = u.ComposeOrder(orderRow)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return resp, nil
}

func (u *orderUseCase) ListOrders(param dto.ListParam[dto.FilterOrderParams]) (resp dto.ListResponse[dto.OrderResponse], err error) {
	pagination := param.Pagination
	param.Pagination.Offset = (pagination.Page - 1) * pagination.Limit

	rows, err := u.orderRepo.SelectAndCountOrder(param)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// Create pagination metadata
	totalItems := rows.Count
	totalPages := int(math.Ceil(float64(totalItems) / float64(pagination.Limit)))

	return dto.ListResponse[dto.OrderResponse]{
		Rows: u.ComposeListOrder(rows.Rows),
		MetaData: dto.Pagination{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
	}, nil
}

func (u *orderUseCase) ComposeListOrder(m []*model.Order) []dto.OrderResponse {
	results := make([]dto.OrderResponse, len(m))
	for idx, el := range m {
		// compose order
		resp, err := u.ComposeOrder(el)
		if err != nil {
			log.Println(err)
			continue
		}

		results[idx] = resp
	}

	return results
}

func (u *orderUseCase) DetailOrder(param dto.DetailParam) (resp dto.OrderResponse, err error) {
	orderRow, err := u.orderRepo.SelectOrderById(param.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, cerror.WrapError(http.StatusNotFound, fmt.Errorf("resource not found"))
		}

		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	// compose order
	resp, err = u.ComposeOrder(orderRow)
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	return resp, nil
}
