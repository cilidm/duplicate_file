package db

import (
	"duplicate_file/app/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB  *gorm.DB
	err error
	// 数据库配置
	User   = "root"
	Pwd    = "123456"
	Host   = "127.0.0.1"
	Port   = "3306"
	DBName = "duplicate_file"
)

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		User, Pwd, Host, Port, DBName)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}
	DB.SingularTable(true)
	DB.LogMode(true)
	DB.AutoMigrate(model.File{})
	return
}

func DropDB() {
	DB.Delete(model.File{})
}

func CloseMysql() {
	if err := DB.Close(); err != nil {
		panic("关闭数据库失败:" + err.Error())
	}
}
