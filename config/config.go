package config

import (
	"os"
	"path/filepath"

	"github.com/astaxie/beego/config"
	"github.com/davyxu/golog"
	"github.com/astaxie/beego/utils"
)

var log *golog.Logger = golog.New("config")

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
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath = filepath.Join(workPath, "config", configFile)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "config", configFile)
		if !utils.FileExists(appConfigPath) {
			log.Errorf("file is not exist!")
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

