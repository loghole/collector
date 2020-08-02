package response

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/lissteron/simplerr"

	"github.com/lissteron/loghole/dashboard/internal/app/codes"
)

type Logger interface {
	Errorf(ctx context.Context, template string, args ...interface{})
}

type BaseResponse struct {
	Status int         `json:"-"`
	Errors []RespError `json:"errors"`
	Data   interface{} `json:"data"`
}

type RespError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func NewBaseResponse() *BaseResponse {
	return &BaseResponse{Status: http.StatusOK}
}

func (r *BaseResponse) Write(ctx context.Context, w http.ResponseWriter, log Logger) {
	w.Header().Add("Content-Type", "application/json")

	if r.Status != 0 {
		w.WriteHeader(r.Status)
	}

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Errorf(ctx, "write response failed: %v", err)
	}
}

func (r *BaseResponse) SetData(v interface{}) {
	r.Data = v
}

func (r *BaseResponse) ParseError(err error) {
	switch e1 := err.(type) {
	case validation.Errors:
		r.Status = http.StatusBadRequest

		for field, e2 := range e1 {
			if e3, ok := e2.(validation.Error); ok {
				r.Errors = append(r.Errors, RespError{
					Code:   e3.Code(),
					Detail: strings.Join([]string{field, e3.Error()}, ": "),
				})
			} else {
				r.ParseError(e2)
			}
		}
	default:
		code := simplerr.GetCode(err)

		r.Status = codes.ToHTTP(code)

		r.Errors = append(r.Errors, RespError{
			Code:   strconv.Itoa(code),
			Detail: simplerr.GetText(err),
		})
	}
}
