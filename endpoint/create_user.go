package endpoint

import (
	"08/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(c *gin.Context) {
	userName := c.GetString("user_name")
	if userName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "仅管理员可创建用户"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := service.CreateUserService.CreateUser(c.Request.Context(), req.Username, req.Password); err != nil {
		// 判断是否唯一索引冲突
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}

		// 其他错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户创建成功"})
}
