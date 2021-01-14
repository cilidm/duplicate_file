package main

import (
	"duplicate_file/app/db"
	"duplicate_file/app/router"
	"duplicate_file/app/service"
	"duplicate_file/app/util"
	"flag"
	"os"
)

var (
	dir   string
	isNew bool
	md    bool
)

func init() {
	db.InitMysql()
	flag.BoolVar(&isNew, "n", false, "是否创建新任务")
	flag.StringVar(&dir, "p", "", "查找路径")
	flag.BoolVar(&md, "md", true, "是否校验md")
	flag.Parse()
}

func main() {
	if isNew && util.CheckFileIsExist(util.LockKey) {
		os.Remove(util.LockKey)
	}
	if util.CheckFileIsExist(util.LockKey) == false {
		if service.InitSearchFiles(dir, md) == false {
			return
		}
	}
	router.RunServer()
}
