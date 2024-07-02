package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/handler"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	"github.com/mqnoy/logistics-app/core/pkg/cvalidator"
)

type orderHandler struct {
	mux          *chi.Mux
	orderUseCase domain.OrderUseCase
}

func New(mux *chi.Mux, middlewareAuthorization domain.MiddlewareAuthorization, orderUseCase domain.OrderUseCase) {

	handler := orderHandler{
		mux:          mux,
		orderUseCase: orderUseCase,
	}

	mux.Route("/orders", func(r chi.Router) {
		r.Use(middlewareAuthorization.AuthorizationJWT)
		r.Post("/goods/in", handler.PostOrderIn)
		r.Post("/goods/out", handler.PostOrderOut)
		r.Get("/goods", handler.GetListOrders)

	})

}

func (h orderHandler) PostOrderIn(w http.ResponseWriter, r *http.Request) {
	var request dto.OrderInRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	param := dto.CreateParam[dto.OrderInRequest]{
		CreateValue: request,
		Session:     dto.GetAuthorizedUser(r.Context()),
	}

	// call usecase
	result, err := h.orderUseCase.OrderIn(r.Context(), param)

	handler.ParseResponse(w, r, "PostOrderIn", result, err)
}

func (h orderHandler) PostOrderOut(w http.ResponseWriter, r *http.Request) {
	var request dto.OrderInRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	param := dto.CreateParam[dto.OrderInRequest]{
		CreateValue: request,
		Session:     dto.GetAuthorizedUser(r.Context()),
	}

	// call usecase
	result, err := h.orderUseCase.OrderOut(r.Context(), param)

	handler.ParseResponse(w, r, "PostOrderOut", result, err)
}

func (h orderHandler) GetListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(handler.DefaultQuery(r, "page", "1"))
	limit, _ := strconv.Atoi(handler.DefaultQuery(r, "limit", "10"))
	offset, _ := strconv.Atoi(handler.DefaultQuery(r, "offset", "0"))
	orderType, _ := strconv.Atoi(handler.DefaultQuery(r, "type", "0"))
	requestAtRange := handler.ParseQueryToInt64Array(handler.DefaultQuery(r, "request_at_range", "[]"))
	orders := handler.DefaultQuery(r, "orders", "id desc")

	param := dto.ListParam[dto.FilterOrderParams]{
		Filters: dto.FilterOrderParams{
			RequestAt: requestAtRange,
			OrderType: orderType,
		},
		Orders: orders,
		Pagination: dto.Pagination{
			Page:   page,
			Limit:  limit,
			Offset: offset,
		},
		Session: dto.GetAuthorizedUser(r.Context()),
	}

	// Call usecase
	result, err := h.orderUseCase.ListOrders(param)

	handler.ParseResponse(w, r, "GetListOrders", result, err)
}
