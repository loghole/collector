package v1

import (
	"net/http"

	"github.com/loghole/tracing/tracelog"

	"github.com/loghole/collector/config"
)

type info struct {
	ServiceName string `json:"service_name,omitempty"`
	AppName     string `json:"app_name,omitempty"`
	GitHash     string `json:"git_hash,omitempty"`
	Version     string `json:"version,omitempty"`
	BuildAt     string `json:"build_at,omitempty"`
}

type InfoHandlers struct {
	logger tracelog.Logger
	info   info
}

func NewInfoHandlers(logger tracelog.Logger) *InfoHandlers {
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
	resp, ctx := NewBaseResponse(), r.Context()
	defer resp.Write(ctx, w, h.logger)

	resp.SetData(h.info)
}
