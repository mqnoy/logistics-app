package dto

import "net/http"

type GoodOrderResponse struct {
	Code string `json:"code"`
}

type OrderResponse struct {
	ID        string            `json:"id"`
	RequestAt int64             `json:"request_at"`
	Type      OrderTypeResponse `json:"type"`
	Good      GoodOrderResponse `json:"good"`
	Total     int               `json:"total"`
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
}
