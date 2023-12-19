package server

import "net/http"

func (s *Server) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		userId, err := s.db.GetUserIdFromApiKey(key)
		if err == nil {
			r.Header.Del("id")
			r.Header.Add("id", userId)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

	})
}
