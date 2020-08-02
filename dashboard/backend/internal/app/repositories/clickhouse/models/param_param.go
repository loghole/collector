package models

import (
	"strings"

	"github.com/Masterminds/squirrel"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/repositories/clickhouse/tools"
)

type ColumnParam struct {
	domain.QueryParam
}

func ColumnParamFromDomain(param *domain.QueryParam) *ColumnParam {
	return &ColumnParam{*param}
}

// nolint:golint,stylecheck,gocritic
func (p *ColumnParam) ToSql() (query string, args []interface{}, err error) {
	return buildSQL(p)
}

func (p *ColumnParam) getIn() (query string, args []interface{}, err error) {
	return squirrel.Eq{p.Key: p.Value.List}.ToSql()
}

func (p *ColumnParam) getNotIn() (query string, args []interface{}, err error) {
	return squirrel.NotEq{p.Key: p.Value.List}.ToSql()
}

func (p *ColumnParam) getLike() (query string, args []interface{}, err error) {
	return squirrel.Like{p.Key: tools.CreateLike(p.Value.Item)}.ToSql()
}

func (p *ColumnParam) getNotLike() (query string, args []interface{}, err error) {
	return squirrel.NotLike{p.Key: tools.CreateLike(p.Value.Item)}.ToSql()
}

func (p *ColumnParam) getLikeWithValue(val string) (query string, args []interface{}, err error) {
	return squirrel.Like{p.Key: tools.CreateLike(val)}.ToSql()
}

func (p *ColumnParam) getNotLikeWithValue(val string) (query string, args []interface{}, err error) {
	return squirrel.NotLike{p.Key: tools.CreateLike(val)}.ToSql()
}

func (p *ColumnParam) getLtGtString() (query string, args []interface{}, err error) {
	return p.getDefault()
}

func (p *ColumnParam) getLtGtFloat(val float64) (query string, args []interface{}, err error) {
	return p.getDefault()
}

func (p *ColumnParam) getDefault() (query string, args []interface{}, err error) {
	return strings.Join([]string{p.Key, getOperator(p.Operator), "?"}, ""), []interface{}{p.Value.Item}, nil
}
