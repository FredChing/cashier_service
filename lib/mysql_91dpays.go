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
	Mysql_91dpays *sql.DB
	DbMap_91dpays *gorp.DbMap
)

func Mysql91dpaysInit() (err error) {
	fmt.Println("Connect to mysql_91dpays...")
	dsn := beego.AppConfig.String("mysql_91dpays::dsn")
	Mysql_91dpays, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("connect to Mysql_91dpays failed")
	} else {
		fmt.Println("Ok")
	}

	maxIdleNum, _ := beego.AppConfig.Int("mysql_91dpays::maxIdleConns")
	maxOpenNum, _ := beego.AppConfig.Int("mysql_91dpays::maxOpenConns")
	if maxIdleNum != 0 {
		Mysql_91dpays.SetMaxIdleConns(maxIdleNum)
	}
	if maxOpenNum != 0 {
		Mysql_91dpays.SetMaxOpenConns(maxOpenNum)
	}

	err = Mysql_91dpays.Ping()
	if err != nil {
		panic("Mysql_91dpays connect failed")
	}
	c := make(chan error)
	go MonitorMysql_91dpays(c)

	DbMap_91dpays = &gorp.DbMap{Db: Mysql_91dpays, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	return
}

func MonitorMysql_91dpays(c chan error) {
	beat := beego.AppConfig.String("mysql_91dpays::maxOpenConns")
	d, err := time.ParseDuration(beat)
	var chant <-chan time.Time
	if err != nil {
		chant = time.Tick(30 * time.Second)
	} else {
		chant = time.Tick(d)
	}
	for _ = range chant {
		err = Mysql_91dpays.Ping()
		if err != nil {
			fmt.Println("Mysql_91dpays connect lost")
		}
	}
}

func Insert_91dpays(object interface{}) (result sql.Result, err error) {
	/*
		val := reflect.ValueOf(object)
		table_name := getTableNamekKngaroo(val)
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
	stmt, err := Mysql_91dpays.Prepare(sql)
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
func getInsertSql91dpays(object interface{}) (sql string, values []interface{}) {
	val := reflect.ValueOf(object)
	table_name := getTableName91dpays(val)
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
func getTableName91dpays(val reflect.Value) string {
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

//func (u *User) Update(p map[string]interface{}, condition map[string]interface{}) (err error) {
//	sql := "UPDATE user SET "
//	for k, v := range p {
//		sql = sql + k + "=?,"
//	}
//
//	for k, v := range condition {
//		switch v.(type) {
//		case map[string]interface{}:
//			for condtion, value := range v {
//			}
//
//		}
//	}
//	//	sql := fmt.Sprintf("UPDATE user SET ")
//	return
//}
