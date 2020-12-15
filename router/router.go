package router

import (
	"duplicate_file/api"
	"duplicate_file/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func Server() {
	r := gin.Default()

	r.Static("/static", "static")
	r.LoadHTMLGlob("views/*")

	r.GET("/file", func(c *gin.Context) {
		c.HTML(http.StatusOK, "file_manager.html", gin.H{})
	})

	r.POST("/file_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		resp, count := service.GetFileByPage(page, limit)
		c.JSON(http.StatusOK, gin.H{"count": count, "images": resp})
	})

	r.GET("show", func(c *gin.Context) {
		imageName := c.Query("imageName")
		c.File(imageName)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "start.html", gin.H{})
	})

	r.GET("/main", func(c *gin.Context) { // 重复文件
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	r.GET("/json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		resp, count := service.GetRepeatByPage(page, limit)
		var data []service.File
		for _, v := range resp {
			data = append(data, v.(service.File))
		}
		api.SuccessResp(c).SetCode(0).SetData(data).SetCount(count).WriteJsonExit()
	})

	r.POST("/del", func(c *gin.Context) {
		p := c.PostForm("path")
		if p == "" {
			api.ErrorResp(c).SetMsg("路径不能为空").WriteJsonExit()
			return
		}
		err := os.Remove(p)
		if err != nil {
			api.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
			return
		}
		api.SuccessResp(c).WriteJsonExit()
	})

	r.GET("/index", func(c *gin.Context) { // 全部文件
		c.HTML(http.StatusOK, "all.html", gin.H{})
	})

	r.GET("/index_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		resp, count := service.GetDupByPage(page, limit)
		var data []service.File
		for _, v := range resp {
			data = append(data, v.(service.File))
		}
		api.SuccessResp(c).SetCode(0).SetData(data).SetCount(count).WriteJsonExit()
	})

	r.GET("/pic", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pic.html", gin.H{})
	})

	r.GET("/pic_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		resp, count := service.GetPicByPage(page, limit)
		var data []service.File
		for _, v := range resp {
			data = append(data, v.(service.File))
		}
		api.SuccessResp(c).SetCode(0).SetData(data).SetCount(count).WriteJsonExit()
	})

	r.GET("/video", func(c *gin.Context) {
		c.HTML(http.StatusOK, "video.html", gin.H{})
	})

	r.GET("/video_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		resp, count := service.GetVideoByPage(page, limit)
		var data []service.File
		for _, v := range resp {
			data = append(data, v.(service.File))
		}
		api.SuccessResp(c).SetCode(0).SetData(data).SetCount(count).WriteJsonExit()
	})

	r.Run(":8000")
}
