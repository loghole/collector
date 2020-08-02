package domain

const (
	TypeKey    = "key"
	TypeColumn = "column"
)

const (
	OperatorLt      = "<"
	OperatorLte     = "<="
	OperatorGt      = ">"
	OperatorGte     = ">="
	OperatorEq      = "="
	OperatorNotEq   = "!="
	OperatorIn      = "IN"
	OperatorNotIn   = "NOT IN"
	OperatorLike    = "LIKE"
	OperatorNotLike = "NOT LIKE"
)

type Query struct {
	Params []*QueryParam
	Limit  int64
	Offset int64
}

type QueryParam struct {
	Type     string
	Key      string
	Value    ParamValue
	Operator string
}

func (p *QueryParam) IsOperator(operator string) bool {
	return p.Operator == operator
}

func (p *QueryParam) IsList() bool {
	return len(p.Value.List) != 0
}

func (p *QueryParam) IsLike() bool {
	return OperatorLike == p.Operator || OperatorNotLike == p.Operator
}

func (p *QueryParam) IsLtGt() bool {
	switch p.Operator {
	case ">", ">=", "<", "<=":
		return true
	}

	return false
}

func (p *QueryParam) IsIn() bool {
	return p.Operator == "=" && len(p.Value.List) > 0
}

func (p *QueryParam) IsNotIn() bool {
	return p.Operator == "!=" && len(p.Value.List) > 0
}

func (p *QueryParam) IsTypeJSON() bool {
	return p.Type == TypeKey
}

func (p *QueryParam) GetValueList() []string {
	return p.Value.List
}

func (p *QueryParam) GetValueItem() string {
	return p.Value.Item
}

type ParamValue struct {
	Item string
	List []string
}
