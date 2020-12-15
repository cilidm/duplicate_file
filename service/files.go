package service

import (
	"duplicate_file/util"
	"encoding/json"
	"fmt"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	FileID   uint64 = util.InitialID
	fileChan        = make(chan File, 10)
	maxRun          = make(chan bool, 5)
	done            = make(chan bool)
)

type File struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Dir     string `json:"dir"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
	Md      string `json:"md"`
	Status  uint8  `json:"status"`
}

func InitSearchFiles(dir string) bool {
	start := time.Now()
	if dir == "" {
		fmt.Println("no path")
		os.Exit(1)
	}
	go getFileFromChan()
	WalkDir(dir)
	fmt.Println("已完成检索，共发现", FileID-100000001, "个文件，耗时", time.Since(start), "，文件信息入库中，请稍后")
	select {
	case <-done:
		fmt.Println("信息更新完毕,耗时", time.Since(start))
	}

	// 查询重复文件
	SearchRepeat()

	os.Create("lock")
	logging.Info("执行完毕，耗时", time.Since(start))
	return true
}

func SearchRepeat() {
	var (
		repeIndex = util.InitialID // 重复文件索引
		indexNum  = util.InitialID // 全部文件索引
		picNum    = util.InitialID
		videoNum  = util.InitialID
	)
	files, err := levelDB.GetServer().FindByPrefix(util.DupKey)
	if err != nil {
		logging.Fatal(err.Error())
	}
	for _, v := range files {
		var file File
		json.Unmarshal([]byte(v), &file)

		levelDB.GetServer().Insert(util.GetDupIndexKey(gconv.String(indexNum)), file)
		atomic.AddUint64(&indexNum, 1)

		ext := path.Ext(file.Name) // 查找图片
		if str.IsContain(util.PicExt, ext) {
			levelDB.GetServer().Insert(util.GetPicIndexKey(gconv.String(picNum)), file)
			atomic.AddUint64(&picNum, 1)
		}
		if str.IsContain(util.VideoExt, ext) { // 查找视频
			levelDB.GetServer().Insert(util.GetVideoIndexKey(gconv.String(videoNum)), file)
			atomic.AddUint64(&videoNum, 1)
		}

		if file.Status == util.StatusRepeat {
			continue
		}
		// 完全匹配 查找是否在重复列表中
		has, err := levelDB.GetServer().IsExistFromLevelDB(util.GetRepeatKey(file.Md, file.Path))
		if err != nil {
			logging.Error(err.Error())
			continue
		}
		// 前面筛选过 还是非repeat status的 进行更新
		if has {
			file.Status = util.StatusRepeat
			levelDB.GetServer().Insert(util.GetDupKey(file.Md, file.Path), file)
			continue
		}
		fs, err := levelDB.GetServer().FindByPrefix(util.GetDupPreKey(file.Md))
		if err != nil {
			logging.Error(err.Error())
			continue
		}
		// 选择大于1个的
		if len(fs) > 1 {
			for _, val := range fs {
				var fal File
				json.Unmarshal([]byte(val), &fal)
				fmt.Println("发现重复文件", fal.Path)
				fal.Status = util.StatusRepeat
				levelDB.GetServer().Insert(util.GetRepeatKey(fal.Md, fal.Path), fal)             // 插入重复数据
				levelDB.GetServer().Insert(util.GetDupKey(fal.Md, fal.Path), fal)                // 更新全部数据里的status todo 在这里删除即可
				levelDB.GetServer().Insert(util.GetRepeatIndexKey(gconv.String(repeIndex)), fal) // 重复数据索引
				atomic.AddUint64(&repeIndex, 1)
			}
		}
	}
	fmt.Println("重复文件查找完毕")
}

func InsertLdb(file *File) {
	md, err := util.GetMdByPath(file.Path)
	if err != nil {
		<-maxRun
		logging.Error(err)
		return
	}
	file.Md = md
	err = levelDB.GetServer().Insert(util.GetDupKey(file.Md, file.Path), *file)
	if err != nil {
		<-maxRun
		logging.Error(err)
		return
	}
	err = levelDB.GetServer().Insert(util.GetDupSucKey(file.Path), *file)
	if err != nil {
		<-maxRun
		logging.Error(err)
		return
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
				go InsertLdb(&val)
			}
		}
	}
}

func WalkDir(dir string) {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if runtime.GOOS == "darwin" {
			if strings.HasSuffix(info.Name(), ".app") || strings.Contains(p, ".app") {
				return nil
			}
		}
		if info.Size() < 1 || info.IsDir() {
			return nil
		}
		fmt.Println("发现文件", p)
		key := util.GetDupSucKey(p) // 判断是否已完成的
		has, err := levelDB.GetServer().IsExistFromLevelDB(key)
		if err != nil {
			logging.Fatal(err)
			return err
		}
		if has == false {
			dir, _ := path.Split(p)
			fileChan <- File{
				ID:      FileID,
				Name:    info.Name(),
				Dir:     dir,
				Path:    p,
				Size:    info.Size(),
				Status:  util.StatusBegin,
				ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			}
			atomic.AddUint64(&FileID, 1)
		}
		return nil
	})
	defer close(fileChan)
	if err != nil {
		logging.Error(err)
	}
}
