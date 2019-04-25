package lib

import (
	"database/sql"
	logs "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type CTx91dpays struct {
	Db *sql.DB
	*sql.Tx
}

func NewTx91dpays() (*CTx91dpays, error) {
	tx, err := Mysql_91dpays.Begin()
	ctx := &CTx91dpays{Mysql_91dpays, tx}
	return ctx, err
}

func (ctx *CTx91dpays) Insert91dpays(object interface{}) (result sql.Result, err error) {
	sql, values := getInsertSql91dpays(object)
	logs.Infof("Insert91dpays::sql:%s",sql)
	result, err = ctx.Exec(sql, values...)
	return
}
