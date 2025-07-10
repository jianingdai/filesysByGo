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

type moveFileService struct{}

var MoveFileService = new(moveFileService)

func (s *moveFileService) Move(c context.Context, fileID, targetParentID, userID uint) (*models.File, error) {
	// 查询源文件/文件夹信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		return nil, errors.New("源文件不存在")
	}

	// 权限校验
	if uint(file.UserID) != userID {
		return nil, errors.New("无权限移动")
	}

	// 检查目标父目录是否存在且有权限
	if targetParentID != 0 {
		targetParent, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(targetParentID))).First()
		if err != nil {
			return nil, errors.New("目标目录不存在")
		}
		if uint(targetParent.UserID) != userID {
			return nil, errors.New("无权限在目标目录操作")
		}
	}

	// 不能移动到自己或自己的子孙目录下
	if fileID == targetParentID {
		return nil, errors.New("不能移动到自身目录下")
	}
	if s.isDescendant(c, fileID, targetParentID) {
		return nil, errors.New("不能移动到自己的子目录下")
	}

	// 处理重名
	newName := s.getUniqueNameInParent(c, file.Name, int32(targetParentID), int32(userID))

	// 更新 ParentID 和 Name
	_, err = dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).Updates(map[string]interface{}{
		"parent_id": targetParentID,
		"name":      newName,
		"mtime":     int32(time.Now().Unix()),
	})
	if err != nil {
		return nil, fmt.Errorf("数据库更新失败: %w", err)
	}

	// 返回最新信息
	updated, _ := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	return updated, nil
}

// 判断目标目录是否为 fileID 的子孙目录（防止循环）
func (s *moveFileService) isDescendant(c context.Context, fileID, targetParentID uint) bool {
	currID := targetParentID
	for currID != 0 {
		parent, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(currID))).First()
		if err != nil {
			break
		}
		if uint(parent.ParentID) == fileID {
			return true
		}
		currID = uint(parent.ParentID)
	}
	return false
}

// 获取在指定父目录下的唯一文件名（处理重名）
func (s *moveFileService) getUniqueNameInParent(c context.Context, originalName string, parentID, userID int32) string {
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
