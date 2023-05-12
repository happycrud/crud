package mgo

import (
	"time"

	"github.com/happycrud/mgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string
	Age   int
	Sex   bool
	Mtime time.Time
}

func (u *User) GetObjectID() (primitive.ObjectID, string) {
	if u.IsNil() {
		return primitive.NilObjectID, ID
	}
	return u.ID, ID
}

func (u *User) SetObjectID(id primitive.ObjectID) {
	if u.IsNil() {
		return
	}
	u.ID = id
}

func (u *User) IsNil() bool {
	return u == nil
}

const (
	tableName = "user"
	ID        = "_id"
	Name      = "name"
	Age       = "age"
	Sex       = "sex"
	Mtime     = "mtime"
)

func Collection(db *mongo.Database) *mongo.Collection {
	return db.Collection(tableName)
}

func Create(col *mongo.Collection) *mgo.InsertExecutor[*User] {
	return mgo.NewInsertExecutor[*User](col)
}

func Delete(col *mongo.Collection) *mgo.DeleteExecutor[*User] {
	return mgo.NewDeleteExecutor[*User](col)
}

func Find(col *mongo.Collection) *mgo.FinderExecutor[*User] {
	return mgo.NewFinderExecutor[*User](col)
}

func Update(col *mongo.Collection) *Updater {
	return &Updater{mgo.NewUpdateExecutor[*User](col)}
}

type Updater struct {
	*mgo.UpdateExecutor[*User]
}

func (u *Updater) SetID(a primitive.ObjectID) *Updater {
	u.UpdateExecutor.Set(ID, a)
	return u
}
func (u *Updater) SetName(a string) *Updater {
	u.UpdateExecutor.Set(Name, a)
	return u
}
func (u *Updater) SetAge(a int) *Updater {
	u.UpdateExecutor.Set(Age, a)
	return u
}
func (u *Updater) SetSex(a bool) *Updater {
	u.UpdateExecutor.Set(Sex, a)
	return u
}
func (u *Updater) SetMtime(a time.Time) *Updater {
	u.UpdateExecutor.Set(Mtime, a)
	return u
}
