package handlers

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gadavy/tracing"
	"github.com/lissteron/simplerr"

	"github.com/lissteron/loghole/collector/internal/app/codes"
	"github.com/lissteron/loghole/collector/internal/app/controllers/http/response"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, template string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, template string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, template string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, template string, args ...interface{})
}

type StoreEntryList interface {
	Do(ctx context.Context, data []byte) (err error)
}

type EntryHandlers struct {
	storeList StoreEntryList
	logger    Logger
	tracer    *tracing.Tracer
}

func NewEntryHandlers(
	storeList StoreEntryList,
	logger Logger,
	tracer *tracing.Tracer,
) *EntryHandlers {
	return &EntryHandlers{
		storeList: storeList,
		logger:    logger,
		tracer:    tracer,
	}
}

func (h *EntryHandlers) StoreListHandler(w http.ResponseWriter, r *http.Request) {
	span := h.tracer.NewSpan().WithName(r.URL.String()).Build()
	defer span.Finish()

	resp, ctx := response.NewBaseResponse(), span.Context(r.Context())
	defer resp.Write(ctx, w, h.logger)

	data, err := readData(r.Body)
	if err != nil {
		h.logger.Errorf(ctx, "read body failed: %v", err)
		resp.ParseError(err)

		return
	}

	err = h.storeList.Do(ctx, data)
	if err != nil {
		h.logger.Errorf(ctx, "store entry list failed: %v", err)
		resp.ParseError(err)

		return
	}
}

func (h *EntryHandlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	span := h.tracer.NewSpan().WithName(r.URL.String()).Build()
	defer span.Finish()

	response.NewBaseResponse().Write(span.Context(r.Context()), w, h.logger)
}

func readData(r io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, simplerr.WrapWithCode(err, codes.SystemError, "system error")
	}

	return data, nil
}
