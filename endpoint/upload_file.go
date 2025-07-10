package endpoint

import (
	"08/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	parentIDStr := c.Param("file_id")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil || parentID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}

	// 这里假设 service 层有处理文件保存的逻辑
	userID := c.GetInt("user_id")
	savedFile, err := service.UploadFileService.Upload(c, file, uint(parentID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, savedFile)
}
