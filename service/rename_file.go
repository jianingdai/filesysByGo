package service

import (
	"08/dao"
	"08/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type renameFileService struct{}

var RenameFileService = new(renameFileService)

func (s *renameFileService) Rename(c context.Context, fileID, userID uint, newName string) (*models.File, error) {
	// 查询文件信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		return nil, errors.New("文件不存在")
	}
	if uint(file.UserID) != userID {
		return nil, errors.New("无权限重命名")
	}
	// 处理重名
	uniqueName := s.getUniqueNameInParent(c, newName, file.ParentID, file.UserID)
	// 更新名称
	_, err = dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).Updates(map[string]interface{}{
		"name":  uniqueName,
		"mtime": int32(time.Now().Unix()),
	})
	if err != nil {
		return nil, fmt.Errorf("数据库更新失败: %w", err)
	}
	// 返回最新信息
	updated, _ := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	return updated, nil
}

// 获取在指定父目录下的唯一文件名（处理重名）
func (s *renameFileService) getUniqueNameInParent(c context.Context, originalName string, parentID, userID int32) string {
	baseName := originalName
	extension := ""
	if dotIndex := strings.LastIndex(originalName, "."); dotIndex != -1 {
		baseName = originalName[:dotIndex]
		extension = originalName[dotIndex:]
	}
	name := originalName
	counter := 1
	for {
		_, err := dao.Q.File.WithContext(c).Where(
			dao.File.ParentID.Eq(parentID),
			dao.File.UserID.Eq(userID),
			dao.File.Name.Eq(name),
		).First()
		if err != nil {
			return name
		}
		name = fmt.Sprintf("%s(%d)%s", baseName, counter, extension)
		counter++
	}
}
