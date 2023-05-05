package xsql

import (
	"fmt"
	"testing"
)

func Test_EQ(t *testing.T) {
	var a = FieldOp[int]("age")

	f := a.EQ(1)
	x := Select("xx").From(Table("user"))
	f(x)
	fmt.Println(x.Query())

}
