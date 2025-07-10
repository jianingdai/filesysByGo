package service

import (
	"08/dao"
	"context"
	"errors"
	"os"
)

type deleteFileService struct{}

var DeleteFileService = new(deleteFileService)

func (s *deleteFileService) Delete(c context.Context, fileID, userID uint) error {
	// 查询文件信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		return errors.New("文件不存在")
	}
	// 权限校验（可选，通常中间件已做）
	if uint(file.UserID) != userID {
		return errors.New("无权限删除")
	}
	// 如果是文件夹，可以递归删除其下所有内容（此处略，可后续完善）

	// 删除物理文件（如果是文件）
	if file.Type != "" && file.StoreKey != "" {
		_ = os.Remove("data/" + file.StoreKey)
	}
	// 删除数据库记录
	_, err = dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).Delete()
	return err
}
