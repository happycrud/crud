package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/happycrud/crud/xsql"
)

func NewDB(c *xsql.Config) (*xsql.DB, error) {
	return xsql.NewDB("mysql", c)
}
