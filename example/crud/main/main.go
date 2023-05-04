package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	user "github.com/happycrud/crud/example/crud/user3"
)

var ctx = context.Background()

func main() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")
	a := &user.User{
		Id:    0,
		Name:  "xxx",
		Age:   11,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	fmt.Println(user.Create(db).SetUser(a).Save(ctx))
	fmt.Println(a)

	list, err := user.Find(db).
		Select(user.Id, user.Age).
		Where(user.And(
			user.IdOp.GT(0),
			user.NameOp.NEQ("aa"),
		)).
		All(ctx)
	b, _ := json.Marshal(list)
	fmt.Println(string(b), err)
	user.Update(db).SetName("aaa").Where(user.IdOp.EQ(1)).Save(ctx)

}
