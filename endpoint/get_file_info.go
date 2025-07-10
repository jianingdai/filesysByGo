package endpoint

import (
	"08/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取文件或文件夹信息
func GetFileInfo(c *gin.Context) {
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

	userID := c.GetInt("user_id")
	if uint(file.UserID) != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	c.JSON(http.StatusOK, file)
}
