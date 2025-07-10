package endpoint

import (
	"08/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取文件的版本历史列表
func GetFileVersions(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 file_id"})
		return
	}

	// 查询文件信息
	file, err := dao.Q.File.WithContext(c).Where(dao.File.ID.Eq(int32(fileID))).First()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	userID := c.GetInt("user_id")
	if uint(file.UserID) != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	if file.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹没有版本历史"})
		return
	}

	// 查询历史版本
	versions, err := dao.Q.Version.WithContext(c).Where(
		dao.Version.FileID.Eq(int32(fileID)),
	).Order(dao.Version.VerNum.Desc()).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询版本历史失败"})
		return
	}

	// 构造响应数据，包含当前版本
	type VersionInfo struct {
		VerNum    int32 `json:"ver_num"`
		Size      int32 `json:"size"`
		Ctime     int32 `json:"ctime"`
		IsCurrent bool  `json:"is_current"`
	}

	var result []VersionInfo

	// 添加当前版本
	result = append(result, VersionInfo{
		VerNum:    file.VerNum,
		Size:      file.Size,
		Ctime:     file.Mtime,
		IsCurrent: true,
	})

	// 添加历史版本
	for _, version := range versions {
		result = append(result, VersionInfo{
			VerNum:    version.VerNum,
			Size:      version.Size,
			Ctime:     version.Ctime,
			IsCurrent: false,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":  fileID,
		"versions": result,
	})
}
