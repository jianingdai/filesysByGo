package endpoint

import (
	"08/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RenameFileRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

func RenameFileOrFolder(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	var req RenameFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID := c.GetInt("user_id")
	file, err := service.RenameFileService.Rename(c, uint(fileID), uint(userID), req.NewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重命名失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, file)
}
