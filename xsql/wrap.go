package xsql

import (
	"context"
	"database/sql"
	"errors"

	"time"

	"github.com/cleancrud/crud/xsql/dialect"
)

type Entity interface {
	Values() []any
	GetAutoIncrPk() (int64, string)
	SetAutoIncrPk(id int64)
	ScanDst(aa any, columns []string) []any
	Columns() []string
	ColumnsSet() map[string]struct{}
	NewPtr() any
	IsNil() bool
	Dialect() string
	Table() string
	Schema() string
}

type InsertExecutor[T Entity] struct {
	eq      ExecQuerier
	builder *InsertBuilder
	items   []T
	upsert  bool
	timeout time.Duration
	a       T
}

func NewInsertExecutor[T Entity](eq ExecQuerier) *InsertExecutor[T] {
	builder := &InsertExecutor[T]{
		eq: eq,
	}
	builder.builder = Dialect(builder.a.Dialect()).Insert(builder.a.Table()).Schema(builder.a.Schema())
	return builder
}

func (in *InsertExecutor[T]) Timeout(t time.Duration) *InsertExecutor[T] {
	in.timeout = t
	return in
}

func (in *InsertExecutor[T]) SetItems(a ...T) *InsertExecutor[T] {
	in.items = append(in.items, a...)
	return in
}

func (in *InsertExecutor[T]) Upsert(ctx context.Context) (int64, error) {
	in.upsert = true
	return in.Save(ctx)
}

// Save Save one or many records set by SetXXX method
// if insert a record , the LastInsertId  will be setted on the struct's  PrimeKey field
// if insert many records , every struct's PrimeKey field will be setted
// return number of RowsAffected or error
func (in *InsertExecutor[T]) Save(ctx context.Context) (int64, error) {
	if len(in.items) == 0 {
		return 0, errors.New("please set a item")
	}

	_, ctx, cancel := Shrink(ctx, in.timeout)
	defer cancel()

	switch in.a.Dialect() {
	case dialect.MySQL:
		in.builder.Columns(in.a.Columns()...)
		if in.upsert {
			in.builder.OnConflict(ResolveWithNewValues())
		}
		for _, a := range in.items {
			if a.IsNil() {
				return 0, errors.New("can not insert a nil item")
			}
			in.builder.Values(a.Values()...)
		}
		ins, args := in.builder.Query()
		result, err := in.eq.ExecContext(ctx, ins, args...)
		if err != nil {
			return 0, err
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return rowsAffected, err
		}
		if lastInsertId > 0 && rowsAffected > 0 {
			for _, v := range in.items {
				if id, autoIncrPkName := v.GetAutoIncrPk(); id > 0 || autoIncrPkName == "" {
					continue
				}
				v.SetAutoIncrPk(lastInsertId)
				lastInsertId++
			}
		}

		return result.RowsAffected()
	case dialect.Postgres, dialect.SQLite:
		_, autoIncrPkName := in.a.GetAutoIncrPk()
		if in.upsert {
			in.builder.Columns(in.a.Columns()...)
			in.builder.OnConflict(
				ConflictColumns(autoIncrPkName),
				ResolveWithNewValues())
			for _, a := range in.items {
				if a.IsNil() {
					return 0, errors.New("can not insert a nil item")
				}
				in.builder.Values(a.Values()...)
			}
			ins, args := in.builder.Query()
			result, err := in.eq.ExecContext(ctx, ins, args...)
			if err != nil {
				return 0, err
			}
			return result.RowsAffected()

		} else {
			insertColumnName := in.a.Columns()

			pkIndex := -1
			if autoIncrPkName != "" {
				in.builder.Returning(autoIncrPkName)
				// 移除自增自增主键名称
				for k, v := range insertColumnName {
					if v == autoIncrPkName {
						insertColumnName = append(insertColumnName[:k], insertColumnName[k+1:]...)
						pkIndex = k
						break
					}
				}
			}
			in.builder.Columns(insertColumnName...)
			for _, a := range in.items {
				if a.IsNil() {
					return 0, errors.New("can not insert a nil item")
				}
				insertValues := a.Values()
				if autoIncrPkName != "" {
					// 移除自增 ==0 的自增主键名称
					if pkIndex >= 0 {
						insertValues = append(insertValues[:pkIndex], insertValues[pkIndex+1:]...)
					}
				}
				in.builder.Values(insertValues...)
			}
			ins, args := in.builder.Query()
			q, err := in.eq.QueryContext(ctx, ins, args...)
			if err != nil {
				return 0, err
			}
			defer q.Close()

			index := 0
			for q.Next() {
				var tempid int64
				if e := q.Scan(&tempid); e != nil {
					return 0, e
				}
				if id, autoPkName := in.items[index].GetAutoIncrPk(); id == 0 && autoPkName != "" {
					in.items[index].SetAutoIncrPk(tempid)
				}
				index++
			}

			if q.Err() != nil {
				return 0, q.Err()
			}
			return int64(len(in.items)), nil
		}

	default:
		return 0, errors.New("not support dialect")
	}
}

type DeleteExecutor[T Entity] struct {
	builder *DeleteBuilder
	eq      ExecQuerier
	timeout time.Duration
	a       T
}

func NewDeleteExecutor[T Entity](eq ExecQuerier) *DeleteExecutor[T] {
	builder := &DeleteExecutor[T]{
		eq: eq,
	}
	builder.builder = Dialect(builder.a.Dialect()).Delete(builder.a.Table()).Schema(builder.a.Schema())
	return builder
}

func (d *DeleteExecutor[T]) Timeout(t time.Duration) *DeleteExecutor[T] {
	d.timeout = t
	return d
}

func (d *DeleteExecutor[T]) Where(p ...WhereFunc) *DeleteExecutor[T] {
	s := &Selector{}
	for _, v := range p {
		v(s)
	}
	d.builder = d.builder.Where(s.P())
	return d
}

func (d *DeleteExecutor[T]) Exec(ctx context.Context) (int64, error) {
	_, ctx, cancel := Shrink(ctx, d.timeout)
	defer cancel()
	del, args := d.builder.Query()
	res, err := d.eq.ExecContext(ctx, del, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type SelectExecutor[T Entity] struct {
	builder *Selector
	eq      ExecQuerier
	timeout time.Duration
	a       T
}

func NewSelectExecutor[T Entity](eq ExecQuerier) *SelectExecutor[T] {
	sel := &SelectExecutor[T]{
		eq: eq,
	}
	sel.builder = Dialect(sel.a.Dialect()).Select().From(Table(sel.a.Table()).Schema(sel.a.Schema()))
	return sel
}

func (s *SelectExecutor[T]) Timeout(t time.Duration) *SelectExecutor[T] {
	s.timeout = t
	return s
}

// Select Select
func (s *SelectExecutor[T]) Select(columns ...string) *SelectExecutor[T] {
	s.builder.Select(columns...)
	return s
}

func (s *SelectExecutor[T]) Count(columns ...string) *SelectExecutor[T] {
	s.builder.Count(columns...)
	return s
}

func (s *SelectExecutor[T]) Where(p ...WhereFunc) *SelectExecutor[T] {
	sel := &Selector{}
	for _, v := range p {
		v(sel)
	}
	s.builder = s.builder.Where(sel.P())
	return s
}

func (s *SelectExecutor[T]) WhereP(ps ...*Predicate) *SelectExecutor[T] {
	for _, v := range ps {
		s.builder.Where(v)
	}
	return s
}

func (s *SelectExecutor[T]) Offset(offset int32) *SelectExecutor[T] {
	s.builder = s.builder.Offset(int(offset))
	return s
}

func (s *SelectExecutor[T]) Limit(limit int32) *SelectExecutor[T] {
	s.builder = s.builder.Limit(int(limit))
	return s
}

func (s *SelectExecutor[T]) OrderDesc(field string) *SelectExecutor[T] {
	s.builder = s.builder.OrderBy(Desc(field))
	return s
}

func (s *SelectExecutor[T]) OrderAsc(field string) *SelectExecutor[T] {
	s.builder = s.builder.OrderBy(Asc(field))
	return s
}

// ForceIndex ForceIndex  FORCE INDEX (`index_name`)
func (s *SelectExecutor[T]) ForceIndex(indexName ...string) *SelectExecutor[T] {
	s.builder.For(LockUpdate)
	return s
}

func (s *SelectExecutor[T]) GroupBy(fields ...string) *SelectExecutor[T] {
	s.builder.GroupBy(fields...)
	return s
}

func (s *SelectExecutor[T]) Having(p *Predicate) *SelectExecutor[T] {
	s.builder.Having(p)
	return s
}

func (s *SelectExecutor[T]) Slice(ctx context.Context, dstSlice interface{}) error {
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return err
	}
	defer q.Close()
	return ScanSlice(q, dstSlice)
}

func (s *SelectExecutor[T]) One(ctx context.Context) (T, error) {
	s.builder.Limit(1)
	results, err := s.All(ctx)
	if err != nil {
		return s.a, err
	}
	if len(results) <= 0 {
		return s.a, sql.ErrNoRows
	}
	return results[0], nil
}

func (s *SelectExecutor[T]) Int64(ctx context.Context) (int64, error) {
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	return Int64(ctx, s.builder, s.eq)
}

func (s *SelectExecutor[T]) Int64s(ctx context.Context) ([]int64, error) {
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	return Int64s(ctx, s.builder, s.eq)
}

func (s *SelectExecutor[T]) String(ctx context.Context) (string, error) {
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	return String(ctx, s.builder, s.eq)
}

func (s *SelectExecutor[T]) Strings(ctx context.Context) ([]string, error) {
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	return Strings(ctx, s.builder, s.eq)
}

func (s *SelectExecutor[T]) selectCheck(columns []string) error {
	set := s.a.ColumnsSet()
	for _, v := range columns {
		if _, ok := set[v]; !ok {
			return errors.New("struct not have field:" + v)
		}
	}
	return nil
}

func (s *SelectExecutor[T]) All(ctx context.Context) ([]T, error) {
	var selectedColumns []string
	if s.builder.NoColumnSelected() {
		s.builder.Select(s.a.Columns()...)
		selectedColumns = s.a.Columns()
	} else {
		selectedColumns = s.builder.SelectedColumns()
		if err := s.selectCheck(selectedColumns); err != nil {
			return nil, err
		}
	}
	_, ctx, cancel := Shrink(ctx, s.timeout)
	defer cancel()
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var result []T
	for q.Next() {
		x := s.a.NewPtr()
		dst := s.a.ScanDst(x, selectedColumns)
		if err := q.Scan(dst...); err != nil {
			return nil, err
		}
		result = append(result, x.(T))
	}
	if q.Err() != nil {
		return nil, q.Err()
	}
	return result, nil
}

type UpdateExecutor[T Entity] struct {
	builder *UpdateBuilder
	eq      ExecQuerier
	timeout time.Duration
	a       T
}

func NewUpdateExecutor[T Entity](eq ExecQuerier) *UpdateExecutor[T] {
	builder := &UpdateExecutor[T]{
		eq: eq,
	}
	builder.builder = Dialect(builder.a.Dialect()).Update(builder.a.Table()).Schema(builder.a.Schema())
	return builder
}

func (u *UpdateExecutor[T]) Timeout(t time.Duration) *UpdateExecutor[T] {
	u.timeout = t
	return u
}

func (u *UpdateExecutor[T]) Where(p ...WhereFunc) *UpdateExecutor[T] {
	s := &Selector{}
	for _, v := range p {
		v(s)
	}
	u.builder = u.builder.Where(s.P())
	return u
}

func (u *UpdateExecutor[T]) Set(name string, arg any) *UpdateExecutor[T] {
	u.builder.Set(name, arg)
	return u
}

func (u *UpdateExecutor[T]) Add(name string, arg any) *UpdateExecutor[T] {
	u.builder.Add(name, arg)
	return u
}

func (u *UpdateExecutor[T]) Save(ctx context.Context) (int64, error) {
	_, ctx, cancel := Shrink(ctx, u.timeout)
	defer cancel()
	up, args := u.builder.Query()
	result, err := u.eq.ExecContext(ctx, up, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
