package service

import (
	"08/dao"
	"context"
	"errors"
)

type listFileService struct{}

var ListFileService = new(listFileService)

type FileInfo struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Size     int32  `json:"size"`
	Type     string `json:"type"`
	IsFolder bool   `json:"is_folder"`
	Ctime    int32  `json:"ctime"`
	Mtime    int32  `json:"mtime"`
	VerNum   int32  `json:"ver_num"`
}

func (s *listFileService) List(c context.Context, parentID, userID uint) ([]FileInfo, error) {
	// 如果 parentID 不为 0，检查父目录是否存在且有权限
	if parentID != 0 {
		parent, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(parentID))).First()
		if err != nil {
			return nil, errors.New("父目录不存在")
		}
		if uint(parent.UserID) != userID {
			return nil, errors.New("无权限访问该目录")
		}
		// 检查是否为文件夹
		if parent.Type != "folder" {
			return nil, errors.New("指定的不是文件夹")
		}
	}

	// 查询指定父目录下的所有文件和文件夹
	files, err := dao.Q.File.WithContext(c).Where(
		dao.File.ParentID.Eq(int32(parentID)),
		dao.File.UserID.Eq(int32(userID)),
	).Order(dao.File.Type, dao.File.Name).Find() // 按类型和名称排序，文件夹在前

	if err != nil {
		return nil, errors.New("查询文件列表失败")
	}

	// 转换为响应格式
	var result []FileInfo
	for _, file := range files {
		fileInfo := FileInfo{
			ID:       file.ID,
			Name:     file.Name,
			Size:     file.Size,
			Type:     file.Type,
			IsFolder: file.Type == "folder", // Type 为空表示文件夹
			Ctime:    file.Ctime,
			Mtime:    file.Mtime,
			VerNum:   file.VerNum,
		}
		result = append(result, fileInfo)
	}

	return result, nil
}
