package middleware

import (
	dto "backend/dto/result"
	jwtToken "backend/pkg/jwt"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type Result struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")

		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response := dto.ErrorResult{Status: http.StatusUnauthorized, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := Result{Status: http.StatusUnauthorized, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		ctx := context.WithValue(r.Context(), "authInfo", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
