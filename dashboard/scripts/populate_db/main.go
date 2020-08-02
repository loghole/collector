package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gadavy/lhw"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	viper.AutomaticEnv()

	time.Sleep(viper.GetDuration("SLEEP"))

	// init writer
	writer, err := lhw.NewWriter(
		lhw.NodeWithAuth(viper.GetString("COLLECTOR_URI"), viper.GetString("COLLECTOR_AUTH")),
		lhw.WithInsecure(), lhw.WithLogger(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		log.Fatalln("init writer failed", err)
	}

	defer writer.Close() // flushes storage, if contain any data

	// init logger
	logger := zap.New(zapcore.NewCore(
		getEncoder(),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer)),
		zapcore.DebugLevel)).Sugar()

	generator := NewErrorGenerator()

	log.Println("start write logs")

	for i := 0; i < viper.GetInt("COUNT"); i++ {
		entry, err := generator.GenerateEntry()
		if err != nil {
			log.Println("GenerateEntry", err)
			continue
		}

		l := logger.With("namespace", entry.Namespace).
			With("source", entry.Source).
			With("host", entry.Host).
			With("trace_id", entry.TraceID).
			With("build_commit", entry.BuildCommit).
			With("config_hash", entry.ConfigHash)

		if entry.JSON != nil {
			l = l.With("json_key", entry.JSON)
		}

		switch entry.Level {
		case "debug":
			l.Debug(entry.Message)
		case "info":
			l.Info(entry.Message)
		case "warn":
			l.Warn(entry.Message)
		case "error":
			l.Error(entry.Message)
		}

		time.Sleep(time.Millisecond * 5)
	}

	log.Println("success")
}

type ErrorGenerator struct {
	template     []string
	errorTxt     []string
	loggerParams []params
	rnd          *rand.Rand
}

func NewErrorGenerator() *ErrorGenerator {
	return &ErrorGenerator{
		template:     errTemplates,
		errorTxt:     errText,
		loggerParams: loggerParams,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *ErrorGenerator) randomMessage() string {
	return fmt.Sprintf(e.template[e.rnd.Intn(len(e.template)-1)], e.errorTxt[e.rnd.Intn(len(e.errorTxt)-1)])
}

type Entry struct {
	Namespace   string
	Source      string
	Host        string
	Level       string
	TraceID     string
	Message     string
	Params      string
	BuildCommit string
	ConfigHash  string
	JSON        interface{}
}

func (e *ErrorGenerator) GenerateEntry() (*Entry, error) {
	param := e.randomParams()

	entry := &Entry{
		Namespace:   param.namespace,
		Source:      param.source,
		Host:        param.host,
		Level:       e.randomLevel(),
		TraceID:     strconv.FormatInt(time.Now().UnixNano(), 32) + strconv.FormatInt(rand.Int63n(1000), 32),
		Message:     e.randomMessage(),
		Params:      "",
		BuildCommit: param.buildCommit,
		ConfigHash:  param.configHash,
		JSON:        e.randomJSON(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	entry.Params = string(data)

	return entry, nil
}

func (e *ErrorGenerator) randomParams() params {
	return e.loggerParams[e.rnd.Intn(len(e.loggerParams)-1)]
}

func (e *ErrorGenerator) randomJSON() interface{} {
	if e.rnd.Uint64()&(1<<63) != 0 {
		return jsonRand[e.rnd.Intn(len(jsonRand)-1)]
	}

	return nil
}

func (e *ErrorGenerator) randomLevel() string {
	switch e.rnd.Intn(4) {
	case 0:
		return "debug"
	case 1:
		return "info"
	case 2:
		return "warn"
	case 3:
		return "error"
	default:
		return "info"
	}
}

// getEncoder return log hole json encoding
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func RFC3339NanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339Nano))
}

type params struct {
	buildCommit string
	configHash  string
	host        string
	source      string
	namespace   string
}

var loggerParams = []params{
	{
		buildCommit: "db957a22b3c1d6e508c0828917a5e14c572fb007",
		configHash:  "130362a6fd10cf2f939dd0cfc0ab222cee6a99ec",
		host:        "127.100.0.1:50000",
		source:      "app_1",
		namespace:   "prod",
	},
	{
		buildCommit: "7fca4555eb6f5eb20ce7f81900595b9c4ee5f47e",
		configHash:  "f365c5c9e9217df5cd3271d3cd25a5ebbf94cfeb",
		host:        "127.11.10.1:45673",
		source:      "app_1",
		namespace:   "dev",
	},
	{
		buildCommit: "0448b96036e3683623a9d6712fd6f6f891607b90",
		configHash:  "4de34c86ce4234270d487c4dd36fdfefe4566bd9",
		host:        "127.11.10.2:45655",
		source:      "app_2",
		namespace:   "dev",
	},
	{
		buildCommit: "d3a8986332402061f9667aee033b33bd890bb96f",
		configHash:  "1c163de678aac3a8ce1582828fd5bea86773ba19",
		host:        "127.100.10.2:43355",
		source:      "app_2",
		namespace:   "prod",
	},
}

var errTemplates = []string{
	"read login request failed: %v",
	"authenticate password: %v",
	"read registration request failed: %v",
	"register user failed: %v",
	"read logout request failed: %v",
	"remove session token failed: %v",
	"read logout request failed: %v",
	"remove all tokens failed: %v",
	"authenticate token failed: %v",
	"read create catalog request failed: %v",
	"create catalog failed: %v",
	"read find catalog request failed: %v",
	"find catalog by url failed: %v",
	"read update catalog request failed: %v",
	"get catalogs list failed: %v",
	"read update catalog request failed: %v",
	"update catalog failed: %v",
	"read confirm email request failed: %v",
	"confirm email failed: %v",
	"read find user confirm request failed: %v",
	"find user confirm failed: %v",
	"read resend confirm request failed: %v",
	"resend confirm failed: %v",
	"read create item request failed: %v",
	"create item failed: %v",
	"read find item request failed: %v",
	"find item failed: %v",
	"read list item request failed: %v",
	"get items list failed: %v",
	"read update item request failed: %v",
	"update item failed: %v",
	"read create notification addr request failed: %v",
	"create notification addr failed: %v",
	"read list notification addr request failed: %v",
	"list notification addr failed: %v",
	"read update notification addr request failed: %v",
	"update notification addr failed: %v",
	"read find order request failed: %v",
	"find order failed: %v",
	"read find order request failed: %v",
	"find order failed: %v",
	"read list orders request failed: %v",
	"list orders failed: %v",
	"read create order request failed: %v",
	"create order failed: %v",
	"read update order request failed: %v",
	"update order failed: %v",
	"read create shop domain request failed: %v",
	"hold shop domain failed: %v",
	"read create shop request failed: %v",
	"create shop failed: %v",
	"read find shop request failed: %v",
	"find shop failed: %v",
	"read list shops request failed: %v",
	"get shops list failed: %v",
	"read find shop request failed: %v",
	"update shop failed: %v",
	"read list measures request failed: %v",
	"get list measures failed: %v",
	"read exec request failed: %v",
	"close exec request body failed: %v",
	"wrap error: %v",
	"read failed: %v",
	"do failed: %v",
	"read failed: %v",
	"do failed: %v",
	"read failed: %v",
	"do failed: %v",
	"read failed: %v",
	"do failed: %v",
	"list active photos failed: %v",
	"list minio photos failed: %v",
	"generate url failed: %v",
	"process image failed: %v",
	"put image failed: %v",
	"remove image failed: %v",
	"find token key failed: %v",
	"update expire in redis failed: %v",
	"generate session key failed: %v",
	"set data to redis failed: %v",
	"generate sign failed: %v",
	"get data from redis failed: %v",
	"remove data from redis failed: %v",
	"get data from redis failed: %v",
	"create catalog failed: %v",
	"find catalog by url failed: %v",
	"create catalog failed: %v",
	"find shop by domain failed: %v",
	"find catalog failed: %v",
	"find shop by domain failed: %v",
	"get catalogs list failed: %v",
	"find catalogs info failed: %v",
	"update catalog failed: %v",
	"find catalog by id failed: %v",
	"find catalog by url failed: %v",
	"update catalog failed: %v",
	"find catalog failed: %v",
	"find item by url failed: %v",
	"create item failed: %v",
	"find shop by domain failed: %v",
	"find item failed: %v",
	"find shop by domain failed: %v",
	"get items list failed: %v",
	"find item by id failed: %v",
	"update item failed: %v",
	"find item by url failed: %v",
	"check catalog owner failed: %v",
	"find measures failed: %v",
	"find user failed: %v",
	"list notification addr failed: %v",
	"find notification addr failed: %v",
	"store notification addr failed: %v",
	"find user failed: %v",
	"find notification addr failed: %v",
	"find notification failed: %v",
	"update notification failed: %v",
	"find shop failed: %v",
	"confirm order failed: %v",
	"find shop failed: %v",
	"find order failed: %v",
	"get order list by shop id failed: %v",
	"get order items failed: %v",
	"find shop failed: %v",
	"store order item failed: %v",
	"find order by id failed: %v",
	"update order failed: %v",
	"create shop failed: %v",
	"find shop failed: %v",
	"hold shop domain failed: %v",
	"find shops by owner id failed: %v",
	"find shops info failed: %v",
	"find shop failed: %v",
	"update shop failed: %v",
	"authenticate token failed: %v",
	"find shop by domain failed: %v",
	"validate email confirmation failed: %v",
	"confirm user failed: %v",
	"find user failed: %v",
	"find user failed: %v",
	"find user by id failed: %v",
	"find user by id failed: %v",
	"send email confirmation failed: %v",
	"find user by email failed: %v",
	"generate password hash failed: %v",
	"create user failed: %v",
	"generate token failed: %v",
	"add user session failed: %v",
	"send email confirmation message failed: %v",
	"find user failed: %v",
	"check active sessions failed: %v",
	"generate token failed: %v",
	"add user session failed: %v",
	"find user by id failed: %v",
	"update user failed: %v",
	"get file stat failed: %v",
	"convert message failed: %v",
	"convert notification message failed: %v",
	"send order notification failed: %v",
	"convert confirmation message failed: %v",
	"send confirmation failed: %v",
	"store notification failed: %v",
	"list shop notification addr failed: %v",
	"send notification failed: %v",
	"store send info failed: %v",
	"remove image failed: %v",
}

var errText = []string{
	"some message 1",
	"some message 2",
	"some message 3",
	"some message 4",
	"some message 5",
	"some message 6",
}

var jsonRand = []interface{}{
	map[string]interface{}{
		"key1": "value1",
		"key2": "85526454644",
		"key3": []string{"array_val1", "array_val2", "array_val3", "array_val4"},
	},
	"some value",
	[]int{1, 2, 3, 4, 5, 6, 7, 8, 8, 99},
	[]float64{66666666666.66666666666666666, 2344233423, 432432435},
}
