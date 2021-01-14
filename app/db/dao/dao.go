package dao

import (
	"duplicate_file/app/db"
	"duplicate_file/app/model"
	"github.com/cilidm/toolbox/gconv"
	"strings"
)

func Insert(f model.File) error {
	err := db.DB.Create(&f).Error
	return err
}

func Update(id uint64, attr map[string]interface{}) error {
	err := db.DB.Where("id = ?", id).Update(attr).Error
	return err
}

func FindByPath(p string) (f model.File, err error) {
	err = db.DB.Model(model.File{}).Where("path = ?", p).First(&f).Error
	return f, err
}

func FindByPage(pageNum, limit int, filters ...interface{}) (files []model.File, count int, err error) {
	offset := (pageNum - 1) * limit
	var queryArr []string
	var values []interface{}
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			queryArr = append(queryArr, gconv.String(filters[k]))
			values = append(values, filters[k+1])
		}
	}
	query := db.DB.Model(model.File{})
	query.Where(strings.Join(queryArr, " OR "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " OR "), values...).
		Order("id desc").Limit(limit).Offset(offset).Find(&files).Error
	return
}

func FindByID(id int) (f model.File, err error) {
	err = db.DB.Where("id = ?", id).First(&f).Error
	return
}

func FindByRaw(sql string) (f []model.File, err error) {
	err = db.DB.Raw(sql).Find(&f).Error
	return
}
