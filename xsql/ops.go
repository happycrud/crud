package xsql

import "golang.org/x/exp/constraints"

type WhereFunc func(*Selector)

type OpType interface {
	constraints.Float | constraints.Integer | []byte | string
}
type FieldOp[T OpType] string

func (f FieldOp[T]) EQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(EQ(string(f), arg))
	}
}
func (f FieldOp[T]) NEQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(NEQ(string(f), arg))
	}
}
func (f FieldOp[T]) LT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LT(string(f), arg))
	}
}
func (f FieldOp[T]) LTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LTE(string(f), arg))
	}
}
func (f FieldOp[T]) GT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GT(string(f), arg))
	}
}
func (f FieldOp[T]) GTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GTE(string(f), arg))
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
		s.Where(In(string(f), v...))
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
		s.Where(NotIn(string(f), v...))
	}
}

type StrFieldOp string

func (f StrFieldOp) EQ(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(EQ(string(f), arg))
	}
}
func (f StrFieldOp) NEQ(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(NEQ(string(f), arg))
	}
}
func (f StrFieldOp) LT(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(LT(string(f), arg))
	}
}
func (f StrFieldOp) LTE(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(LTE(string(f), arg))
	}
}
func (f StrFieldOp) GT(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(GT(string(f), arg))
	}
}
func (f StrFieldOp) GTE(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(GTE(string(f), arg))
	}
}
func (f StrFieldOp) In(args ...string) WhereFunc {
	return func(s *Selector) {
		if len(args) == 0 {
			s.Where(False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(In(string(f), v...))
	}
}
func (f StrFieldOp) NotIn(args ...string) WhereFunc {
	return func(s *Selector) {
		if len(args) == 0 {
			s.Where(Not(False()))
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(NotIn(string(f), v...))
	}
}

func (f StrFieldOp) IsNull() WhereFunc {
	return func(s *Selector) {
		s.Where(IsNull(string(f)))
	}
}

func (f StrFieldOp) NotNull() WhereFunc {
	return func(s *Selector) {
		s.Where(NotNull(string(f)))
	}
}

func (f StrFieldOp) HasPrefix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasPrefix(string(f), arg))
	}
}

func (f StrFieldOp) HasSuffix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasSuffix(string(f), arg))
	}
}

func (f StrFieldOp) Contains(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(Contains(string(f), arg))
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
