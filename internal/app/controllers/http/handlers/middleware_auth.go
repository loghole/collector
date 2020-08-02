package handlers

import (
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

type AuthMiddleware struct {
	enabled bool
	tokens  map[string]struct{}
}

func NewAuthMiddleware(enabled bool, tokens []string) *AuthMiddleware {
	middleware := &AuthMiddleware{
		enabled: enabled,
		tokens:  make(map[string]struct{}, len(tokens)),
	}

	for _, token := range tokens {
		middleware.tokens[strings.TrimSpace(token)] = struct{}{}
	}

	return middleware
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	if !m.enabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimSpace(r.Header.Get(authorizationHeader))

		if auth == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parts := strings.Split(auth, " ")

		if len(parts) < 2 || !strings.EqualFold(parts[0], "bearer") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if _, ok := m.tokens[parts[1]]; !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// if we authenticated successfully, go ahead and remove the bearer token so that no one
		// is ever tempted to use it inside of the API server
		r.Header.Del(authorizationHeader)

		next.ServeHTTP(w, r)
	})
}
