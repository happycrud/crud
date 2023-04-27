package xsql

import (
	"fmt"
	"strings"
)

var opMap = map[string]Op{
	"=":           OpEQ,
	"<>":          OpNEQ,
	">":           OpGT,
	">=":          OpGTE,
	"<":           OpLT,
	"<=":          OpLTE,
	"IN":          OpIn,
	"NOT IN":      OpNotIn,
	"LIKE":        OpLike,
	"IS NULL":     OpIsNull,
	"IS NOT NULL": OpNotNull,
	"+":           OpAdd,
	"-":           OpSub,
	"*":           OpMul,
	"/":           OpDiv,
	"%":           OpMod,
}

func VailedOp(op string) (vailed bool, t Op) {
	o := strings.ToUpper(strings.TrimSpace(op))
	if p, ok := opMap[o]; ok {
		return true, p
	}
	return false, Op(-1)
}

func GenP(field, op, value string) (*Predicate, error) {
	v, o := VailedOp(op)
	if !v {
		return nil, fmt.Errorf("op:%s is not support", op)
	}
	switch o {
	case OpEQ:
		return EQ(field, value), nil
	case OpNEQ:
		return NEQ(field, value), nil
	case OpGT:
		return GT(field, value), nil
	case OpGTE:
		return GTE(field, value), nil
	case OpLT:
		return LT(field, value), nil
	case OpLTE:
		return LTE(field, value), nil
	case OpIn:
		vs := strings.Split(value, ",")
		is := make([]interface{}, 0, len(vs))
		for _, i := range vs {
			is = append(is, i)
		}
		return In(field, is...), nil
	case OpNotIn:
		vs := strings.Split(value, ",")
		is := make([]interface{}, 0, len(vs))
		for _, i := range vs {
			is = append(is, i)
		}
		return NotIn(field, is...), nil
	case OpLike:
		return Like(field, value), nil
	default:
		return nil, fmt.Errorf("op:%s is not support", op)
	}

}

// NoColumnSelected returns the selected columns in the Selector.
func (s *Selector) NoColumnSelected() bool {
	return len(s.selection) == 0
}
