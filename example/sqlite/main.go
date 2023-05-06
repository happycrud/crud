package main

import (
	"context"
	"fmt"
	"time"

	"github.com/happycrud/crud/example/sqlite/crud/user"
	"github.com/happycrud/crud/xsql"
	"github.com/happycrud/crud/xsql/sqlite3"
)

var db *xsql.DB
var ctx = context.Background()

func main() {
	var err error
	db, err = sqlite3.NewDB(&xsql.Config{
		DSN:          "/Users/hongshengjie/db/sqlite.db",
		ReadDSN:      []string{"/Users/hongshengjie/db/sqlite.db"},
		Active:       20,
		Idle:         20,
		IdleTimeout:  time.Hour * 24,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
	debugdb := xsql.Debug(db)
	a := &user.User{
		Id:    0,
		Name:  "a",
		Age:   11,
		Ctime: time.Now().Unix(),
		Mtime: time.Now().Unix(),
	}
	b := &user.User{
		Id:    1,
		Name:  "xa",
		Age:   11,
		Ctime: time.Now().Unix(),
		Mtime: time.Now().Unix(),
	}
	r, err := user.Create(debugdb).SetUser(a, b).Upsert(ctx)
	fmt.Println(a, b, err, r)

	//db.User.Update().SetName("xxx").Where(user.IdOp.EQ(4005)).Save(ctx)

	//list, err := db.User.Find().Select().Where(user.AgeOp.EQ(11)).All(ctx)
	//b, _ := json.Marshal(list)
	//fmt.Println(string(b), err)

	//db.User.Delete().Where(user.IdOp.EQ(a.Id)).Exec(ctx)

}
