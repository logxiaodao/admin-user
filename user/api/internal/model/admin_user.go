package model

import (
	"time"
)

/******sql******
CREATE TABLE `admin_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账户名',
  `password` varchar(255) NOT NULL COMMENT '密码（哈希加盐）',
  `nick_name` varchar(255) NOT NULL DEFAULT '无' COMMENT '昵称',
  `phone` varchar(20) NOT NULL DEFAULT '无' COMMENT '手机号',
  `email` varchar(255) NOT NULL DEFAULT '无' COMMENT '电子邮箱',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '平台id',
  `is_super_admin` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否为超级管理员，0为false，1为true',
  `is_ban` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否封禁  0  false  1 true',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_UNIQUE` (`phone`),
  UNIQUE KEY `email_UNIQUE` (`email`),
  UNIQUE KEY `account_UNIQUE` (`account`),
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `platform_id` (`platform_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-用户表'
******sql******/
// AdminUser 管理端-用户表
type AdminUser struct {
	ID           uint      `gorm:"primaryKey;column:id;type:int unsigned;not null" json:"id"`
	Account      string    `gorm:"unique;column:account;type:varchar(255);not null" json:"account"`                                       // 账户名
	Password     string    `gorm:"column:password;type:varchar(255);not null" json:"password"`                                            // 密码（哈希加盐）
	NickName     string    `gorm:"column:nick_name;type:varchar(255);not null;default:无" json:"nickName"`                                 // 昵称
	Phone        string    `gorm:"unique;column:phone;type:varchar(20);not null;default:无" json:"phone"`                                  // 手机号
	Email        string    `gorm:"unique;column:email;type:varchar(255);not null;default:无" json:"email"`                                 // 电子邮箱
	PlatformID   uint      `gorm:"index:platform_id;column:platform_id;type:int unsigned;not null;default:0" json:"platformId"`           // 平台id
	IsSuperAdmin uint8     `gorm:"column:is_super_admin;type:tinyint unsigned;not null;default:0" json:"isSuperAdmin"`                    // 是否为超级管理员，0为false，1为true
	IsBan        uint8     `gorm:"column:is_ban;type:tinyint unsigned;not null;default:0" json:"isBan"`                                   // 是否封禁  0  false  1 true
	CreatedAt    time.Time `gorm:"index:created_at;column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`                  // 更新时间
	Creater      uint      `gorm:"column:creater;type:int unsigned;not null;default:0" json:"creater"`                                    // 创建者
	Updater      uint      `gorm:"column:updater;type:int unsigned;not null;default:0" json:"updater"`                                    // 更新者
}

// TableName get sql table name.获取数据库表名
func (m *AdminUser) TableName() string {
	return "admin_user"
}

// AdminUserColumns get sql column name.获取数据库列名
var AdminUserColumns = struct {
	ID           string
	Account      string
	Password     string
	NickName     string
	Phone        string
	Email        string
	PlatformID   string
	IsSuperAdmin string
	IsBan        string
	CreatedAt    string
	UpdatedAt    string
	Creater      string
	Updater      string
}{
	ID:           "id",
	Account:      "account",
	Password:     "password",
	NickName:     "nick_name",
	Phone:        "phone",
	Email:        "email",
	PlatformID:   "platform_id",
	IsSuperAdmin: "is_super_admin",
	IsBan:        "is_ban",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	Creater:      "creater",
	Updater:      "updater",
}
