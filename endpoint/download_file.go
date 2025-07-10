package endpoint

import (
	"08/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 下载当前版本文件内容
func DownloadFileContent(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	if file.Type == "folder" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能下载文件夹"})
		return
	}
	if file.StoreKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件存储信息缺失"})
		return
	}

	c.FileAttachment("data/"+file.StoreKey, file.Name)
}

// 下载指定版本的文件内容
func DownloadFileVersionContent(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	verNumStr := c.Param("ver_num")
	verNum, err := strconv.Atoi(verNumStr)
	if err != nil || verNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本号"})
		return
	}

	// 查询文件信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	userID := c.GetInt("user_id")
	if uint(file.UserID) != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	if file.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能下载文件夹"})
		return
	}

	// 如果请求的是当前版本，直接返回当前文件
	if int32(verNum) == file.VerNum {
		if file.StoreKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件存储信息缺失"})
			return
		}
		c.FileAttachment("data/"+file.StoreKey, file.Name)
		return
	}

	// 查询历史版本信息
	version, err := dao.Q.Version.WithContext(c).Where(
		dao.Version.FileID.Eq(int32(fileID)),
		dao.Version.VerNum.Eq(int32(verNum)),
	).First()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "版本不存在"})
		return
	}

	if version.StoreKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "版本文件存储信息缺失"})
		return
	}

	c.FileAttachment("data/"+version.StoreKey, file.Name)
}
