package middleware

import (
	"08/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FilePermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileIDStr := c.Param("file_id")
		fileID, err := strconv.Atoi(fileIDStr)
		if err != nil || fileID < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
			c.Abort()
			return
		}
		if fileID == 0 {
			// 如果是根目录，直接放行
			c.Next()
			return
		}

		userID := c.GetInt("user_id") // 假设已通过 AuthMiddleware 写入 user_id
		//根据file_id查询文件信息
		file, err := dao.Q.File.WithContext(c.Request.Context()).Where(dao.File.ID.Eq(int32(fileID))).First()
		if err != nil {
			// 假设使用 GORM，检查是否为记录未找到错误
			if err == gorm.ErrRecordNotFound {
				// 如果文件不存在，返回404
				c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
				c.Abort()
				return
			}
			// 其他数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			c.Abort()
			return
		}
		// 检查用户是否有权限访问该文件
		if file.UserID != int32(userID) {
			// 如果用户ID不匹配，返回403 Forbidden
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该文件"})
			c.Abort()
			return
		}

		c.Next()
	}
}
