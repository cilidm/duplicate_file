package service

import (
	"duplicate_file/app/db/dao"
	"duplicate_file/app/model"
	"duplicate_file/app/util"
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	FileID   uint64 = 1 // 计数器
	fileChan        = make(chan model.File, 10)
	maxRun          = make(chan bool, 5)
	done            = make(chan bool)
	checkMd  bool
)

func InitSearchFiles(dir string, check bool) bool {
	start := time.Now()
	if dir == "" {
		fmt.Println("路径不能为空")
		os.Exit(1)
	}
	checkMd = check
	go getFileFromChan()
	WalkDir(dir)
	fmt.Println("已完成检索，共发现", FileID, "个文件，耗时", time.Since(start), "，文件信息入库中，请稍后")
	select {
	case <-done:
		fmt.Println("信息更新完毕,耗时", time.Since(start))
	}
	os.Create("lock")
	logging.Info("执行完毕，耗时", time.Since(start))
	return true
}

func InsertFileInfo(file *model.File) {
	if checkMd {
		md, err := util.GetMdByPath(file.Path)
		if err != nil {
			logging.Error(err)
			<-maxRun
			return
		}
		file.Md = md
	}
	err := dao.Insert(*file)
	if err != nil {
		logging.Error(err.Error())
	}
	<-maxRun
}

func getFileFromChan() {
	for {
		val, ok := <-fileChan
		if !ok {
			if len(maxRun) == 0 {
				done <- true
				break
			}
		} else {
			if val.Md == "" {
				maxRun <- true
				go InsertFileInfo(&val)
			}
		}
	}
}

func WalkDir(dir string) {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Error(err.Error())
			return nil
		}
		if runtime.GOOS == "darwin" {
			if strings.HasSuffix(info.Name(), ".app") || strings.Contains(p, ".app") {
				return nil
			}
		}
		if info.IsDir() == false {
			if info.Size() > 0 {
				fmt.Println("发现文件", p)
				lev := strings.Split(p, "/")
				ext := path.Ext(p)
				if ext == "" {
					sp := strings.Split(lev[len(lev)-1], ".")
					if len(sp) > 1 {
						ext = "." + sp[len(sp)-1]
					}
				}
				has, _ := dao.FindByPath(p)
				if has.ID < 1 {
					dir, _ := path.Split(p)
					fileChan <- model.File{
						Name:      info.Name(),
						Dir:       dir,
						Path:      p,
						Size:      info.Size(),
						Status:    util.StatusBegin,
						Level:     uint8(len(lev)),
						ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
						Ext:       ext,
						CreatedAt: time.Now(),
					}
					atomic.AddUint64(&FileID, 1)
				}
			}
		}
		return nil
	})
	defer close(fileChan)
	if err != nil {
		logging.Error(err)
	}
}
