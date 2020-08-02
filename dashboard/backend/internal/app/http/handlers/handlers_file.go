// +build !bindata

package handlers

import (
	"net/http"
)

type FilesHandlers struct {
	path string
}

func NewFilesHandlers(path string) *FilesHandlers {
	return &FilesHandlers{path: path}
}

func (h *FilesHandlers) Handler() http.Handler {
	return http.FileServer(http.Dir(h.path))
}
