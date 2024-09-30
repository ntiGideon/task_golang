package middleware

import (
	"awesomeProject2/data/user"
	"awesomeProject2/helpers"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			helpers.WriteResponseBody(w, &user.WebResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization token is missing!",
				Data:    nil,
			}, http.StatusUnauthorized)
			return
		}

		claims, err := helpers.ValidateToken(bearerToken)

		if err != nil {
			helpers.WriteResponseBody(w, &user.WebResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			}, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.Id)
		next(w, r.WithContext(ctx), params)
	}
}
