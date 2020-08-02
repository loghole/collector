package models

import (
	"errors"
	"strconv"
	"strings"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
)

var (
	ErrNotImplemented   = errors.New("not implemented")
	ErrArrayNotAccepted = errors.New("arrays not accepted")
)

func buildSQL(p ParamInt) (query string, args []interface{}, err error) {
	switch {
	case p.IsIn():
		return p.getIn()
	case p.IsNotIn():
		return p.getNotIn()
	case p.IsLike():
		return prepareParamLike(p)
	case p.IsLtGt():
		return prepareParamLtGt(p)
	default:
		return p.getDefault()
	}
}

func prepareParamLike(p ParamInt) (query string, args []interface{}, err error) {
	if p.IsList() {
		return prepareParamListLike(p)
	}

	switch {
	case p.IsOperator(domain.OperatorLike):
		return p.getLike()
	case p.IsOperator(domain.OperatorNotLike):
		return p.getNotLike()
	default:
		panic(ErrNotImplemented)
	}
}

func prepareParamListLike(p ParamInt) (query string, args []interface{}, err error) {
	var (
		queries = make([]string, 0, len(p.GetValueList()))
		a       = make([]interface{}, 0)
		q       string
	)

	for _, value := range p.GetValueList() {
		switch {
		case p.IsOperator(domain.OperatorLike):
			q, a, err = p.getLikeWithValue(value)
		case p.IsOperator(domain.OperatorNotLike):
			q, a, err = p.getNotLikeWithValue(value)
		default:
			panic(ErrNotImplemented)
		}

		if err != nil {
			return "", nil, err
		}

		queries = append(queries, q)
		args = append(args, a...)
	}

	return strings.Join(queries, " AND "), args, nil
}

func prepareParamLtGt(p ParamInt) (query string, args []interface{}, err error) {
	if p.IsList() {
		return "", nil, ErrArrayNotAccepted
	}

	if value, ok := valueToFloat(p.GetValueItem()); ok {
		return p.getLtGtFloat(value)
	}

	return p.getLtGtString()
}

func valueToFloat(val string) (float64, bool) {
	f, err := strconv.ParseFloat(val, 64)

	return f, err == nil
}

func getOperator(operator string) string {
	switch operator {
	case domain.OperatorLt:
		return domain.OperatorLt
	case domain.OperatorLte:
		return domain.OperatorLte
	case domain.OperatorGt:
		return domain.OperatorGt
	case domain.OperatorGte:
		return domain.OperatorGte
	case domain.OperatorNotEq:
		return domain.OperatorNotEq
	case domain.OperatorLike:
		return domain.OperatorLike
	case domain.OperatorNotLike:
		return domain.OperatorNotLike
	default:
		return domain.OperatorEq
	}
}
