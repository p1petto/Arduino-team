package server

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware function for authentication.
func (s *Server) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			op := "AuthMiddleware"
			reqToken := r.Header.Get("Authorization")
			if reqToken == "" {
				http.Error(w, "Incorrect token", http.StatusBadRequest)
				return
			}
			splitToken := strings.Split(reqToken, "Bearer")
			if len(splitToken) != 2 {
				http.Error(w, "Incorrect token", http.StatusBadRequest)
				return
			}

			reqToken = strings.TrimSpace(splitToken[1])
			s.log.Debug("got auth token", "token", reqToken)
			client, err := s.CheckToken(reqToken)
			if err != nil {
				s.log.Error(op+"bad token", slErr(err))
				http.Error(w, "Bad Token", http.StatusUnauthorized)
				return
			} else {
				ctx := context.WithValue(r.Context(), "user", client)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
