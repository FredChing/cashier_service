package lib

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type CTx struct {
	Db *sql.DB
	*sql.Tx
}

func NewTx() (*CTx, error) {
	tx, err := Mysql.Begin()
	ctx := &CTx{Mysql, tx}
	return ctx, err
}

func (ctx *CTx) Insert(object interface{}) (result sql.Result, err error) {
	sql, values := getInsertSql(object)
	result, err = ctx.Exec(sql, values...)
	return
}
