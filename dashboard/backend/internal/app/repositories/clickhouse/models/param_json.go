package models

import (
	"fmt"
	"strings"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/repositories/clickhouse/tools"
)

const (
	stringParams = "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] %s %s)"
	floatParams  = "(has(params_float.keys, ?) AND params_float.values[indexOf(params_float.keys, ?)] %s %s)"
)

type JSONParam struct {
	domain.QueryParam
}

func JSONParamFromDomain(param *domain.QueryParam) *JSONParam {
	return &JSONParam{*param}
}

// nolint:golint,stylecheck,gocritic
func (p *JSONParam) ToSql() (query string, args []interface{}, err error) {
	return buildSQL(p)
}

func (p *JSONParam) getIn() (query string, args []interface{}, err error) {
	builder := make([]string, 0, len(p.Value.List))
	args = append(args, p.Key, p.Key)

	for _, value := range p.GetValueList() {
		builder = append(builder, "?")
		args = append(args, value)
	}

	param := strings.Join([]string{"(", strings.Join(builder, ","), ")"}, "")

	return fmt.Sprintf(stringParams, domain.OperatorIn, param), args, nil
}

func (p *JSONParam) getNotIn() (query string, args []interface{}, err error) {
	builder := make([]string, 0, len(p.Value.List))
	args = append(args, p.Key, p.Key)

	for _, value := range p.GetValueList() {
		builder = append(builder, "?")
		args = append(args, value)
	}

	param := strings.Join([]string{"(", strings.Join(builder, ","), ")"}, "")

	return fmt.Sprintf(stringParams, domain.OperatorNotIn, param), args, nil
}

func (p *JSONParam) getLike() (query string, args []interface{}, err error) {
	query = fmt.Sprintf(stringParams, "LIKE", "?")
	args = []interface{}{p.Key, p.Key, tools.CreateLike(p.GetValueItem())}

	return query, args, nil
}

func (p *JSONParam) getNotLike() (query string, args []interface{}, err error) {
	query = fmt.Sprintf(stringParams, "NOT LIKE", "?")
	args = []interface{}{p.Key, p.Key, tools.CreateLike(p.GetValueItem())}

	return query, args, nil
}

func (p *JSONParam) getLikeWithValue(val string) (query string, args []interface{}, err error) {
	query = fmt.Sprintf(stringParams, "LIKE", "?")
	args = []interface{}{p.Key, p.Key, tools.CreateLike(val)}

	return query, args, nil
}

func (p *JSONParam) getNotLikeWithValue(val string) (query string, args []interface{}, err error) {
	query = fmt.Sprintf(stringParams, "NOT LIKE", "?")
	args = []interface{}{p.Key, p.Key, tools.CreateLike(val)}

	return query, args, nil
}

func (p *JSONParam) getLtGtString() (query string, args []interface{}, err error) {
	return fmt.Sprintf(stringParams, getOperator(p.Operator), "?"), []interface{}{p.Key, p.Key, p.Value.Item}, nil
}

func (p *JSONParam) getLtGtFloat(val float64) (query string, args []interface{}, err error) {
	return fmt.Sprintf(floatParams, getOperator(p.Operator), "?"), []interface{}{p.Key, p.Key, val}, nil
}

func (p *JSONParam) getDefault() (query string, args []interface{}, err error) {
	return fmt.Sprintf(stringParams, getOperator(p.Operator), "?"), []interface{}{p.Key, p.Key, p.Value.Item}, nil
}
