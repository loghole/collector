package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLSCertFile  string
	TLSKeyFile   string
}

type HTTP struct {
	config *Config
	server *http.Server
	router *mux.Router
}

func NewHTTP(config *Config) *HTTP {
	router := mux.NewRouter()

	server := &HTTP{
		config: config,
		router: router,
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      router,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
	}

	return server
}

func (h *HTTP) ListenAndServe() (err error) {
	switch {
	case h.config.TLSCertFile != "" && h.config.TLSKeyFile != "":
		err = h.server.ListenAndServeTLS(h.config.TLSCertFile, h.config.TLSKeyFile)
	default:
		err = h.server.ListenAndServe()
	}

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (h *HTTP) Router() *mux.Router {
	return h.router
}

func (h *HTTP) Addr() string {
	return fmt.Sprintf("http://%s", h.config.Addr)
}

func (h *HTTP) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
