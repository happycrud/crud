package xsql

import (
	"fmt"
	"testing"
)

func Test_EQ(t *testing.T) {
	var a StrFieldOps
	a.Name = "xxx"
	f := a.EQ("asaa")
	x := Select("xx").From(Table("use"))
	f(x)
	fmt.Println(x.Query())

}
