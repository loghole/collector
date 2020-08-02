package models

type paramDomain interface {
	IsIn() bool
	IsNotIn() bool
	IsLike() bool
	IsList() bool
	IsLtGt() bool
	IsOperator(operator string) bool
	GetValueList() []string
	GetValueItem() string
}

type paramIntGet interface {
	getIn() (string, []interface{}, error)
	getNotIn() (string, []interface{}, error)
	getLike() (string, []interface{}, error)
	getNotLike() (string, []interface{}, error)

	getLikeWithValue(val string) (string, []interface{}, error)
	getNotLikeWithValue(val string) (string, []interface{}, error)

	getLtGtString() (string, []interface{}, error)
	getLtGtFloat(val float64) (string, []interface{}, error)

	getDefault() (string, []interface{}, error)
}

type ParamInt interface {
	paramDomain
	paramIntGet
}
