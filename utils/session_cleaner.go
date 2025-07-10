package utils

import (
	"08/dao"
	"context"
	"log"
	"time"
)

func StartSessionCleaner() {
	go func() {
		// 每10分钟清理一次过期的session
		log.Println("定期清理过期session已启动")
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			//阻塞等待定时器的下一个时间点。
			<-ticker.C
			_, err := dao.Q.Session.WithContext(context.Background()).
				Where(dao.Session.Etime.Lt(int32(time.Now().Unix()))).Delete()
			if err != nil {
				log.Println("定期清理过期session失败:", err)
			}
		}
	}()
}
