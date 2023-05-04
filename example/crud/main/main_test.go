package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	user "github.com/happycrud/crud/example/crud/user3"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")

}

func Test_main(t *testing.B) {
	a := &user.User{
		Id:    0,
		Name:  "xxx",
		Age:   11,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	fmt.Println(a)

}
