package model

import (
	"time"
)

/******sql******
CREATE TABLE `admin_api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '接口名',
  `http_method` varchar(45) NOT NULL COMMENT 'http请求方式',
  `http_path` varchar(255) NOT NULL COMMENT 'http请求路径',
  `platform_id` int unsigned NOT NULL COMMENT '平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_open` tinyint NOT NULL DEFAULT '0' COMMENT '0 接口权限接受配置 1 接口对所有人开放 ',
  `is_super` tinyint NOT NULL DEFAULT '0' COMMENT '0 接口权限接受配置  1 接口只对超级管理员开放',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`platform_id`,`name`),
  UNIQUE KEY `api_identify` (`platform_id`,`http_method`,`http_path`)
) ENGINE=InnoDB AUTO_INCREMENT=7518 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-api表'
******sql******/
// AdminAPI 管理端-api表
type AdminAPI struct {
	ID         uint      `gorm:"primaryKey;column:id;type:int unsigned;not null" json:"id"`
	Name       string    `gorm:"uniqueIndex:name;column:name;type:varchar(255);not null" json:"name"`                                       // 接口名
	HTTPMethod string    `gorm:"uniqueIndex:api_identify;column:http_method;type:varchar(45);not null" json:"httpMethod"`                   // http请求方式
	HTTPPath   string    `gorm:"uniqueIndex:api_identify;column:http_path;type:varchar(255);not null" json:"httpPath"`                      // http请求路径
	PlatformID uint      `gorm:"uniqueIndex:name;uniqueIndex:api_identify;column:platform_id;type:int unsigned;not null" json:"platformId"` // 平台id
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`                      // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`                      // 更新时间
	IsOpen     int8      `gorm:"column:is_open;type:tinyint;not null;default:0" json:"isOpen"`                                              // 0 接口权限接受配置 1 接口对所有人开放
	IsSuper    int8      `gorm:"column:is_super;type:tinyint;not null;default:0" json:"isSuper"`                                            // 0 接口权限接受配置  1 接口只对超级管理员开放
}

// TableName get sql table name.获取数据库表名
func (m *AdminAPI) TableName() string {
	return "admin_api"
}

// AdminAPIColumns get sql column name.获取数据库列名
var AdminAPIColumns = struct {
	ID         string
	Name       string
	HTTPMethod string
	HTTPPath   string
	PlatformID string
	CreatedAt  string
	UpdatedAt  string
	IsOpen     string
	IsSuper    string
}{
	ID:         "id",
	Name:       "name",
	HTTPMethod: "http_method",
	HTTPPath:   "http_path",
	PlatformID: "platform_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	IsOpen:     "is_open",
	IsSuper:    "is_super",
}
