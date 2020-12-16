package service

import (
	"duplicate_file/util"
	"encoding/json"
	"github.com/cilidm/toolbox/OS"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
	"path"
	"strings"
)

func GetDupByPage(page, limit int) ([]interface{}, int) {
	count, _ := levelDB.GetServer().FindByPrefix(util.DupIndexKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit) // InitialID - 1
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count))
	}
	datas, _ := levelDB.GetServer().FindLimit(util.DupIndexKey+":"+gconv.String(begin), util.DupIndexKey+":"+gconv.String(end))
	var resp []interface{}
	for _, v := range datas {
		var f File
		json.Unmarshal([]byte(v), &f)
		if OS.IsWindows() {
			f.Dir = strings.TrimRight(f.Path, f.Name)
		}
		resp = append(resp, f)
	}
	util.SortBodyByMd(resp)
	return resp, len(count)
}

func GetRepeatByPage(page, limit int) ([]interface{}, int) {
	count, _ := levelDB.GetServer().FindByPrefix(util.RepeatIndexKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit)
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count)) // 因取值是[)的 因此结尾要多1
	}
	datas, _ := levelDB.GetServer().FindLimit(util.RepeatIndexKey+":"+gconv.String(begin), util.RepeatIndexKey+":"+gconv.String(end))
	var resp []interface{}
	for _, v := range datas {
		var f File
		json.Unmarshal([]byte(v), &f)
		if OS.IsWindows() {
			f.Dir = strings.TrimRight(f.Path, f.Name)
		}
		resp = append(resp, f)
	}
	util.SortBodyByMd(resp)
	return resp, len(count)
}

func GetPicByPage(page, limit int) ([]interface{}, int) {
	count, _ := levelDB.GetServer().FindByPrefix(util.PicKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit)
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count)) // 因取值是[)的 因此结尾要多1
	}
	datas, _ := levelDB.GetServer().FindLimit(util.PicKey+":"+gconv.String(begin), util.PicKey+":"+gconv.String(end))
	var resp []interface{}
	for k, v := range datas {
		var f File
		json.Unmarshal([]byte(v), &f)
		if OS.IsWindows() {
			f.Dir = strings.TrimRight(f.Path, f.Name)
		}
		f.Path = "/detail?index=" + k // 用索引替换path
		resp = append(resp, f)
	}
	util.SortBodyByMd(resp)
	return resp, len(count)
}

func GetVideoByPage(page, limit int) ([]interface{}, int) {
	count, _ := levelDB.GetServer().FindByPrefix(util.VideoKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit)
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count)) // 因取值是[)的 因此结尾要多1
	}
	datas, _ := levelDB.GetServer().FindLimit(util.VideoKey+":"+gconv.String(begin), util.VideoKey+":"+gconv.String(end))
	var resp []interface{}
	for k, v := range datas {
		var f File
		json.Unmarshal([]byte(v), &f)
		if OS.IsWindows() {
			f.Dir = strings.TrimRight(f.Path, f.Name)
		}
		f.Path = "video_play?index=" + k
		resp = append(resp, f)
	}
	util.SortBodyByMd(resp)
	return resp, len(count)
}

type FileManager struct {
	Thumb string `json:"thumb"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Path  string `json:"path"`
	Index string `json:"index"`
}

func GetFileByPage(page, limit int) ([]FileManager, int) {
	//count, _ := levelDB.GetServer().FindByPrefix(util.DupIndexKey)
	count, _ := levelDB.GetServer().FindByPrefix(util.PicKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit)
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count)) // 因取值是[)的 因此结尾要多1
	}
	//datas, err := levelDB.GetServer().FindLimit(util.DupIndexKey+":"+gconv.String(begin), util.DupIndexKey+":"+gconv.String(end))
	datas, err := levelDB.GetServer().FindLimit(util.PicKey+":"+gconv.String(begin), util.PicKey+":"+gconv.String(end))
	if err != nil {
		logging.Error(err.Error())
	}
	var resp []FileManager
	for k, v := range datas {
		var (
			f     File
			thumb string
		)
		json.Unmarshal([]byte(v), &f)
		if str.IsContain(util.PicExt, path.Ext(f.Name)) || str.IsContain(util.VideoExt, path.Ext(f.Name)) {
			thumb = "/show?imageName=" + f.Path
		} else {
			thumb = f.Path
		}
		ext := path.Ext(f.Name)
		if ext != "" {
			// path.Ext获取的是.jpg这样的样式，前端需要jpg这样的样式，因此去掉最前面
			ext = strings.TrimLeft(ext, ".")
		}
		resp = append(resp, FileManager{
			Thumb: thumb,
			Name:  f.Name,
			Type:  ext,
			Path:  f.Path,
			Index: "/detail?index=" + k,
		})
	}
	return resp, len(count)
}

// 按长度截取字符串 超出的部分用..替代 （发现直接用css样式控制更简单，此方法弃用）
func OmitStr(str string, num int) (newStr string) {
	runeStr := []rune(str)
	if len(runeStr) > num {
		newStr = string(runeStr[:num]) + ".."
	} else {
		newStr = str
	}
	return
}

func GetIndexDetail(key string) (f File) {
	data, err := levelDB.GetServer().FindByKey(key)
	if err != nil {
		logging.Error(err.Error())
		return f
	}
	err = json.Unmarshal(data, &f)
	if err != nil {
		logging.Error(err.Error())
		return f
	}
	return f
}
