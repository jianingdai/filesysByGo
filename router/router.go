package router

import (
	"08/endpoint"
	"08/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/web", "frontend")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/web/index.html")
	})

	r.POST("/login", endpoint.Login)

	auth := r.Group("/api", middleware.AuthMiddleware())
	auth.POST("/user/create", endpoint.CreateUser)

	fileGroup := r.Group("/api/file", middleware.AuthMiddleware(), middleware.FilePermissionMiddleware())
	// 添加创建文件夹的路由
	fileGroup.POST("/:file_id/new", endpoint.CreateFolder)
	// 添加上传文件的路由
	fileGroup.POST("/:file_id/upload", endpoint.UploadFile)
	// 添加删除文件或文件夹的路由
	fileGroup.DELETE("/:file_id", endpoint.DeleteFileOrFolder)
	// 添加复制文件或文件夹的路由
	fileGroup.POST("/:file_id/copy", endpoint.CopyFileOrFolder)
	// 添加移动文件或文件夹的路由
	fileGroup.POST("/:file_id/move", endpoint.MoveFileOrFolder)
	// 添加获取文件夹列表的路由
	fileGroup.GET("/:file_id/list", endpoint.ListFiles)
	// 添加重命名文件或文件夹的路由
	fileGroup.POST("/:file_id/rename", endpoint.RenameFileOrFolder)
	// 添加下载文件内容的路由
	fileGroup.GET("/:file_id/content", endpoint.DownloadFileContent)
	// 添加获取文件或文件夹信息的路由
	fileGroup.GET("/:file_id", endpoint.GetFileInfo)
	// 添加更新文件的路由
	fileGroup.POST("/:file_id/update", endpoint.UpdateFile)
	// 添加下载历史版本文件内容的路由
	fileGroup.GET("/:file_id/version/:ver_num/content", endpoint.DownloadFileVersionContent)
	// 添加获取文件版本历史的路由
	fileGroup.GET("/:file_id/versions", endpoint.GetFileVersions)

	return r
}
