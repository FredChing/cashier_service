package lib

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Mysql *sql.DB
	DbMap *gorp.DbMap
)

func MysqlInit() (err error) {
	fmt.Println("Connect to mysql...")
	dsn := beego.AppConfig.String("mysql::dsn")
	Mysql, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("connect to mysql failed")
	} else {
		fmt.Println("Ok")
	}

	maxIdleNum, _ := beego.AppConfig.Int("mysql::maxIdleConns")
	maxOpenNum, _ := beego.AppConfig.Int("mysql::maxOpenConns")
	if maxIdleNum != 0 {
		Mysql.SetMaxIdleConns(maxIdleNum)
	}
	if maxOpenNum != 0 {
		Mysql.SetMaxOpenConns(maxOpenNum)
	}

	err = Mysql.Ping()
	if err != nil {
		panic("Mysql connect failed")
	}
	c := make(chan error)
	go MonitorMysql(c)

	DbMap = &gorp.DbMap{Db: Mysql, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	return
}

func MonitorMysql(c chan error) {
	beat := beego.AppConfig.String("mysql::maxOpenConns")
	d, err := time.ParseDuration(beat)
	var chant <-chan time.Time
	if err != nil {
		chant = time.Tick(30 * time.Second)
	} else {
		chant = time.Tick(d)
	}
	for _ = range chant {
		err = Mysql.Ping()
		if err != nil {
			fmt.Println("mysql connect lost")
		}
	}
}

func Insert(object interface{}) (result sql.Result, err error) {
	/*
		val := reflect.ValueOf(object)
		table_name := getTableName(val)
		val_elem := val.Elem()
		t := reflect.TypeOf(object).Elem()

		sql := ""
		columns_str := ""
		values_str := ""
		values := make([]interface{}, 0)
		for i := 0; i < t.NumField(); i++ {
			type_field := t.Field(i)
			if type_field.Tag.Get("skip") != "true" {
				switch type_field.Tag.Get("auto") {
				case "true":
				default:
					columns_str = columns_str + "," + strings.ToLower(type_field.Name)
					values_str += ",?"
					values = append(values, val_elem.Field(i).Interface())
				}
			}

		}
		columns_str = strings.TrimPrefix(columns_str, ",")
		values_str = strings.TrimPrefix(values_str, ",")
		sql = "INSERT INTO " + table_name + "(" + columns_str + ") VALUES(" + values_str + ")"
	*/

	sql, values := getInsertSql(object)
	stmt, err := Mysql.Prepare(sql)
	if err != nil {
		log.Printf("sql prepare failed:%v, sql::%v", err, sql)
		return
	}
	defer stmt.Close()

	result, err = stmt.Exec(values...)
	if err != nil {
		log.Printf("%v save error:%v", sql, err)
	}
	return result, err
}
func getInsertSql(object interface{}) (sql string, values []interface{}) {
	val := reflect.ValueOf(object)
	table_name := getTableName(val)
	val_elem := val.Elem()
	t := reflect.TypeOf(object).Elem()

	sql = ""
	columns_str := ""
	values_str := ""
	values = make([]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		type_field := t.Field(i)
		if type_field.Tag.Get("skip") != "true" {
			switch type_field.Tag.Get("auto") {
			case "true":
			default:
				columns_str = columns_str + "," + strings.ToLower(type_field.Name)
				values_str += ",?"
				values = append(values, val_elem.Field(i).Interface())
			}
		}

	}
	columns_str = strings.TrimPrefix(columns_str, ",")
	values_str = strings.TrimPrefix(values_str, ",")
	sql = "INSERT INTO " + table_name + "(" + columns_str + ") VALUES(" + values_str + ")"
	return
}

// get table name. method, or field name.
func getTableName(val reflect.Value) string {
	ind := reflect.Indirect(val)
	fun := val.MethodByName("TableName")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}
	return ind.Type().Name()
}