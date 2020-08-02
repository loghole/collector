package request

import (
	"encoding/json"
	"net/http"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/lissteron/simplerr"

	"github.com/lissteron/loghole/dashboard/internal/app/codes"
	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases"
)

type ListSuggestRequest struct {
	Type  string `json:"-"`
	Value string `json:"value"`
}

func ReadListSuggestRequest(r *http.Request) (*ListSuggestRequest, error) {
	req := &ListSuggestRequest{Type: mux.Vars(r)["type"]}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, simplerr.WrapWithCode(err, codes.InvalidJSONError, "invalid json struct")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}

func (r *ListSuggestRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Type, r.typeRules()...),
		validation.Field(&r.Value, r.valueRules()...),
	)
}

func (r *ListSuggestRequest) ToInput() *usecases.ListSuggestIn {
	return &usecases.ListSuggestIn{
		Type:  r.Type,
		Value: r.Value,
	}
}

func (r *ListSuggestRequest) typeRules() []validation.Rule {
	return []validation.Rule{
		validation.In(
			domain.SuggestLevel,
			domain.SuggestHost,
			domain.SuggestSource,
			domain.SuggestNamespace).ErrorCode(codes.ValidSuggestTypeIn.String()),
	}
}

func (r *ListSuggestRequest) valueRules() []validation.Rule {
	return []validation.Rule{
		validation.Length(0, maxSuggestValue).ErrorCode(codes.ValidSuggestValueLength.String()),
	}
}
