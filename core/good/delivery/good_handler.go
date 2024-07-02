package delivery

import (
	"net/http"

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
