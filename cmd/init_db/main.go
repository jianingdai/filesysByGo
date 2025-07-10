package main

import (
	model_def "08/models_def"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 清空数据库所有数据（保留表结构）
func clearAllData(db *gorm.DB) error {
	// 按依赖关系逆序删除，避免外键约束问题
	if err := db.Where("1 = 1").Delete(&model_def.Session{}).Error; err != nil {
		return fmt.Errorf("清空 sessions 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.StoreRef{}).Error; err != nil {
		return fmt.Errorf("清空 store_refs 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.Version{}).Error; err != nil {
		return fmt.Errorf("清空 versions 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.File{}).Error; err != nil {
		return fmt.Errorf("清空 files 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.User{}).Error; err != nil {
		return fmt.Errorf("清空 users 表失败: %v", err)
	}
	return nil
}

func main() {
	// 1. 链接数据库
	db, err := gorm.Open(sqlite.Open("filesys.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到数据库: ", err)
	}

	// 2. 自动建表(生产环境中通常不建议使用AutoMigrate, 但在这里为了简化示例)
	var migrateErrs []string
	if err := db.AutoMigrate(&model_def.User{}); err != nil {
		migrateErrs = append(migrateErrs, "User")
	}
	if err := db.AutoMigrate(&model_def.File{}); err != nil {
		migrateErrs = append(migrateErrs, "File")
	}
	if err := db.AutoMigrate(&model_def.Version{}); err != nil {
		migrateErrs = append(migrateErrs, "Version")
	}
	if err := db.AutoMigrate(&model_def.StoreRef{}); err != nil {
		migrateErrs = append(migrateErrs, "StoreRef")
	}
	if err := db.AutoMigrate(&model_def.Session{}); err != nil {
		migrateErrs = append(migrateErrs, "Session")
	}
	if len(migrateErrs) > 0 {
		log.Fatalf("以下表建表失败: %v", migrateErrs)
	}

	// 3. 清空数据
	log.Println("开始清空所有表数据...")
	if err := clearAllData(db); err != nil {
		log.Printf("清空数据失败: %v", err)
	}

	// 4. 创建默认管理员,数据库里面存储sha256加密后的密码
	adminUser := &model_def.User{Name: "admin", Password: "admin123"}
	hash := sha256.Sum256([]byte(adminUser.Password))
	hashStr := hex.EncodeToString(hash[:])
	adminUser.Password = hashStr
	adminUser.Ctime = time.Now().Unix()
	adminUser.Mtime = adminUser.Ctime
	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("创建默认用户失败: %v", err)
	}

	log.Println("所有表数据已清空")
}
