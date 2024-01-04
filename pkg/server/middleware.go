package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/senyc/jason/pkg/auth"
	"github.com/senyc/jason/pkg/types"
)

func (s *Server) autorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		userId, err := s.db.GetUserIdFromApiKey(auth.EncryptApiKey(key))

		if err == nil {
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			s.logger.Panic(err)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (s *Server) jwtAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			decodedJwt *jwt.Token
			err        error
		)
		bearerToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearerToken, "Bearer") {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			token := strings.TrimPrefix(bearerToken, "Bearer")
			token = strings.TrimSpace(token)
			decodedJwt, err = jwt.ParseWithClaims(token, &types.JwtClaims{}, func(tok *jwt.Token) (any, error) {
				privateKey, err := auth.GetJwtPrivateKey()
				return &privateKey.PublicKey, err
			})
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				s.logger.Panic(fmt.Errorf("jwt auth failure %v", err))
			}
			if claims, ok := decodedJwt.Claims.(*types.JwtClaims); ok && decodedJwt.Valid {
				ctx := context.WithValue(r.Context(), "userId", claims.Uuid)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
