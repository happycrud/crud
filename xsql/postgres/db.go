package postgres

import (
	"github.com/cleancrud/crud/xsql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(c *xsql.Config) (*xsql.DB, error) {
	return xsql.NewDB("pgx", c)
}
