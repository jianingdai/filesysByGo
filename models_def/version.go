package model_def

type Version struct {
	ID       int64 `gorm:"primaryKey"` // int64类型，主键
	FileID   int64 `gorm:"index"`
	Size     int64
	VerNum   int64  // 版本号，文件夹为0, 文件从1开始递增
	StoreKey string `gorm:"type:varchar(100)"` // 存储的唯一标识符，用雪花算法生成
	Ctime    int64
}

func (Version) TableName() string {
	return "tb_version"
}
