package model_def

type File struct {
	ID       int64  `gorm:"primaryKey"` // int64类型，主键
	UserID   int64  `gorm:"uniqueIndex:uniq_user_parent_name"`
	ParentID int64  `gorm:"uniqueIndex:uniq_user_parent_name"` // 父目录ID，根目录为0
	Name     string `gorm:"uniqueIndex:uniq_user_parent_name;type:varchar(1024)"`
	Size     int64  // 文件大小，文件夹为0
	Type     string `gorm:"type:varchar(10)"` // 文件类型，只有两种："file"和"folder"
	VerNum   int64  // 版本号，文件夹为0, 文件从1开始递增
	StoreKey string `gorm:"type:varchar(100)"` // 存储的唯一标识符，用雪花算法生成
	Ctime    int64  // 创建时间
	Mtime    int64  // 修改时间
}

func (File) TableName() string {
	return "tb_file"
}
