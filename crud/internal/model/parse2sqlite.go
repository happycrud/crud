package model

import "github.com/xwb1989/sqlparser"

func Sqlite3Table(db, path, relative string, notint64 bool, dialect string) *Table {
	return nil
}

func Sqlite3Column(ddl *sqlparser.DDL, notint64 bool) ([]*Column, error) {
	return nil, nil
}
