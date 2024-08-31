package model

import (
	"log"
	"os"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/test_driver"
	"github.com/pingcap/parser/types"
)

func MysqlTable(db, path, relative string, dialect string) *Table {
	sql, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	p, _, err := parser.New().Parse(string(sql), "", "")
	if err != nil {
		log.Fatal(err)
	}
	stmt, ok := p[0].(*ast.CreateTableStmt)
	if !ok {
		log.Fatal("please check sql file statement is DDL")
	}
	var buf strings.Builder
	if err = stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)); err != nil {
		log.Fatal(err)
	}
	tableName := stmt.Table.Name.String()

	gotableName := GoCamelCase(tableName)
	mytable := &Table{
		Database:    db,
		TableName:   tableName,
		GoTableName: gotableName,
		PackageName: strings.ToLower(gotableName),
		Dialect:     dialect,
	}
	columns, err := MysqlColumn(stmt)
	if err != nil {
		log.Fatal(err)
	}
	if len(columns) <= 0 {
		log.Fatal("schema or table not exist")
	}
	mytable.Fields = columns
	for _, v := range columns {
		if v.IsPrimaryKey {
			mytable.PrimaryKey = v
		}
		if v.GoColumnType == "time.Time" {
			mytable.ImportTime = true
		}
		v.ProtoType = GoTypeToProtoType(v.GoColumnType)
	}

	mytable.GenerateWhereCol = mytable.Fields
	mytable.RelativePath = relative
	return mytable
}

func MysqlColumn(ddl *ast.CreateTableStmt) ([]*Column, error) {
	res := []*Column{}
	for k, v := range ddl.Cols {

		var ct string
		if mysql.HasUnsignedFlag(v.Tp.Flag) {
			ct = "unsigned"
		}
		var defaultTs bool
		var notNull bool
		var autoIncrement bool
		var primaryKey bool
		var comment string
		for _, v2 := range v.Options {
			switch v2.Tp {
			case ast.ColumnOptionDefaultValue:
				if v2.Expr.GetFlag() == ast.FlagHasFunc && v2.Expr.(*ast.FuncCallExpr).FnName.L == "current_timestamp" {
					defaultTs = true
				}
			case ast.ColumnOptionComment:
				comment = v2.Expr.(*test_driver.ValueExpr).GetString()
			case ast.ColumnOptionNotNull:
				notNull = true
			case ast.ColumnOptionAutoIncrement:
				autoIncrement = true
			case ast.ColumnOptionPrimaryKey:
				primaryKey = true
			default:
			}
		}
		columnType := types.TypeToStr(v.Tp.Tp, v.Tp.Charset)
		c := &Column{
			OrdinalPosition:           k,
			ColumnName:                v.Name.Name.L,
			DataType:                  columnType,
			ColumnType:                ct,
			ColumnComment:             comment,
			NotNull:                   notNull,
			IsPrimaryKey:              primaryKey,
			IsAutoIncrment:            autoIncrement,
			IsDefaultCurrentTimestamp: defaultTs,
			GoColumnName:              "",
			GoColumnType:              "",
			BigType:                   0,
			GoConditionType:           "",
		}

		c.GoColumnName = GoCamelCase(c.ColumnName)
		c.GoColumnType, c.BigType = MysqlToGoFieldType(c.DataType, c.ColumnType)
		if strings.Contains(c.GoColumnType, "int") {
			c.GoColumnType = "int64"
		}
		c.GoConditionType = c.GoColumnType
		if c.BigType == bigtypeCompareTime {
			c.GoConditionType = "string"
		}
		res = append(res, c)

	}
	var primaryKey string
	for _, v := range ddl.Constraints {
		if v.Tp == ast.ConstraintPrimaryKey {
			if len(v.Keys) > 1 {
				log.Fatal("primary key is not single column")
			}
			primaryKey = v.Keys[0].Column.Name.String()
		}
	}
	for _, v := range res {
		if v.ColumnName == primaryKey {
			v.IsPrimaryKey = true
		}
		v.ProtoType = GoTypeToProtoType(v.GoColumnType)
	}
	return res, nil
}

// MysqlToGoFieldType MysqlToGoFieldType
func MysqlToGoFieldType(dt, ct string) (string, int) {
	var unsigned bool
	if strings.Contains(ct, "unsigned") {
		unsigned = true
	}
	var typ string
	var gtp int
	switch dt {
	case "bit":
		typ = "[]byte"
		gtp = bigtypeCompareBit
	case "bool", "boolean":
		typ = "bool"
	case "char", "varchar":
		typ = "string"
		gtp = bigtypeCompareString
	case "tinytext", "text", "mediumtext", "longtext", "json":
		typ = "string"
	case "tinyint":
		typ = "int8"
		if unsigned {
			typ = "uint8"
		}
		gtp = bigtypeCompare
	case "smallint":
		typ = "int16"
		if unsigned {
			typ = "uint16"
		}
		gtp = bigtypeCompare
	case "mediumint", "int", "integer":
		typ = "int32"
		if unsigned {
			typ = "uint32"
		}
		gtp = bigtypeCompare
	case "bigint":
		typ = "int64"
		if unsigned {
			typ = "uint64"
		}
		gtp = bigtypeCompare
	case "float":
		typ = "float32"
		gtp = bigtypeCompare
	case "decimal", "double":
		typ = "float64"
		gtp = bigtypeCompare
	case "binary", "varbinary":
		typ = "[]byte"
		gtp = bigtypeCompare
	case "tinyblob", "blob", "mediumblob", "longblob":
		typ = "[]byte"
	case "timestamp", "datetime", "date":
		typ = "time.Time"
		gtp = bigtypeCompareTime
	case "time", "year", "enum", "set":
		typ = "string"
		gtp = bigtypeCompare
	default:
		typ = "UNKNOWN"
	}
	return typ, gtp
}
