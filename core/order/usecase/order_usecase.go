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

	var goodResponse *dto.GoodResponse
	if goodSnapshotCol != nil && len(goodSnapshotCol) != 2 {
		if err := json.Unmarshal(goodSnapshotCol, &goodResponse); err != nil {
			log.Println(err)
		}
	}

	// compose orderItems
	var items []dto.OrderItemResponse
	if len(m.OrderItem) != 0 {
		items = u.ComposeListOrderItem(m.OrderItem)
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
		CountItem:            m.CountItem,
		Items:                items,
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

func (u *orderUseCase) MultipleOrderInOut(ctx context.Context, param dto.CreateParam[dto.OrderCreateMultipleRequest]) (resp dto.OrderCreateMultipleResponse, err error) {
	createValue := param.CreateValue

	// perform with transaction
	trx := u.txManager.AcquireTx(ctx)
	ctx = context.WithValue(ctx, constant.TrxKey, trx)

	var successOrderItems []dto.OrderItemTemp
	var failedOrderItems []dto.OrderItemTemp

	// persist insert order
	orderRow, err := u.orderRepo.WithTrx(trx).InsertOrder(model.Order{
		RequestAt: util.GetCurrentTime(),
		Type:      int(createValue.Type),
	})
	if err != nil {
		// rollback transaction
		u.txManager.CommitOrRollback(ctx, trx, true)

		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	var orderItemRows []*model.OrderItem
	orderId := orderRow.UUIDColumn.ID
	for _, orderItem := range createValue.Items {
		// determine goods by code
		goodSnapshot, goodRow, err := u.goodUseCase.SnapshotGood(orderItem.Code)
		if err != nil {
			log.Println(err)

			failedOrderItems = append(failedOrderItems, dto.OrderItemTemp{
				Code:   orderItem.Code,
				GoodID: "",
				Total:  orderItem.Total,
				Reason: err,
			})

			continue
		}

		// check barang should active
		if !goodRow.IsActive {
			failedOrderItems = append(failedOrderItems, dto.OrderItemTemp{
				Code:   orderItem.Code,
				GoodID: goodRow.ID,
				Total:  orderItem.Total,
				Reason: fmt.Errorf("resource inactive"),
			})
			continue
		}

		// check stock if order is type out
		if createValue.Type == enum.ORDER_OUT {
			if err := u.goodUseCase.CheckAvailabilityStock(goodRow.GoodStock.Total, orderItem.Total); err != nil {
				failedOrderItems = append(failedOrderItems, dto.OrderItemTemp{
					Code:   orderItem.Code,
					GoodID: goodRow.ID,
					Total:  orderItem.Total,
					Reason: fmt.Errorf("out of stock"),
				})
				continue
			}
		}

		goodsId := goodRow.ID

		// persist insert orderItem
		itemRow, err := u.orderRepo.WithTrx(trx).InsertOrderItem(model.OrderItem{
			OrderID: orderId,
			Total:   orderItem.Total,
			GoodID:  goodRow.ID,
			GoodSnapShotColumn: model.GoodSnapShotColumn{
				GoodSnapShot: goodSnapshot.Snapshot,
			},
		})
		if err != nil {
			log.Println(err)

			failedOrderItems = append(failedOrderItems, dto.OrderItemTemp{
				Code:   orderItem.Code,
				GoodID: "",
				Total:  orderItem.Total,
				Reason: err,
			})
			continue
		}
		orderItemRows = append(orderItemRows, itemRow)

		// append success order
		successOrderItems = append(successOrderItems, dto.OrderItemTemp{
			Code:   orderItem.Code,
			GoodID: goodsId,
			Total:  orderItem.Total,
			Reason: nil,
		})
	}

	// call goodUseCase for increasing/decrease stock
	if createValue.Type == enum.ORDER_IN {
		if err := u.OrderTypeInFn(ctx, successOrderItems); err != nil {
			return resp, err
		}
	} else if createValue.Type == enum.ORDER_OUT {
		if err := u.OrderTypeOutFn(ctx, successOrderItems); err != nil {
			return resp, err
		}
	} else {
		// rollback transaction
		u.txManager.CommitOrRollback(ctx, trx, true)

		return resp, cerror.WrapError(http.StatusBadRequest, fmt.Errorf("error determine order type"))
	}

	// update count_item
	u.orderRepo.WithTrx(trx).UpdateOrderById(orderRow.ID, map[string]interface{}{
		"count_item": len(successOrderItems),
	})

	// commit transaction
	if len(successOrderItems) != 0 {
		u.txManager.CommitOrRollback(ctx, trx, false)
	} else {
		u.txManager.CommitOrRollback(ctx, trx, true)
	}

	// compose response
	orderRow.OrderItem = orderItemRows
	return u.ComposeOrderCreateMultiple(orderRow, successOrderItems, failedOrderItems), nil
}

func (u *orderUseCase) OrderTypeInFn(ctx context.Context, successOrderItems []dto.OrderItemTemp) (err error) {
	for _, orderItem := range successOrderItems {
		if err := u.goodUseCase.IncreaseStock(ctx, dto.UpdateParam[dto.GoodStockRequest]{
			ID: orderItem.GoodID,
			UpdateValue: dto.GoodStockRequest{
				Total: orderItem.Total,
			},
		}); err != nil {
			// rollback transaction on IncreaseStock
			return err
		}
	}

	return nil
}

func (u *orderUseCase) OrderTypeOutFn(ctx context.Context, successOrderItems []dto.OrderItemTemp) (err error) {
	for _, orderItem := range successOrderItems {
		if err := u.goodUseCase.DecreaseStockV2(ctx, dto.UpdateParam[dto.GoodStockRequest]{
			ID: orderItem.GoodID,
			UpdateValue: dto.GoodStockRequest{
				Total: orderItem.Total,
			},
		}); err != nil {
			// rollback transaction on IncreaseStock
			return err
		}
	}

	return nil
}

func (u *orderUseCase) ComposeOrderCreateMultiple(m *model.Order, success, failed []dto.OrderItemTemp) dto.OrderCreateMultipleResponse {
	var successArr = make([]dto.OrderItemMultipleResponse, len(success))
	var reason string
	for i, s := range success {
		if s.Reason != nil {
			reason = s.Reason.Error()
		}
		successArr[i] = dto.OrderItemMultipleResponse{
			Code:   s.Code,
			Total:  s.Total,
			Reason: reason,
		}
	}

	var failedArr = make([]dto.OrderItemMultipleResponse, len(failed))
	for i, s := range failed {
		if s.Reason != nil {
			reason = s.Reason.Error()
		}

		failedArr[i] = dto.OrderItemMultipleResponse{
			Code:   s.Code,
			Total:  s.Total,
			Reason: reason,
		}
	}

	return dto.OrderCreateMultipleResponse{
		ID:        m.ID,
		Success:   successArr,
		Failed:    failedArr,
		Timestamp: dto.ComposeTimestamp(m.TimestampColumn),
	}
}

func (u *orderUseCase) ComposeOrderItem(m *model.OrderItem) (resp dto.OrderItemResponse, err error) {
	goodSnapshotCol, err := m.GoodSnapShotColumn.ParseGoodSnapshot()
	if err != nil {
		log.Println(err)
		return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}

	var goodResponse dto.GoodResponse
	if goodSnapshotCol != nil {
		if err := json.Unmarshal(goodSnapshotCol, &goodResponse); err != nil {
			log.Println(err)
			return resp, cerror.WrapError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		}
	}

	return dto.OrderItemResponse{
		ID:                   m.ID,
		GoodSnapshotResponse: goodResponse,
		Total:                m.Total,
		Timestamp:            dto.ComposeTimestamp(m.TimestampColumn),
	}, nil
}

func (u *orderUseCase) ComposeListOrderItem(m []*model.OrderItem) []dto.OrderItemResponse {
	results := make([]dto.OrderItemResponse, len(m))
	for idx, el := range m {
		// compose order
		resp, err := u.ComposeOrderItem(el)
		if err != nil {
			log.Println(err)
			continue
		}

		results[idx] = resp
	}

	return results
}
