package dto

import (
	"net/http"

	"github.com/mqnoy/logistics-app/core/enum"
)

type GoodOrderResponse struct {
	Code string `json:"code"`
}

type OrderResponse struct {
	ID                   string              `json:"id"`
	RequestAt            int64               `json:"request_at"`
	Type                 OrderTypeResponse   `json:"type"`
	GoodSnapshotResponse *GoodResponse       `json:"good_snapshot"`
	Total                int                 `json:"total"`
	CountItem            int                 `json:"count_item"`
	Items                []OrderItemResponse `json:"items"`
	Timestamp
}

type OrderTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GoodOrderRequest struct {
	Code string `json:"code" validate:"min=1,max=10,required"`
}

type OrderInRequest struct {
	Good  GoodOrderRequest `json:"good"`
	Total int              `json:"total" validate:"min=1,number,required"`
}

func (t *OrderInRequest) Bind(r *http.Request) error {
	return nil
}

type FilterOrderParams struct {
	RequestAt []int64
	OrderType int
	GoodId    string
}

type OrderCreateMultipleRequest struct {
	Items []OrderItem `json:"items" validate:"required,unique"`
	Type  enum.OrderType
}

// Bind implements render.Binder.
func (o *OrderCreateMultipleRequest) Bind(r *http.Request) error {
	return nil
}

type OrderItem struct {
	Code  string `json:"code" validate:"min=1,max=10,required"`
	Total int    `json:"total" validate:"min=1,number,required"`
}

type OrderItemTemp struct {
	Code   string
	GoodID string
	Total  int
	Reason error
}

type OrderItemResponse struct {
	ID                   string       `json:"id"`
	GoodSnapshotResponse GoodResponse `json:"good_snapshot"`
	Total                int          `json:"total"`
	Timestamp
}

type OrderItemMultipleResponse struct {
	Code   string `json:"code"`
	Total  int    `json:"total"`
	Reason string `json:"reason"`
}

type OrderCreateMultipleResponse struct {
	ID      string                      `json:"id"`
	Success []OrderItemMultipleResponse `json:"success"`
	Failed  []OrderItemMultipleResponse `json:"failed"`
	Timestamp
}
