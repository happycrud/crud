package xsql

import "golang.org/x/exp/constraints"

type WhereFunc func(*Selector)

type Ops interface {
	constraints.Float | constraints.Integer | []byte | string
}
type FieldOps[T Ops] struct {
	name string
}

func NewFieldOps[T Ops](name string) FieldOps[T] {
	return FieldOps[T]{name: name}
}

func (f FieldOps[T]) EQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(EQ(f.name, arg))
	}
}
func (f FieldOps[T]) NEQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(NEQ(f.name, arg))
	}
}
func (f FieldOps[T]) LT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LT(f.name, arg))
	}
}
func (f FieldOps[T]) LTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LTE(f.name, arg))
	}
}
func (f FieldOps[T]) GT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GT(f.name, arg))
	}
}
func (f FieldOps[T]) GTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GTE(f.name, arg))
	}
}
func (f FieldOps[T]) In(args ...T) WhereFunc {
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
func (f FieldOps[T]) NotIn(args ...T) WhereFunc {
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

type StrFieldOps struct {
	FieldOps[string]
}

func NewStrFieldOps(name string) StrFieldOps {
	return StrFieldOps{FieldOps[string]{name: name}}
}

func (f StrFieldOps) HasPrefix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasPrefix(f.name, arg))
	}
}
func (f StrFieldOps) HasSuffix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasSuffix(f.name, arg))
	}
}
func (f StrFieldOps) Contains(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(Contains(f.name, arg))
	}
}

// And groups predicates with the AND operator between them.
func AndOps(predicates ...WhereFunc) WhereFunc {
	return func(s *Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	}
}

// Or groups predicates with the OR operator between them.
func OrOps(predicates ...WhereFunc) WhereFunc {
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
func NotOps(p WhereFunc) WhereFunc {
	return func(s *Selector) {
		p(s.Not())
	}
}
