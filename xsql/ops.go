package xsql

import "golang.org/x/exp/constraints"

type WhereFunc func(*Selector)

type OpType interface {
	constraints.Float | constraints.Integer | []byte | string
}
type FieldOp[T OpType] struct {
	name string
}

func NewFieldOp[T OpType](name string) FieldOp[T] {
	return FieldOp[T]{name: name}
}

func (f FieldOp[T]) EQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(EQ(f.name, arg))
	}
}
func (f FieldOp[T]) NEQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(NEQ(f.name, arg))
	}
}
func (f FieldOp[T]) LT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LT(f.name, arg))
	}
}
func (f FieldOp[T]) LTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LTE(f.name, arg))
	}
}
func (f FieldOp[T]) GT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GT(f.name, arg))
	}
}
func (f FieldOp[T]) GTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GTE(f.name, arg))
	}
}
func (f FieldOp[T]) In(args ...T) WhereFunc {
	return func(s *Selector) {
		if len(args) == 0 {
			s.Where(False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(In(f.name, v...))
	}
}
func (f FieldOp[T]) NotIn(args ...T) WhereFunc {
	return func(s *Selector) {
		if len(args) == 0 {
			s.Where(Not(False()))
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(NotIn(f.name, v...))
	}
}

type StrFieldOp struct {
	FieldOp[string]
}

func NewStrFieldOp(name string) StrFieldOp {
	return StrFieldOp{FieldOp[string]{name: name}}
}

func (f StrFieldOp) HasPrefix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasPrefix(f.name, arg))
	}
}
func (f StrFieldOp) HasSuffix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasSuffix(f.name, arg))
	}
}
func (f StrFieldOp) Contains(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(Contains(f.name, arg))
	}
}

// And groups predicates with the AND operator between them.
func AndOp(predicates ...WhereFunc) WhereFunc {
	return func(s *Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	}
}

// Or groups predicates with the OR operator between them.
func OrOp(predicates ...WhereFunc) WhereFunc {
	return func(s *Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	}
}

// Not applies the not operator on the given predicate.
func NotOp(p WhereFunc) WhereFunc {
	return func(s *Selector) {
		p(s.Not())
	}
}
