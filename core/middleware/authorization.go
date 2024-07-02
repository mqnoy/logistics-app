package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mqnoy/logistics-app/core/config"
	"github.com/mqnoy/logistics-app/core/constant"
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/handler"
	"github.com/mqnoy/logistics-app/core/pkg/cerror"
	"github.com/mqnoy/logistics-app/core/pkg/token"
	"github.com/mqnoy/logistics-app/core/util"
)

type authorizationMiddleware struct {
	userUseCase domain.UserUseCase
}

func NewAuthorizationMiddleware(userUseCase domain.UserUseCase) domain.MiddlewareAuthorization {
	return &authorizationMiddleware{
		userUseCase: userUseCase,
	}
}

func (am *authorizationMiddleware) AuthorizationJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract authorization from header
		tokenString, err := util.ExtractTokenBearerFromHeader(r)
		if err != nil {
			handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusUnauthorized, err))
			return
		}

		// Verify the token jwt
		mapClaims := make(jwt.MapClaims)
		key := []byte(config.AppConfig.JWT.Key)
		verifyToken, err := token.Verify(mapClaims, key, tokenString)

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusUnauthorized, err))
				return
			}

			handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusUnauthorized, err))
			return
		}

		if !verifyToken.Valid {
			handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusUnauthorized, jwt.ErrTokenNotValidYet))
			return
		}

		claims, ok := verifyToken.Claims.(jwt.MapClaims)
		if ok {
			subject, err := claims.GetSubject()
			if err != nil {
				handler.ParseResponse(w, r, "", nil, cerror.WrapError(http.StatusUnauthorized, err))
				return
			}

			// TODO: Validate role and privileges

			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.AuthorizedUserCtxKey, dto.AuthorizedUser{
				UserID: subject,
			})

			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		handler.ParseToErrorMsg(w, r, http.StatusForbidden, cerror.WrapError(http.StatusUnauthorized, fmt.Errorf("you don't have permission to access resource")))

		h.ServeHTTP(w, r)
	})
}
