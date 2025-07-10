package main

import (
	"08/dao"
	"08/router"
	"08/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("filesys.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dao.SetDefault(db)          //让gen生成的代码使用这个db
	utils.StartSessionCleaner() // 启动定时清理过期session的协程
	r := router.InitRouter()
	r.Run(":8080")
}
