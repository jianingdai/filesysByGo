package endpoint

import (
	"08/service"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 更新文件内容
func UpdateFile(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	// 读取请求体中的文件内容
	fileContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取文件内容失败"})
		return
	}

	userID := c.GetInt("user_id")
	updatedFile, err := service.UpdateFileService.Update(c, uint(fileID), uint(userID), fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFile)
}
