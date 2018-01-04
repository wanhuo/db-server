package config

import (
	"os"
	"fmt"
	"path/filepath"

	"github.com/jingwanglong/golog"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/utils"
)


var (
	// AppPath is the absolute path to the app
	AppPath string
	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "json"
	//default config file
	configFile = "prod.conf.json"
)


var(
	AppConfig    config.Configer
	Host         = make(map[string]interface{})
	DatabaseList  []map[string]interface{}
)


func parseConfig (){
	dbs, err := AppConfig.DIY("database")
	if err != nil {
		panic(err)
	}
	DBList := dbs.([]interface{})
	if DBList == nil {
		panic("db not []interface{}")
	}
	for _, conf := range DBList {
		confCasted := conf.(map[string]interface{})
		if confCasted == nil {
			panic("conf is not a map")
		}
		DatabaseList = append(DatabaseList, confCasted)
	}

	value, err := AppConfig.DIY("host")
	if err != nil{
		panic(err)
	}
	host := value.(map[string]interface{})
	if host == nil {
		panic("get host conf error")
	}
	Host = host
}

func init()  {
	workPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	logPath := filepath.Join(workPath, "logs", "db.log")
	err = golog.SetOutputLogger("*", logPath)
	err = golog.SetLevelByString("*", "info")
	if err != nil {
		panic(err)
	}

	appConfigPath = filepath.Join(workPath, "config", configFile)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "config", configFile)
		if !utils.FileExists(appConfigPath) {
			fmt.Printf("work path is %s\n", workPath)
			panic("file is not exist!")
			os.Exit(1)
		}
	}
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		panic(err)
	}
	AppConfig = ac
	parseConfig()
}

