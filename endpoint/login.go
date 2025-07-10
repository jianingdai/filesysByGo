package endpoint

import (
	"08/dao"
	"08/models"
	"08/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Loginrequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req Loginrequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	uid, err := service.LoginService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	// 设置会话cookie
	// 生成uuid作为会话ID
	sessionID := uuid.New().String()

	dao.Q.Session.WithContext(c.Request.Context()).Create(&models.Session{
		UserID:    uid,
		SessionID: sessionID,
		Ctime:     int32(time.Now().Unix()),
		Etime:     int32(time.Now().Add(time.Hour).Unix()), // 设置过期时间为1小时后
	})
	c.SetCookie("sid", sessionID, 3600, "/", "", false, true) // 设置cookie，1小时后过期
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "sid": sessionID})
}
