package model_def

type StoreRef struct {
	ID       int64  `gorm:"primaryKey"`                    // int64类型，主键
	StoreKey string `gorm:"uniqueIndex;type:varchar(100)"` // 存储的唯一标识符
	RefCount int64  // 引用计数
	Ctime    int64  // 创建时间
	Mtime    int64  // 修改时间
}

func (StoreRef) TableName() string {
	return "tb_store_ref"
}
