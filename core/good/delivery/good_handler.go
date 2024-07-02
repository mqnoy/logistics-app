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

type goodHandler struct {
	mux         *chi.Mux
	goodUseCase domain.GoodUseCase
}

func New(mux *chi.Mux, goodUseCase domain.GoodUseCase) {

	handler := goodHandler{
		mux:         mux,
		goodUseCase: goodUseCase,
	}

	mux.Route("/goods", func(r chi.Router) {
		r.Post("/", handler.PostCreateGood)
		r.Get("/{id}", handler.GetDetailGood)
		r.Get("/", handler.GetListGoods)
		r.Put("/{id}", handler.PutUpdateGood)
	})

}

func (h goodHandler) PostCreateGood(w http.ResponseWriter, r *http.Request) {
	var request dto.GoodCreateRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	param := dto.CreateParam[dto.GoodCreateRequest]{
		CreateValue: request,
	}

	// call usecase
	result, err := h.goodUseCase.CreateGood(param)

	handler.ParseResponse(w, r, "PostCreateGood", result, err)
}

func (h goodHandler) GetDetailGood(w http.ResponseWriter, r *http.Request) {
	param := dto.DetailParam{
		ID: chi.URLParam(r, "id"),
	}

	// Call usecase
	result, err := h.goodUseCase.DetailGood(param)

	handler.ParseResponse(w, r, "GetDetailGood", result, err)
}

func (h goodHandler) GetListGoods(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(handler.DefaultQuery(r, "page", "1"))
	limit, _ := strconv.Atoi(handler.DefaultQuery(r, "limit", "10"))
	offset, _ := strconv.Atoi(handler.DefaultQuery(r, "offset", "0"))
	keyword, _ := handler.GetQuery(r, "keyword")

	qIsActive, _ := handler.GetQuery(r, "is_active")
	IsActive := handler.ParseQueryToBool(qIsActive)

	orders := handler.DefaultQuery(r, "orders", "id desc")

	param := dto.ListParam[dto.FilterCommonParams]{
		Filters: dto.FilterCommonParams{
			Keyword:  keyword,
			IsActive: IsActive,
		},
		Orders: orders,
		Pagination: dto.Pagination{
			Page:   page,
			Limit:  limit,
			Offset: offset,
		},
	}

	// Call usecase
	result, err := h.goodUseCase.ListGoods(param)

	handler.ParseResponse(w, r, "GetListGoods", result, err)
}

func (h goodHandler) PutUpdateGood(w http.ResponseWriter, r *http.Request) {
	var request dto.GoodUpdateRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", err, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// Validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	param := dto.UpdateParam[dto.GoodUpdateRequest]{
		UpdateValue: request,
		ID:          chi.URLParam(r, "id"),
	}

	// Call usecase
	result, err := h.goodUseCase.UpdateGood(param)

	handler.ParseResponse(w, r, "PutUpdateGood", result, err)
}
