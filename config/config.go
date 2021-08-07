package config

import (
	"fmt"
	"strings"
	"time"

	// init clickhouse driver.
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/google/uuid"
	"github.com/loghole/database"
	"github.com/loghole/lhw/zap"
	"github.com/loghole/tracing"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"

	"github.com/loghole/collector/pkg/server"
)

const (
	_defaultServiceName = "collector"

	_defaultServerIdleTimeout = time.Minute * 10

	_defaultClickhouseReadTimeoutSeconds  = 10
	_defaultClickhouseWriteTimeoutSeconds = 20

	_defaultServerWriterCapacity = 1000
)

// nolint:gochecknoglobals // build args
var (
	InstanceUUID = uuid.New()
	ServiceName  string
	AppName      string
	GitHash      string
	Version      string
	BuildAt      string
)

func Init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("json")
	viper.SetConfigName(_defaultServiceName)
	viper.AddConfigPath("./configs/")

	_ = viper.ReadInConfig()

	viper.SetDefault("server.read.timeout", time.Minute)
	viper.SetDefault("server.write.timeout", time.Minute)
	viper.SetDefault("server.idle.timeout", _defaultServerIdleTimeout)

	viper.SetDefault("clickhouse.read.timeout", _defaultClickhouseReadTimeoutSeconds)
	viper.SetDefault("clickhouse.write.timeout", _defaultClickhouseWriteTimeoutSeconds)
	viper.SetDefault("service.writer.capacity", _defaultServerWriterCapacity)
	viper.SetDefault("service.writer.period", time.Second)
}

func ClickhouseConfig() *database.Config {
	return &database.Config{
		Addr:         viper.GetString("clickhouse.uri"),
		User:         viper.GetString("clickhouse.user"),
		Database:     viper.GetString("clickhouse.database"),
		ReadTimeout:  viper.GetString("clickhouse.read.timeout"),
		WriteTimeout: viper.GetString("clickhouse.write.timeout"),
		Type:         database.ClickhouseDatabase,
	}
}

func TracerConfig() *config.Configuration {
	return tracing.DefaultConfiguration(serviceName(), viper.GetString("jaeger.uri"))
}

func ServerConfig() *server.Config {
	return &server.Config{
		Addr:         fmt.Sprintf("0.0.0.0:%s", viper.GetString("server.http.port")),
		ReadTimeout:  viper.GetDuration("server.read.timeout"),
		WriteTimeout: viper.GetDuration("server.write.timeout"),
		IdleTimeout:  viper.GetDuration("server.idle.timeout"),
		TLSCertFile:  viper.GetString("server.tls.cert"),
		TLSKeyFile:   viper.GetString("server.tls.key"),
	}
}

func LoggerConfig() *zap.Config {
	return &zap.Config{
		Level:         viper.GetString("logger.level"),
		CollectorURL:  viper.GetString("logger.collector.url"),
		Namespace:     viper.GetString("logger.namespace"),
		Source:        serviceName(),
		BuildCommit:   GitHash,
		DisableStdout: false,
	}
}

func serviceName() string {
	switch {
	case ServiceName != "":
		return ServiceName
	case viper.GetString("service.name") != "":
		return viper.GetString("service.name")
	default:
		return _defaultServiceName
	}
}
