package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/lissteron/simplerr"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracelog"

	"github.com/loghole/collector/internal/app/codes"
)

type EntryService interface {
	Ping(ctx context.Context) error
	StoreItem(ctx context.Context, remoteIP string, data []byte) (err error)
	StoreList(ctx context.Context, remoteIP string, data []byte) (err error)
}

type EntryHandlers struct {
	service EntryService
	logger  tracelog.Logger
	tracer  *tracing.Tracer
}

func NewEntryHandlers(
	service EntryService,
	logger tracelog.Logger,
	tracer *tracing.Tracer,
) *EntryHandlers {
	return &EntryHandlers{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

func (h *EntryHandlers) StoreItemHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	data, err := readData(r.Body)
	if err != nil {
		h.logger.Errorf(ctx, "read body failed: %v", err)
		resp.ParseError(err)

		return
	}

	err = h.service.StoreItem(ctx, r.RemoteAddr, data)
	if err != nil {
		h.logger.Errorf(ctx, "store entry item failed: %v", err)
		resp.ParseError(err)
	}
}

func (h *EntryHandlers) StoreListHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	data, err := readData(r.Body)
	if err != nil {
		h.logger.Errorf(ctx, "read body failed: %v", err)
		resp.ParseError(err)

		return
	}

	err = h.service.StoreList(ctx, r.RemoteAddr, data)
	if err != nil {
		h.logger.Errorf(ctx, "store entry list failed: %v", err)
		resp.ParseError(err)
	}
}

func (h *EntryHandlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	if err := h.service.Ping(ctx); err != nil {
		h.logger.Errorf(ctx, "ping failed: %v", err)
		resp.ParseError(err)
	}
}

func readData(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, simplerr.WrapWithCode(err, simplerr.InternalCode(codes.SystemError), "system error")
	}

	return data, nil
}
