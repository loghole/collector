package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracelog"
)

type EntryService interface {
	StoreItem(ctx context.Context, remoteIP string, data []byte) (err error)
	StoreList(ctx context.Context, remoteIP string, data []byte) (err error)
}

type Message struct {
	Event Event `json:"event"`
}

type Event struct {
	Line json.RawMessage `json:"line"`
}

type SplunkHandler struct {
	service EntryService
	logger  tracelog.Logger
	tracer  *tracing.Tracer
}

func NewSplunkHandler(
	service EntryService,
	logger tracelog.Logger,
	tracer *tracing.Tracer,
) *SplunkHandler {
	return &SplunkHandler{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

func (h *SplunkHandler) Handler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read data", http.StatusInternalServerError)

		return
	}

	if len(data) == 0 {
		return
	}

	var (
		ctx  = r.Context()
		addr = r.RemoteAddr
	)

	for {
		idx := bytes.Index(data, []byte("}{"))
		if idx == -1 {
			break
		}

		if err := h.handleMessage(ctx, addr, data[:idx+1]); err != nil {
			http.Error(w, "handle message", http.StatusInternalServerError)

			return
		}

		data = data[idx+1:]
	}

	if err := h.handleMessage(ctx, addr, data); err != nil {
		http.Error(w, "handle message", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SplunkHandler) handleMessage(ctx context.Context, addr string, data []byte) error {
	var dest Message

	if err := json.Unmarshal(data, &dest); err != nil {
		h.logger.Errorf(ctx, "unmarshal message: %v", err)

		return nil
	}

	if err := h.service.StoreItem(ctx, addr, dest.Event.Line); err != nil {
		h.logger.Errorf(ctx, "store item: %v", err)

		return fmt.Errorf("store item: %w", err)
	}

	return nil
}
