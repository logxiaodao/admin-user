package model

import (
	"time"
)

/******sql******
CREATE TABLE `admin_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `role_name` varchar(255) NOT NULL COMMENT '角色名称',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_name` (`platform_id`,`role_name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-角色表'
******sql******/
// AdminRole 管理端-角色表
type AdminRole struct {
	ID         uint      `gorm:"primaryKey;column:id;type:int unsigned;not null" json:"id"`                                       // 角色id
	RoleName   string    `gorm:"uniqueIndex:role_name;column:role_name;type:varchar(255);not null" json:"roleName"`               // 角色名称
	PlatformID uint      `gorm:"uniqueIndex:role_name;column:platform_id;type:int unsigned;not null;default:0" json:"platformId"` // 关联的平台id
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`            // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`            // 更新时间
	Creater    uint      `gorm:"column:creater;type:int unsigned;not null;default:0" json:"creater"`                              // 创建者
	Updater    uint      `gorm:"column:updater;type:int unsigned;not null;default:0" json:"updater"`                              // 更新者
}

// TableName get sql table name.获取数据库表名
func (m *AdminRole) TableName() string {
	return "admin_role"
}

// AdminRoleColumns get sql column name.获取数据库列名
var AdminRoleColumns = struct {
	ID         string
	RoleName   string
	PlatformID string
	CreatedAt  string
	UpdatedAt  string
	Creater    string
	Updater    string
}{
	ID:         "id",
	RoleName:   "role_name",
	PlatformID: "platform_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	Creater:    "creater",
	Updater:    "updater",
}
