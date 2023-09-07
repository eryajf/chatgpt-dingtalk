package db

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
)

// 全局数据库对象
var DB *gorm.DB

// 初始化数据库
func InitDB() {
	DB = ConnSqlite()

	dbAutoMigrate()
}

// 自动迁移表结构
func dbAutoMigrate() {
	_ = DB.AutoMigrate(
		Chat{},
	)
}

func ConnSqlite() *gorm.DB {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return nil
	}
	db, err := gorm.Open(sqlite.Open("data/dingtalkbot.sqlite"), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Fatal("failed to connect sqlite3: %v", err)
	}
	dbObj, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sqlite3 obj: %v", err)
	}
	// 参见： https://github.com/glebarez/sqlite/issues/52
	dbObj.SetMaxOpenConns(1)
	return db
}
