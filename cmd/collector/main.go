package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/loghole/database"
	"github.com/loghole/lhw/zap"
	"github.com/loghole/lhw/zaplog"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracehttp"
	"github.com/loghole/tracing/tracelog"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/loghole/collector/config"
	entryV1 "github.com/loghole/collector/internal/app/api/entry/v1"
	"github.com/loghole/collector/internal/app/api/middleware"
	splunkV1 "github.com/loghole/collector/internal/app/api/splunk/v1"
	"github.com/loghole/collector/internal/app/repositories/clickhouse"
	"github.com/loghole/collector/internal/app/services/entry"
	"github.com/loghole/collector/pkg/server"
)

const _defaultRetryTry = 10

// nolint: funlen,gocritic
func main() {
	// Init config, logger, exit chan
	config.Init()

	logger, err := zaplog.NewLogger(config.LoggerConfig(), zaplog.AddCaller())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "init logger failed: %v", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.With(
		"InstanceUUID", config.InstanceUUID,
		"Version", config.Version,
		"GitHash", config.GitHash,
		"BuildAt", config.BuildAt,
		"ServiceName", config.ServiceName,
		"AppName", config.AppName,
	).Info("application init")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Init jaeger tracer.
	tracer, err := tracing.NewTracer(config.TracerConfig())
	if err != nil {
		logger.Fatalf("init tracing client failed: %v", err)
	}

	traceLogger := tracelog.NewTraceLogger(logger.SugaredLogger)

	// Init clients
	clickhouseDB, err := database.New(
		config.ClickhouseConfig(),
		database.WithReconnectHook(),
		clockhouseRetryFunc(logger),
	)
	if err != nil {
		logger.Fatalf("can't connect to clickhouse db: %v", err)
	}

	// Init repository
	repository := clickhouse.NewEntryRepository(
		clickhouseDB,
		traceLogger,
		viper.GetInt("service.writer.capacity"),
		viper.GetDuration("service.writer.period"),
	)

	// Init service
	entryService := entry.NewService(repository, traceLogger)

	// Init handlers
	var (
		entryHandlers = entryV1.NewEntryHandlers(entryService, traceLogger, tracer)
		splunkHandler = splunkV1.NewSplunkHandler(entryService, traceLogger, tracer)
		infoHandlers  = entryV1.NewInfoHandlers(traceLogger)

		remoteIPMiddleware = middleware.NewRemoteIPMiddleware("service.ip.header")
		authMiddleware     = middleware.NewAuthMiddleware(
			viper.GetBool("service.auth.enable"),
			viper.GetStringSlice("service.auth.tokens"),
		)
	)

	srv := server.NewHTTP(config.ServerConfig())

	r := srv.Router()
	r.HandleFunc("/api/v1/info", infoHandlers.InfoHandler)

	r1 := r.PathPrefix("/api/v1").Subrouter()
	r1.Use(authMiddleware.Middleware, remoteIPMiddleware.Middleware, tracehttp.NewMiddleware(tracer).Middleware)
	r1.HandleFunc("/store", entryHandlers.StoreItemHandler)
	r1.HandleFunc("/store/list", entryHandlers.StoreListHandler)
	r1.HandleFunc("/ping", entryHandlers.PingHandler)

	r2 := r.PathPrefix("/services/collector/event").Subrouter()
	r2.Use(authMiddleware.Middleware, remoteIPMiddleware.Middleware, tracehttp.NewMiddleware(tracer).Middleware)
	r2.HandleFunc("/1.0", splunkHandler.Handler)

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		logger.Info("start entry writer")

		return repository.Run(ctx)
	})

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

	repository.Stop()

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

func clockhouseRetryFunc(logger *zap.Logger) database.Option {
	return database.WithRetryFunc(func(retryCount int, err error) bool {
		if err == nil {
			return false
		}

		logger.Infof("retry count: %d", retryCount)

		if retryCount > _defaultRetryTry {
			return false
		}

		return strings.Contains(err.Error(), "broken pipe") ||
			strings.Contains(err.Error(), "bad connection") ||
			strings.Contains(err.Error(), "connection timed out")
	})
}
