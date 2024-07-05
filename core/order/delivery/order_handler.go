package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/enum"
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
		r.Get("/{id}", handler.GetDetailOrder)
		r.Post("/multiple/goods/in", handler.PostOrderMultipleIn)
		r.Post("/multiple/goods/out", handler.PostOrderMultipleOut)
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
	goodsId := handler.DefaultQuery(r, "goodId", "")

	param := dto.ListParam[dto.FilterOrderParams]{
		Filters: dto.FilterOrderParams{
			RequestAt: requestAtRange,
			OrderType: orderType,
			GoodId:    goodsId,
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

func (h orderHandler) GetDetailOrder(w http.ResponseWriter, r *http.Request) {
	param := dto.DetailParam{
		ID:      chi.URLParam(r, "id"),
		Session: dto.GetAuthorizedUser(r.Context()),
	}

	// Call usecase
	result, err := h.orderUseCase.DetailOrder(param)

	handler.ParseResponse(w, r, "GetDetailOrder", result, err)
}

func (h orderHandler) PostOrderMultipleIn(w http.ResponseWriter, r *http.Request) {
	var request dto.OrderCreateMultipleRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	// perform with order type in
	request.Type = enum.ORDER_IN
	param := dto.CreateParam[dto.OrderCreateMultipleRequest]{
		CreateValue: request,
		Session:     dto.GetAuthorizedUser(r.Context()),
	}

	// call usecase
	result, err := h.orderUseCase.MultipleOrderInOut(r.Context(), param)

	handler.ParseResponse(w, r, "PostOrderMultipleIn", result, err)
}

func (h orderHandler) PostOrderMultipleOut(w http.ResponseWriter, r *http.Request) {
	var request dto.OrderCreateMultipleRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	// perform with order type in
	request.Type = enum.ORDER_OUT
	param := dto.CreateParam[dto.OrderCreateMultipleRequest]{
		CreateValue: request,
		Session:     dto.GetAuthorizedUser(r.Context()),
	}

	// call usecase
	result, err := h.orderUseCase.MultipleOrderInOut(r.Context(), param)

	handler.ParseResponse(w, r, "PostOrderMultipleOut", result, err)
}
