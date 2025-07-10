package service

import (
	"08/dao"
	"08/models"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type updateFileService struct{}

var UpdateFileService = new(updateFileService)

func (s *updateFileService) Update(c *gin.Context, fileID, userID uint, fileContent []byte) (*models.File, error) {
	// 查询文件信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		return nil, errors.New("文件不存在")
	}

	// 权限校验
	if uint(file.UserID) != userID {
		return nil, errors.New("无权限更新")
	}

	// 只能更新文件，不能更新文件夹
	if file.Type == "" {
		return nil, errors.New("不能更新文件夹")
	}

	// 保存旧版本到 tb_version 表
	if file.StoreKey != "" {
		version := &models.Version{
			FileID:   file.ID,
			Size:     file.Size,
			VerNum:   file.VerNum,
			StoreKey: file.StoreKey,
			Ctime:    file.Mtime, // 使用文件的修改时间作为版本的创建时间
		}
		if err := dao.Q.Version.WithContext(c).Create(version); err != nil {
			return nil, fmt.Errorf("保存历史版本失败: %w", err)
		}
	}

	// 生成新的存储路径
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	newStoreKey := fmt.Sprintf("%d_%d_%d_%s", userID, file.ParentID, time.Now().UnixNano(), file.Name)
	savePath := filepath.Join(dataDir, newStoreKey)

	// 保存新文件内容到磁盘
	if err := os.WriteFile(savePath, fileContent, 0644); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 更新文件记录
	newVerNum := file.VerNum + 1
	_, err = dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).Updates(map[string]interface{}{
		"size":      len(fileContent),
		"ver_num":   newVerNum,
		"store_key": newStoreKey,
		"mtime":     int32(time.Now().Unix()),
	})
	if err != nil {
		// 如果数据库更新失败，删除新创建的文件
		os.Remove(savePath)
		return nil, fmt.Errorf("数据库更新失败: %w", err)
	}

	// 返回更新后的文件信息
	updated, _ := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	return updated, nil
}
