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
