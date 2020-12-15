package main

import (
	"duplicate_file/router"
	"duplicate_file/service"
	"duplicate_file/util"
	"flag"
	"github.com/cilidm/toolbox/levelDB"
	"os"
)

var (
	dir   string
	isNew bool
)

func init() {
	levelDB.InitServer("runtime/")
	flag.StringVar(&dir, "p", "", "查找路径")
	flag.BoolVar(&isNew, "n", false, "是否新开启任务")
	flag.Parse()
	util.CheckTemplate()
}

func main() {
	if isNew && util.CheckFileIsExist("lock") {
		os.Remove("lock")
	}
	if util.CheckFileIsExist("lock") == false {
		if service.InitSearchFiles(dir) == false {
			return
		}
	}
	router.Server()
}
