package service

import (
	"duplicate_file/app/db/dao"
	"duplicate_file/app/model"
	"duplicate_file/app/util"
	"fmt"
	"github.com/cilidm/toolbox/OS"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
	"path"
	"strings"
)

func GetRepeatFile(page, limit int) ([]model.File, int, error) {
	offset := (page - 1) * limit
	sql := "SELECT * FROM file WHERE md IN ( SELECT md FROM file WHERE md != '' GROUP BY md HAVING COUNT(md) > 1 ) ORDER BY md LIMIT " + gconv.String(limit) + " offset " + gconv.String(offset)
	fmt.Println(sql)
	f, err := dao.FindByRaw(sql)
	countSql := "SELECT id FROM file WHERE md IN ( SELECT md FROM file WHERE md != '' GROUP BY md HAVING COUNT(md) > 1 )"
	count, err := dao.FindByRaw(countSql)
	return f, len(count), err
}

func GetFileInfoAll(page, limit int) (files []model.File, count int, err error) {
	files, count, err = dao.FindByPage(page, limit)
	return
}

func GetFileInfo(page, limit int, filters []interface{}) (resp []model.File, count int, err error) {
	var files []model.File
	files, count, err = dao.FindByPage(page, limit, filters...)
	for _, v := range files {
		if OS.IsWindows() {
			v.Dir = strings.TrimRight(v.Path, v.Name)
		}
		v.Path = "/detail?index=" + gconv.String(v.ID) // 用索引替换path
		resp = append(resp, v)
	}
	return
}

func GetFileShowPage(page, limit int) (resp []model.FileManager, count int, err error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	var (
		files   []model.File
		filters []interface{}
	)
	filters = append(filters, "ext in (?)", util.PicExt)
	filters = append(filters, "ext in (?)", util.VideoExt)
	files, count, err = dao.FindByPage(page, limit, filters...)
	for _, v := range files {
		var thumb string
		if str.IsContain(util.PicExt, path.Ext(v.Name)) || str.IsContain(util.VideoExt, path.Ext(v.Name)) {
			thumb = "/show?imageName=" + v.Path
		} else {
			thumb = v.Path
		}
		ext := path.Ext(v.Name)
		if ext != "" {
			// path.Ext获取的是.jpg这样的样式，前端需要jpg这样的样式，因此去掉最前面
			ext = strings.TrimLeft(ext, ".")
		}
		resp = append(resp, model.FileManager{
			Thumb: thumb,
			Name:  v.Name,
			Type:  ext,
			Path:  v.Path,
			Index: "/detail?index=" + gconv.String(v.ID),
		})
	}
	return
}

func GetByID(id string) (f model.File) {
	var err error
	f, err = dao.FindByID(gconv.Int(id))
	if err != nil {
		logging.Error(err.Error())
	}
	return f
}

//func GetPicInfo(page, limit int) (files []model.File, count int, err error)  {
//	files,count,err = dao.FindByPage()
//}
