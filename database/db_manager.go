package database

import (
	"fmt"
	"db-server/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/davyxu/golog"
	"db-server/proto/dbproto"
)
var log *golog.Logger = golog.New("Database")
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
			log.Debugf("default database: ms db")
			orm.RegisterDataBase("default", "mysql", dsn, dbPool, dbPool)
		}else{
			orm.RegisterDataBase(dbName, "mysql", dsn, dbPool, dbPool)
		}

	}
}

type sqlArgs []interface{}
//type runSqlFunc   func(sqlArgs) error
type querySqlFunc func(sqlArgs) ([]*dbproto.OneRow, error)

var (
	//runSqlHandler   = make(map[string]runSqlFunc)
	querySqlHandler = make(map[string]querySqlFunc)
)

//func RegisterRunSqlCB(action string, sqlHandler runSqlFunc){
//	if f := runSqlHandler[action]; f != nil{
//		return
//	}
//	runSqlHandler[action] = sqlHandler
//}

func RegisterQueryCB(action string, sqlHandler querySqlFunc){
	if f := querySqlHandler[action]; f != nil{
		return
	}
	querySqlHandler[action] = sqlHandler
}

func GetQueryHandler(action string) querySqlFunc {
	if f := querySqlHandler[action]; f != nil{
		return querySqlHandler[action]
	}
	return nil
}





