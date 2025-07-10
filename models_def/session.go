package model_def

type Session struct {
	ID        int64 `gorm:"primaryKey"` // int64类型，主键
	UserID    int64
	SessionId string `gorm:"index;type:varchar(100)"` // 会话 ID，用uuid生成
	Ctime     int64
	Etime     int64 // 过期时间
}

func (Session) TableName() string {
	return "tb_session"
}
