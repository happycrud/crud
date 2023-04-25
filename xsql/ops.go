package xsql

type WhereFunc func(*Selector)

type FieldOps[T int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8 | float32 | float64 | []byte | string] struct {
	Name string
}

func (f FieldOps[T]) EQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(EQ(f.Name, arg))
	}
}
func (f FieldOps[T]) NEQ(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(NEQ(f.Name, arg))
	}
}
func (f FieldOps[T]) LT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LT(f.Name, arg))
	}
}
func (f FieldOps[T]) LTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(LTE(f.Name, arg))
	}
}
func (f FieldOps[T]) GT(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GT(f.Name, arg))
	}
}
func (f FieldOps[T]) GTE(arg T) WhereFunc {
	return func(s *Selector) {
		s.Where(GTE(f.Name, arg))
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
		s.Where(In(f.Name, v...))
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
		s.Where(NotIn(f.Name, v...))
	}
}

type StrFieldOps struct {
	FieldOps[string]
}

func (f StrFieldOps) HasPrefix(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(HasPrefix(f.Name, arg))
	}
}

func (f StrFieldOps) Contains(arg string) WhereFunc {
	return func(s *Selector) {
		s.Where(Contains(f.Name, arg))
	}
}
