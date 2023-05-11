package mysql

import (
	"github.com/cleancrud/crud/xsql"
	_ "github.com/go-sql-driver/mysql"
)

func NewDB(c *xsql.Config) (*xsql.DB, error) {
	return xsql.NewDB("mysql", c)
}
