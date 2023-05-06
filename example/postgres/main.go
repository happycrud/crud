package main

import (
	"context"
	"fmt"
	"time"

	"github.com/happycrud/crud/example/postgres/crud/user"
	"github.com/happycrud/crud/xsql"
	"github.com/happycrud/crud/xsql/postgres"
)

var db *xsql.DB
var ctx = context.Background()

func main() {
	var err error
	db, err = postgres.NewDB(&xsql.Config{
		DSN:          "postgres://postgres:123456@localhost:5432/postgres",
		ReadDSN:      []string{"postgres://postgres:123456@localhost:5432/postgres"},
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
		Name:  "sdfs",
		Age:   11,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	b := &user.User{
		Id:    1,
		Name:  "a",
		Age:   22,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	_, err = user.Create(debugdb).SetUser(a, b).Upsert(ctx)
	fmt.Println(a, b, err)

}
