package middleware

import (
	"net/http"
	"strings"
)

const (
	_tokenParts          = 2
	_authorizationHeader = "Authorization"
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
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		auth := strings.TrimSpace(r.Header.Get(_authorizationHeader))

		if auth == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		parts := strings.Split(auth, " ")

		if len(parts) < _tokenParts {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		if _, ok := m.tokens[parts[1]]; !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		// if we authenticated successfully, go ahead and remove the bearer token so that no one
		// is ever tempted to use it inside of the API server
		r.Header.Del(_authorizationHeader)

		next.ServeHTTP(w, r)
	})
}
