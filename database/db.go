package database

import (
	"fmt"
	"sse/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

func init() {
	util.MyLogger.Info("begin init sse")
	dbInfo := viper.GetStringMapString("db")
	if len(dbInfo) == 0 {
		return
	}
	host := viper.GetString("postgres.host")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")
	dbName := viper.GetString("postgres.dbName")
	port := viper.GetString("postgres.port")

	postgresConfig := fmt.Sprintf(
		//"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	// util.MyLogger.Info(postgresConfig)
	dbInstance, err := gorm.Open("postgres", postgresConfig)
	if err != nil {
		util.MyLogger.Error(err.Error())
		panic("DB Error")
	}
	util.MyLogger.Info("connect success")
	//db.SetLogger(util.MyLogger)
	dbInstance.Callback().Update().Remove("gorm:update_time_stamp")
	db = dbInstance
	util.MyLogger.Info("db init complete")
}

func GetDatabase() *gorm.DB {
	return db
}
