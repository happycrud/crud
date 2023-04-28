package sqlite

import (
	"github.com/happycrud/crud/xsql"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLite(c *xsql.Config) (*xsql.DB, error) {
	return xsql.NewDB("sqlite3", c)
}
