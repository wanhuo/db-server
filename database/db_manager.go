package database

import (
	"fmt"
	"strconv"
	"db-server/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"db-server/proto/dbproto"
)

//connect to mysql
func init() {

	orm.RegisterDriver("mysql", orm.DRMySQL)
	for _, DbConfig := range config.DatabaseList{
		dbID     := DbConfig["id"].(string)
		dbName   := DbConfig["name"].(string)
		user     := DbConfig["user"].(string)
		host     := DbConfig["host"].(string)
		port     := int(DbConfig["port"].(float64))
		password := DbConfig["password"].(string)
		dbPool   := int(DbConfig["pool"].(float64))

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, dbName)
		if dbID == "ms_db" {
			orm.RegisterDataBase("default", "mysql", dsn, dbPool, dbPool)
		}else{
			orm.RegisterDataBase(dbName, "mysql", dsn, dbPool, dbPool)
		}

	}
}

type SqlArgList []string
type runSqlFunc   func(SqlArgList) error
type querySqlFunc func(SqlArgList) ([]*dbproto.OneRow, error)

var (
	runSqlHandler   = make(map[string]runSqlFunc)
	querySqlHandler = make(map[string]querySqlFunc)
)

func RegisterRunSqlCB(action string, sqlHandler runSqlFunc){
	if f := runSqlHandler[action]; f != nil{
		return
	}
	runSqlHandler[action] = sqlHandler
}

func RegisterQueryCB(action string, sqlHandler querySqlFunc){
	if f := querySqlHandler[action]; f != nil{
		return
	}
	querySqlHandler[action] = sqlHandler
}

func ProcessQuerySql(action string, params []*dbproto.OneField) (rows []*dbproto.OneRow, err error) {
	handler := querySqlHandler[action]
	if handler == nil{
		err = fmt.Errorf("can not find %s handler", action)
		return
	}
	var paramList SqlArgList
	for _, field := range params {
		log.Debugf("field: %s ", string(field.Value))
		paramList = append(paramList, string(field.Value))
	}
	rows, err = handler(paramList)
	return
}

func ProcessRunSql(action string, params []*dbproto.OneField) (err error) {
	handler := runSqlHandler[action]
	if handler == nil{
		err = fmt.Errorf("can not find %s handler", action)
		return
	}
	var paramList SqlArgList
	for _, field := range params {
		log.Debugf("field: %s ", string(field.Value))
		paramList = append(paramList, string(field.Value))
	}
	err = handler(paramList)
	return
}

func serializeRowDate(params []interface{}) (row *dbproto.OneRow){
	var fieldList []*dbproto.OneField
	for _, field := range params{
		var valueStr string
		switch field.(type) {
		case int:
			value := field.(int)
			valueStr = strconv.Itoa(int(value))
		case int32:
			value := field.(int32)
			valueStr = strconv.Itoa(int(value))
		case int64:
			value := field.(int64)
			valueStr = strconv.Itoa(int(value))
		case string:
			valueStr = field.(string)
		}
		field := dbproto.OneField{
			Value: []byte(valueStr),
		}
		fieldList = append(fieldList, &field)
	}
	row = &dbproto.OneRow{
		OneField:fieldList,
	}
	return
}
