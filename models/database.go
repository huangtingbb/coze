package models

import (
	"coze-agent-platform/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	cfg := config.Cfg.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		log.Println("警告: 数据库不可用，某些功能可能无法正常工作")
		return
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&User{},
		&Agent{},
		&Conversation{},
		&Message{},
	)

	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return
	}

	log.Println("数据库初始化完成")
}
