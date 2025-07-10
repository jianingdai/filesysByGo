package service

import (
	"08/dao"
	"08/models"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type copyFileService struct{}

var CopyFileService = new(copyFileService)

func (s *copyFileService) Copy(c context.Context, fileID, targetParentID, userID uint) (*models.File, error) {
	// 查询源文件/文件夹信息
	sourceFile, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		return nil, errors.New("源文件不存在")
	}

	// 权限校验
	if uint(sourceFile.UserID) != userID {
		return nil, errors.New("无权限复制")
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

	// 判断是文件还是文件夹（文件夹的 Type 为空）
	if sourceFile.Type == "" {
		// 复制文件夹
		return s.copyFolder(c, sourceFile, targetParentID, userID)
	} else {
		// 复制文件
		return s.copyFile(c, sourceFile, targetParentID, userID)
	}
}

// 复制单个文件
func (s *copyFileService) copyFile(c context.Context, sourceFile *models.File, targetParentID uint, userID uint) (*models.File, error) {
	// 处理重名
	newName := s.getUniqueNameInParent(c, sourceFile.Name, int32(targetParentID), int32(userID))

	// 复制物理文件
	newStoreKey := ""
	if sourceFile.StoreKey != "" {
		newStoreKey = fmt.Sprintf("%d_%d_%d_%s", userID, targetParentID, time.Now().UnixNano(), sourceFile.Name)
		if err := s.copyPhysicalFile(sourceFile.StoreKey, newStoreKey); err != nil {
			return nil, fmt.Errorf("复制物理文件失败: %w", err)
		}
	}

	// 创建新文件记录
	newFile := &models.File{
		UserID:   int32(userID),
		ParentID: int32(targetParentID),
		Name:     newName,
		Size:     sourceFile.Size,
		Type:     sourceFile.Type,
		StoreKey: newStoreKey,
		Ctime:    int32(time.Now().Unix()),
		Mtime:    int32(time.Now().Unix()),
		VerNum:   1,
	}

	if err := dao.Q.File.WithContext(c).Create(newFile); err != nil {
		return nil, fmt.Errorf("创建文件记录失败: %w", err)
	}

	return newFile, nil
}

// 复制文件夹（递归）
func (s *copyFileService) copyFolder(c context.Context, sourceFolder *models.File, targetParentID uint, userID uint) (*models.File, error) {
	// 处理重名
	newName := s.getUniqueNameInParent(c, sourceFolder.Name, int32(targetParentID), int32(userID))

	// 创建新文件夹记录
	newFolder := &models.File{
		UserID:   int32(userID),
		ParentID: int32(targetParentID),
		Name:     newName,
		Size:     0,
		Type:     "", // 文件夹的 Type 为空
		StoreKey: "",
		Ctime:    int32(time.Now().Unix()),
		Mtime:    int32(time.Now().Unix()),
		VerNum:   1,
	}

	if err := dao.Q.File.WithContext(c).Create(newFolder); err != nil {
		return nil, fmt.Errorf("创建文件夹记录失败: %w", err)
	}

	// 查询源文件夹下的所有子文件/子文件夹
	children, err := dao.Q.File.WithContext(c).Where(dao.File.ParentID.Eq(sourceFolder.ID)).Find()
	if err != nil {
		return newFolder, nil // 如果查询子文件失败，至少返回已创建的文件夹
	}

	// 递归复制所有子文件/子文件夹
	for _, child := range children {
		_, _ = s.Copy(c, uint(child.ID), uint(newFolder.ID), userID)
	}

	return newFolder, nil
}

// 复制物理文件
func (s *copyFileService) copyPhysicalFile(sourceStoreKey, targetStoreKey string) error {
	sourcePath := filepath.Join("data", sourceStoreKey)
	targetPath := filepath.Join("data", targetStoreKey)

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}

// 获取在指定父目录下的唯一文件名（处理重名）
func (s *copyFileService) getUniqueNameInParent(c context.Context, originalName string, parentID, userID int32) string {
	baseName := originalName
	extension := ""

	// 分离文件名和扩展名
	if dotIndex := strings.LastIndex(originalName, "."); dotIndex != -1 {
		baseName = originalName[:dotIndex]
		extension = originalName[dotIndex:]
	}

	name := originalName
	counter := 1

	for {
		// 检查当前名称是否已存在
		_, err := dao.Q.File.WithContext(c).Where(
			dao.File.ParentID.Eq(parentID),
			dao.File.UserID.Eq(userID),
			dao.File.Name.Eq(name),
		).First()

		if err != nil {
			// 如果查不到，说明名称可用
			return name
		}

		// 如果存在重名，生成新名称
		name = fmt.Sprintf("%s(%d)%s", baseName, counter, extension)
		counter++
	}
}
