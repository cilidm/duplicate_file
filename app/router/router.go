package router

import (
	"duplicate_file/app/api"
	"duplicate_file/app/service"
	"duplicate_file/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func RunServer() {
	r := gin.Default()

	r.Static("/static", "static")
	r.LoadHTMLGlob("views/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "start.html", gin.H{})
	})

	r.GET("/main", func(c *gin.Context) { // 重复文件
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	r.GET("/json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		resp, count, err := service.GetRepeatFile(page, limit)
		if err != nil {
			api.ErrorResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
			return
		}
		api.SuccessResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
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
		resp, count, err := service.GetFileInfoAll(page, limit)
		if err != nil {
			api.ErrorResp(c).SetCode(0).SetData(resp).SetCount(0).WriteJsonExit()
			return
		}
		api.SuccessResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
	})

	r.GET("/pic", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pic.html", gin.H{})
	})

	r.GET("/pic_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		filters := make([]interface{}, 0)
		filters = append(filters, "ext in (?)", util.PicExt)
		resp, count, err := service.GetFileInfo(page, limit, filters)
		if err != nil {
			api.ErrorResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
			return
		}
		api.SuccessResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
	})

	r.GET("/video", func(c *gin.Context) {
		c.HTML(http.StatusOK, "video.html", gin.H{})
	})

	r.GET("/video_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		filters := make([]interface{}, 0)
		filters = append(filters, "ext in (?)", util.VideoExt)
		resp, count, err := service.GetFileInfo(page, limit, filters)
		if err != nil {
			api.ErrorResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
			return
		}
		api.SuccessResp(c).SetCode(0).SetData(resp).SetCount(count).WriteJsonExit()
	})

	// 图片文件展示
	r.GET("/file", func(c *gin.Context) {
		c.HTML(http.StatusOK, "file_manager.html", gin.H{})
	})

	r.POST("/file_json", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		resp, count, err := service.GetFileShowPage(page, limit)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"count": count, "images": resp})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count, "images": resp})
	})

	// 显示硬盘图片
	r.GET("show", func(c *gin.Context) {
		imageName := c.Query("imageName")
		c.File(imageName)
	})

	r.GET("detail", func(c *gin.Context) {
		index := c.Query("index")
		f := service.GetByID(index)
		c.File(f.Path)
	})

	r.GET("video_play", func(c *gin.Context) {
		index := c.Query("index")
		f := service.GetByID(index)
		f.Path = "/play?name=" + f.Path // 路径问题，需要使用c.File转换一下
		c.HTML(http.StatusOK, "video_play.html", f)
	})

	r.GET("play", func(c *gin.Context) {
		name := c.Query("name")
		c.File(name)
	})

	r.Run(":8000")
}
