package postgres

import (
	"github.com/happycrud/crud/xsql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(c *xsql.Config) (*xsql.DB, error) {
	return xsql.NewDB("pgx", c)
}
