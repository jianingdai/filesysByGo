package service

import (
	"08/dao"
	"08/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type createUserService struct{}

var CreateUserService = &createUserService{}

func (s *createUserService) CreateUser(ctx context.Context, username, password string) error {
	// 计算 SHA-256 哈希
	hash := sha256.Sum256([]byte(password))

	// 转为十六进制字符串
	hashStr := hex.EncodeToString(hash[:])

	user := &models.User{
		Name:     username,
		Password: hashStr,
		Ctime:    int32(time.Now().Unix()), // 创建时间
		Mtime:    int32(time.Now().Unix()), // 修改时间
	}
	return dao.Q.User.WithContext(ctx).Create(user)
}
