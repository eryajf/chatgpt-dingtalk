package db

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type ChatType uint

const Q ChatType = 1
const A ChatType = 2

type Chat struct {
	gorm.Model
	Username      string   `gorm:"type:varchar(50);not null;comment:'用户名'" json:"username"`        // 用户名
	Source        string   `gorm:"type:varchar(50);comment:'用户来源：群聊名字，私聊'" json:"source"`          // 对话来源
	ChatType      ChatType `gorm:"type:tinyint(1);default:1;comment:'类型:1问, 2答'" json:"chat_type"` // 状态
	ParentContent uint     `gorm:"default:0;comment:'父消息编号(编号为0时表示为首条)'" json:"parent_content"`
	Content       string   `gorm:"type:varchar(128);comment:'内容'" json:"content"` // 问题或回答的内容
}

type ChatListReq struct {
	Username string `json:"username" form:"username"`
	Source   string `json:"source" form:"source"`
}

// Add 添加资源
func (c Chat) Add() (uint, error) {
	err := DB.Create(&c).Error
	return c.ID, err
}

// Find 获取单个资源
func (c Chat) Find(filter map[string]interface{}, data *Chat) error {
	return DB.Where(filter).First(&data).Error
}

// List 获取数据列表
func (c Chat) List(req ChatListReq) ([]*Chat, error) {
	var list []*Chat
	db := DB.Model(&Chat{}).Order("created_at ASC")

	userName := strings.TrimSpace(req.Username)
	if userName != "" {
		db = db.Where("username = ?", userName)
	}
	source := strings.TrimSpace(req.Source)
	if source != "" {
		db = db.Where("source = ?", source)
	}

	err := db.Find(&list).Error
	return list, err
}

// Exist 判断资源是否存在
func (c Chat) Exist(filter map[string]interface{}) bool {
	var dataObj Chat
	err := DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
