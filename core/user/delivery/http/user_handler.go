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

type userHandler struct {
	mux         *chi.Mux
	userUseCase domain.UserUseCase
}

func New(mux *chi.Mux, userUseCase domain.UserUseCase) {
	handler := userHandler{
		mux:         mux,
		userUseCase: userUseCase,
	}

	mux.Route("/users", func(r chi.Router) {
		r.Post("/register", handler.Register)
	})
}

func (h userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	if err := render.Bind(r, &request); err != nil {
		handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusBadRequest, err))
		return
	}

	// Validate payload
	if err := cvalidator.ValidateStruct(&request); err != nil {
		handler.ParseToErrorValidation(w, r, http.StatusBadRequest, cvalidator.ErrorValidator, err)
		return
	}

	result, err := h.userUseCase.RegisterUser(request)

	handler.ParseResponse(w, r, "Register", result, err)
}
}
