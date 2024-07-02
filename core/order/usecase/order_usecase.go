package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mqnoy/logistics-app/core/constant"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/enum"
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	transaction "github.com/mqnoy/logistics-app/core/transaction_manager/repository"
	"github.com/mqnoy/logistics-app/core/util"
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
		RequestAt: util.DateToEpoch(m.RequestAt),
		Type: dto.OrderTypeResponse{
			ID:   m.Type,
			Name: m.GetOrderTypeName(),
		},
		Good: dto.GoodOrderResponse{
			Code: goodResponse.Code,
		},
		Total:     m.Total,
		Timestamp: dto.ComposeTimestamp(m.TimestampColumn),
	}, nil
}
