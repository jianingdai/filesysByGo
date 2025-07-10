package model_def

type User struct {
	ID       int64  `gorm:"primaryKey"`                    // int64类型，主键
	Name     string `gorm:"uniqueIndex;type:varchar(100)"` // 用户名，唯一
	Password string `gorm:"type:varchar(100)"`             // 密码，存储哈希值
	Ctime    int64  // 创建时间
	Mtime    int64  // 修改时间
}

func (User) TableName() string {
	return "tb_user"
}
