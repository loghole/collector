package handlers

import (
	"context"
	"net/http"

	"github.com/lissteron/loghole/dashboard/internal/app/http/request"
	"github.com/lissteron/loghole/dashboard/internal/app/http/response"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases"
)

type ListSuggest interface {
	Do(ctx context.Context, input *usecases.ListSuggestIn) (result []string, err error)
}

type SuggestHandlers struct {
	listSuggest ListSuggest
	logger      Logger
}

func NewSuggestHandlers(
	listSuggest ListSuggest,
	logger Logger,
) *SuggestHandlers {
	return &SuggestHandlers{
		listSuggest: listSuggest,
		logger:      logger,
	}
}

func (h *SuggestHandlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := response.NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	req, err := request.ReadListSuggestRequest(r)
	if err != nil {
		h.logger.Errorf(ctx, "read request failed: %v", err)
		resp.ParseError(err)

		return
	}

	result, err := h.listSuggest.Do(ctx, req.ToInput())
	if err != nil {
		h.logger.Errorf(ctx, "do failed: %v", err)
		resp.ParseError(err)

		return
	}

	resp.SetData(result)
}
