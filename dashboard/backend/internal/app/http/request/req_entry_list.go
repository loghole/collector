package request

import (
	"encoding/json"
	"net/http"

	validation "github.com/gadavy/ozzo-validation/v4"
	"github.com/lissteron/simplerr"

	"github.com/lissteron/loghole/dashboard/internal/app/codes"
	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases"
)

type ListEntryRequest struct {
	Params []*ListEntryParam `json:"params"`
	Limit  int64             `json:"limit"`
	Offset int64             `json:"offset"`
}

func ReadListEntryRequest(r *http.Request) (*ListEntryRequest, error) {
	req := &ListEntryRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, simplerr.WrapWithCode(err, codes.InvalidJSONError, "invalid json struct")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}

func (r *ListEntryRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Params),
		validation.Field(&r.Limit, r.limitRules()...),
		validation.Field(&r.Offset, r.offsetRules()...),
	)
}

func (r *ListEntryRequest) ToInput() *usecases.ListEntryIn {
	query := &domain.Query{
		Params: make([]*domain.QueryParam, 0, len(r.Params)),
		Limit:  r.Limit,
		Offset: r.Offset,
	}

	for _, val := range r.Params {
		if val.Value.isEmpty() {
			continue
		}

		query.Params = append(query.Params, &domain.QueryParam{
			Type:     val.Type,
			Key:      val.Key,
			Value:    domain.ParamValue{Item: val.Value.Item, List: val.Value.List},
			Operator: val.Operator,
		})
	}

	return &usecases.ListEntryIn{
		Query: query,
	}
}

func (r *ListEntryRequest) limitRules() []validation.Rule {
	return []validation.Rule{
		validation.Min(int64(0)).ErrorCode(codes.ValidMinLimit.String()),
		validation.Max(int64(maxLimit)).ErrorCode(codes.ValidMaxLimit.String()),
	}
}

func (r *ListEntryRequest) offsetRules() []validation.Rule {
	return []validation.Rule{
		validation.Min(int64(0)).ErrorCode(codes.ValidMinOffset.String()),
	}
}

type ListEntryParam struct {
	Type     string      `json:"type"`
	Key      string      `json:"key"`
	Value    *ParamValue `json:"value"`
	Operator string      `json:"operator"`
}

func (p *ListEntryParam) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, p.typeRules()...),
		validation.Field(&p.Key, p.keyRules()...),
		validation.Field(&p.Value, p.valueRules()...),
		validation.Field(&p.Operator, p.operatorRules()...),
	)
}

func (p *ListEntryParam) typeRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.ErrorCode(codes.ValidQueryParamsTypeRequired.String()),
		validation.In(domain.TypeKey, domain.TypeColumn).ErrorCode(codes.ValidQueryParamsTypeIn.String()),
	}
}

func (p *ListEntryParam) keyRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.ErrorCode(codes.ValidQueryParamsKeyRequired.String()),
	}
}

func (p *ListEntryParam) valueRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.ErrorCode(codes.ValidQueryParamsValueRequired.String()),
	}
}

func (p *ListEntryParam) operatorRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.ErrorCode(codes.ValidQueryParamsOperatorRequired.String()),
		validation.In(
			domain.OperatorLt,
			domain.OperatorLte,
			domain.OperatorGt,
			domain.OperatorGte,
			domain.OperatorEq,
			domain.OperatorNotEq,
			domain.OperatorLike,
			domain.OperatorNotLike).ErrorCode(codes.ValidQueryParamsOperatorIn.String()),
	}
}

type ParamValue struct {
	Item string   `json:"item"`
	List []string `json:"list"`
}

func (p *ParamValue) isEmpty() bool {
	return len(p.List) == 0 && p.Item == ""
}

func (p *ParamValue) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Item, p.itemRules()...),
		validation.Field(&p.List, p.listRules()...),
	)
}

func (p *ParamValue) itemRules() []validation.Rule {
	return []validation.Rule{
		validation.When(p.List != nil, validation.In("").ErrorCode(codes.ValidQueryParamsValueItemEmpty.String())),
	}
}

func (p *ParamValue) listRules() []validation.Rule {
	return []validation.Rule{
		validation.When(p.Item != "", validation.Length(0, 0).ErrorCode(codes.ValidQueryParamsValueListEmpty.String())),
	}
}
