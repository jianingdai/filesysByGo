package endpoint

import (
	"08/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MoveFileRequest struct {
	TargetParentID int `json:"target_parent_id" binding:"required"`
}

func MoveFileOrFolder(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	var req MoveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID := c.GetInt("user_id")
	movedFile, err := service.MoveFileService.Move(c, uint(fileID), uint(req.TargetParentID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移动失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, movedFile)
}
