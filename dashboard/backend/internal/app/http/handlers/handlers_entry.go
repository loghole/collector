package handlers

import (
	"context"
	"net/http"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/http/request"
	"github.com/lissteron/loghole/dashboard/internal/app/http/response"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases"
)

type ListEntry interface {
	Do(ctx context.Context, input *usecases.ListEntryIn) ([]*domain.Entry, error)
}

type EntryHandlers struct {
	listEntry ListEntry
	logger    Logger
}

func NewEntryHandlers(
	listEntry ListEntry,
	logger Logger,
) *EntryHandlers {
	return &EntryHandlers{
		listEntry: listEntry,
		logger:    logger,
	}
}

func (h *EntryHandlers) ListEntryHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := response.NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	req, err := request.ReadListEntryRequest(r)
	if err != nil {
		h.logger.Errorf(ctx, "read request failed: %v", err)
		resp.ParseError(err)

		return
	}

	result, err := h.listEntry.Do(ctx, req.ToInput())
	if err != nil {
		h.logger.Errorf(ctx, "do failed: %v", err)
		resp.ParseError(err)

		return
	}

	resp.SetData(result)
}
