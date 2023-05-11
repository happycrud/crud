package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entity interface {
	GetObjectID() (primitive.ObjectID, string)
	SetObjectID(p primitive.ObjectID)
	IsNil() bool
}

type FinderExecutor[T Entity] struct {
	col     *mongo.Collection
	filters primitive.D
	opts    *options.FindOptions
	a       T
}

func NewFinderExecutor[T Entity](col *mongo.Collection) *FinderExecutor[T] {
	return &FinderExecutor[T]{col: col, opts: options.Find()}
}

func (f *FinderExecutor[T]) Filter(filter ...primitive.E) *FinderExecutor[T] {
	f.filters = append(f.filters, filter...)
	return f
}
func (f *FinderExecutor[T]) Limit(l int64) *FinderExecutor[T] {
	f.opts.SetLimit(l)
	return f
}

func (f *FinderExecutor[T]) SortDesc(field string) *FinderExecutor[T] {
	f.opts.SetSort(primitive.D{{Key: field, Value: -1}})
	return f
}
func (f *FinderExecutor[T]) SortAsc(field string) *FinderExecutor[T] {
	f.opts.SetSort(primitive.D{{Key: field, Value: 1}})
	return f
}

func (f *FinderExecutor[T]) Skip(s int64) *FinderExecutor[T] {
	f.opts.SetSkip(s)
	return f
}
func (f *FinderExecutor[T]) One(ctx context.Context) (T, error) {
	f.opts = f.opts.SetLimit(1)
	ret, err := f.All(ctx)
	if err != nil {
		return f.a, err
	}
	if len(ret) == 1 {
		return ret[0], nil
	}
	return f.a, mongo.ErrNoDocuments
}

func (f *FinderExecutor[T]) All(ctx context.Context) ([]T, error) {
	cursor, err := f.col.Find(ctx, f.filters, f.opts)
	if err != nil {
		return nil, err
	}
	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

type InsertExecutor[T Entity] struct {
	col   *mongo.Collection
	items []T
	a     T
}

func NewInsertExecutor[T Entity](col *mongo.Collection) *InsertExecutor[T] {
	return &InsertExecutor[T]{col: col}
}

func (i *InsertExecutor[T]) SetItem(u ...T) *InsertExecutor[T] {
	i.items = append(i.items, u...)
	return i
}

func (i *InsertExecutor[T]) Save(ctx context.Context) error {
	var insertItems []interface{}
	for _, v := range i.items {
		insertItems = append(insertItems, v)
	}
	ret, err := i.col.InsertMany(ctx, insertItems)
	if err != nil {
		return err
	}
	if _, name := i.a.GetObjectID(); name != "" {
		for idx, v := range ret.InsertedIDs {
			if val, _ := i.items[idx].GetObjectID(); val.IsZero() {
				i.items[idx].SetObjectID(v.(primitive.ObjectID))
			}
		}
	}

	return nil
}

type UpdateExecutor[T Entity] struct {
	col *mongo.Collection
	up  primitive.D
	a   T
}

func NewUpdateExecutor[T Entity](col *mongo.Collection) *UpdateExecutor[T] {
	return &UpdateExecutor[T]{col: col}
}

func (u *UpdateExecutor[T]) Set(name string, arg any) *UpdateExecutor[T] {
	u.up = append(u.up, primitive.E{
		Key:   name,
		Value: arg,
	})
	return u
}

func (u *UpdateExecutor[T]) ByID(ctx context.Context, a primitive.ObjectID) (int64, error) {
	ret, err := u.col.UpdateByID(ctx, a, primitive.D{primitive.E{Key: "$set", Value: u.up}})
	if err != nil {
		return 0, err
	}
	return ret.ModifiedCount, nil
}

type DeleteExecutor[T Entity] struct {
	col *mongo.Collection
}

func NewDeleteExecutor[T Entity](col *mongo.Collection) *DeleteExecutor[T] {
	return &DeleteExecutor[T]{col: col}
}
func (d *DeleteExecutor[T]) ByID(ctx context.Context, a primitive.ObjectID) (int64, error) {
	ret, err := d.col.DeleteOne(ctx, primitive.D{primitive.E{Key: "_id", Value: a}})
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}
