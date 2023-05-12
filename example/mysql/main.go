package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/happycrud/crud/example/mysql/crud"
	"github.com/happycrud/crud/example/mysql/crud/user"
	"github.com/happycrud/crud/xsql"
)

var db *crud.Client
var ctx = context.Background()

func main() {
	var err error
	db, err = crud.NewClient(&xsql.Config{
		DSN:          "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Local",
		ReadDSN:      []string{"root:123456@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Local"},
		Active:       20,
		Idle:         20,
		IdleTimeout:  time.Hour * 24,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
	a := &user.User{
		Id:    0,
		Name:  "a",
		Age:   11,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	db.User.Create().SetUser(a).Save(ctx)
	fmt.Println(a)

	db.User.Update().SetName("xxx").Where(user.IdOp.EQ(4005)).Save(ctx)

	list, err := db.User.Find().Select().Where(user.AgeOp.EQ(11)).All(ctx)
	b, _ := json.Marshal(list)
	fmt.Println(string(b), err)

	db.User.Delete().Where(user.IdOp.EQ(a.Id)).Exec(ctx)

}
