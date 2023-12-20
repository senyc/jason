package server

import (
	"context"
	"net/http"

	"github.com/senyc/jason/pkg/auth"
)

func (s *Server) autorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		userId, err := s.db.GetUserIdFromApiKey(auth.EncryptApiKey(key))

		if err == nil {
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

	})
}
