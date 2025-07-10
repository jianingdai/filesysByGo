package endpoint

import (
	"08/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateFolder(c *gin.Context) {
	parentIDStr := c.Param("file_id")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil || parentID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 parent_id"})
		return
	}

	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	folder, err := service.CreateFolderService.Create(c, req.Name, uint(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, folder)
}
