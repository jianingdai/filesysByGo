package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"08/dao"
	"08/models"

	"github.com/gin-gonic/gin"
)

type uploadFileService struct{}

var UploadFileService = new(uploadFileService)

// 保存文件到 data 目录，并写入数据库
func (s *uploadFileService) Upload(c *gin.Context, file *multipart.FileHeader, parentID, userID uint) (*models.File, error) {
	// 生成存储路径
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}
	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%d_%d_%s", userID, parentID, time.Now().UnixNano(), file.Filename)
	savePath := filepath.Join(dataDir, filename)

	// 保存文件到磁盘
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 写入数据库
	newFile := &models.File{
		UserID:   int32(userID),
		ParentID: int32(parentID),
		Name:     file.Filename,
		Size:     int32(file.Size),
		Type:     filepath.Ext(file.Filename),
		StoreKey: filename,
		Ctime:    int32(time.Now().Unix()),
		Mtime:    int32(time.Now().Unix()),
		VerNum:   1,
	}
	if err := dao.Q.File.WithContext(c.Request.Context()).Create(newFile); err != nil {
		return nil, fmt.Errorf("数据库写入失败: %w", err)
	}

	return newFile, nil
}
