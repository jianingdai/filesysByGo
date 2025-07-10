package endpoint

import (
	"08/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteFileOrFolder(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	userID := c.GetInt("user_id")
	err = service.DeleteFileService.Delete(c, uint(fileID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
