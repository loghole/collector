package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gadavy/tracing"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/lissteron/loghole/dashboard/config"
	"github.com/lissteron/loghole/dashboard/internal/app/http/handlers"
	"github.com/lissteron/loghole/dashboard/internal/app/repositories/clickhouse"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases"
	"github.com/lissteron/loghole/dashboard/pkg/clickhouseclient"
	"github.com/lissteron/loghole/dashboard/pkg/log"
	"github.com/lissteron/loghole/dashboard/pkg/server"
)

// nolint: funlen,gocritic
func main() {
	// Init config, logger, exit chan
	config.Init()

	logger, err := log.NewLogger(config.LoggerConfig())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "init logger failed: %v", err)
		os.Exit(1)
	}

	logger.Infof("Version: %s, GitHash: %s, BuildAt: %s", config.Version, config.GitHash, config.BuildAt)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Init jaeger tracer.
	tracer, err := tracing.NewTracer(config.TracerConfig())
	if err != nil {
		logger.Fatalf("init tracing client failed: %v", err)
	}

	traceLogger := tracing.NewTraceLogger(logger)

	// Init clients
	clickhouseDB, err := clickhouseclient.NewClient(config.ClickhouseConfig())
	if err != nil {
		logger.Fatalf("init clickhouse db client failed: %v", err)
	}

	var (
		// Init repositories
		repository = clickhouse.NewRepository(clickhouseDB.Client(), traceLogger)

		// Init use cases
		entryList   = usecases.NewListEntry(repository, traceLogger)
		suggestList = usecases.NewListSuggest(repository, traceLogger)

		// Init http handlers
		listEntryHandlers   = handlers.NewEntryHandlers(entryList, traceLogger)
		listSuggestHandlers = handlers.NewSuggestHandlers(suggestList, traceLogger)
		filesHandlers       = handlers.NewFilesHandlers(viper.GetString("frontend.path"))
		infoHandlers        = handlers.NewInfoHandlers(traceLogger)

		// Init http middleware
		tracingMiddleware  = handlers.NewTracingMiddleware(tracer)
		compressMiddleware = handlers.NewCompressMiddleware(gzip.DefaultCompression, traceLogger)
	)

	// Init http server
	srv := server.NewHTTP(config.ServerConfig())

	// Init v1 routes
	r := srv.Router()
	r.Use(compressMiddleware.Middleware)
	r.HandleFunc("/", handlers.NewRedirectHandler("/ui/"))

	r1 := r.PathPrefix("/ui/")
	r1.Handler(http.StripPrefix("/ui/", filesHandlers.Handler())).Methods("GET")

	r2 := r.PathPrefix("/api/v1").Subrouter()
	r2.Use(tracingMiddleware.Middleware)
	r2.HandleFunc("/info", infoHandlers.InfoHandler).Methods("GET")
	r2.HandleFunc("/entry/list", listEntryHandlers.ListEntryHandler).Methods("POST")
	r2.HandleFunc("/suggest/{type}", listSuggestHandlers.ListHandler).Methods("POST")

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		logger.Infof("start http server on: %s", srv.Addr())
		return srv.ListenAndServe()
	})

	select {
	case <-exit:
		logger.Info("stopping application")
	case <-ctx.Done():
		logger.Error("stopping application with error")
	}

	if err = srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error while stopping web server: %v", err)
	}

	if err = errGroup.Wait(); err != nil {
		logger.Errorf("error while waiting for goroutines: %v", err)
	}

	if err = tracer.Close(); err != nil {
		logger.Errorf("error while stopping tracer: %v", err)
	}

	if err = clickhouseDB.Close(); err != nil {
		logger.Errorf("error while stopping clickhouse db: %v", err)
	}

	logger.Info("application stopped")
}
