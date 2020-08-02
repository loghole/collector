package handlers

import (
	"net/http"

	"github.com/lissteron/loghole/collector/config"
	"github.com/lissteron/loghole/collector/internal/app/controllers/http/response"
)

type info struct {
	ServiceName string `json:"service_name,omitempty"`
	AppName     string `json:"app_name,omitempty"`
	GitHash     string `json:"git_hash,omitempty"`
	Version     string `json:"version,omitempty"`
	BuildAt     string `json:"build_at,omitempty"`
}

type InfoHandlers struct {
	logger Logger
	info   info
}

func NewInfoHandlers(logger Logger) *InfoHandlers {
	return &InfoHandlers{
		logger: logger,
		info: info{
			ServiceName: config.ServiceName,
			AppName:     config.AppName,
			GitHash:     config.GitHash,
			Version:     config.Version,
			BuildAt:     config.BuildAt,
		},
	}
}

func (h *InfoHandlers) InfoHandler(w http.ResponseWriter, r *http.Request) {
	resp, ctx := response.NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	resp.SetData(h.info)
}
