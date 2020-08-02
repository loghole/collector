package handlers

import (
	"net"
	"net/http"
	"strings"
)

type RemoteIPMiddleware struct {
	header string
}

func NewRemoteIPMiddleware(header string) *RemoteIPMiddleware {
	return &RemoteIPMiddleware{header: strings.TrimSpace(header)}
}

func (m *RemoteIPMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var remoteAddr string

		if m.header != "" {
			remoteAddr = r.Header.Get(m.header)
		}

		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}

		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			r.RemoteAddr = host
		}

		r.Header.Del(m.header)

		next.ServeHTTP(w, r)
	})
}
