package rbac

import (
	"net/http"

	"github.com/joabegranvile/auth-go/internal/auth"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Middleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(auth.UserCtxKey).(*auth.Claims)

		if claims.Role != requiredRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
