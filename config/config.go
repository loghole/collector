package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/gadavy/tracing"
	"github.com/spf13/viper"

	"github.com/lissteron/loghole/collector/pkg/clickhouseclient"
	"github.com/lissteron/loghole/collector/pkg/log"
	"github.com/lissteron/loghole/collector/pkg/server"
)

const (
	defaultServiceName = "collector"
)

// nolint:gochecknoglobals // build args
var (
	ServiceName string
	AppName     string
	GitHash     string
	Version     string
	BuildAt     string
)

func Init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs/")

	_ = viper.ReadInConfig()

	viper.SetDefault("server.read.timeout", time.Minute)
	viper.SetDefault("server.write.timeout", time.Minute)
	viper.SetDefault("server.idle.timeout", time.Minute*10) // nolint:gomnd,gocritic

	viper.SetDefault("clickhouse.read.timeout", 10)
	viper.SetDefault("clickhouse.write.timeout", 20)
	viper.SetDefault("service.writer.capacity", 1000)
	viper.SetDefault("service.writer.period", time.Second)
}

func ClickhouseConfig() *clickhouseclient.Config {
	return &clickhouseclient.Config{
		Addr:         viper.GetString("clickhouse.uri"),
		User:         viper.GetString("clickhouse.user"),
		Password:     viper.GetString("clickhouse.password"),
		Database:     viper.GetString("clickhouse.database"),
		ReadTimeout:  viper.GetInt("clickhouse.read.timeout"),
		WriteTimeout: viper.GetInt("clickhouse.write.timeout"),
	}
}

func TracerConfig() *tracing.Config {
	var service string

	switch {
	case ServiceName != "":
		service = ServiceName
	case viper.GetString("service.name") != "":
		service = viper.GetString("service.name")
	default:
		service = defaultServiceName
	}

	return &tracing.Config{
		URI:         viper.GetString("jaeger.uri"),
		Enabled:     viper.GetString("jaeger.uri") != "",
		ServiceName: service,
	}
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

func LoggerConfig() *log.Config {
	return &log.Config{
		Level: viper.GetString("logger.level"),
		Options: []log.Option{
			log.AddCaller(),
		},
	}
}
