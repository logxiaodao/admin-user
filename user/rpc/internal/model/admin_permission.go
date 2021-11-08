package model

import (
	"time"
)

/******sql******
CREATE TABLE `admin_permission` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `permission_name` varchar(255) NOT NULL COMMENT '权限名称',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `permission_name` (`platform_id`,`permission_name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-权限表'
******sql******/
// AdminPermission 管理端-权限表
type AdminPermission struct {
	ID             uint      `gorm:"primaryKey;column:id;type:int unsigned;not null" json:"id"`                                             // 权限id
	PermissionName string    `gorm:"uniqueIndex:permission_name;column:permission_name;type:varchar(255);not null" json:"permissionName"`   // 权限名称
	PlatformID     uint      `gorm:"uniqueIndex:permission_name;column:platform_id;type:int unsigned;not null;default:0" json:"platformId"` // 平台id
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`                  // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`                  // 更新时间
	Creater        uint      `gorm:"column:creater;type:int unsigned;not null;default:0" json:"creater"`                                    // 创建者
	Updater        uint      `gorm:"column:updater;type:int unsigned;not null;default:0" json:"updater"`                                    // 更新者
}

// TableName get sql table name.获取数据库表名
func (m *AdminPermission) TableName() string {
	return "admin_permission"
}

// AdminPermissionColumns get sql column name.获取数据库列名
var AdminPermissionColumns = struct {
	ID             string
	PermissionName string
	PlatformID     string
	CreatedAt      string
	UpdatedAt      string
	Creater        string
	Updater        string
}{
	ID:             "id",
	PermissionName: "permission_name",
	PlatformID:     "platform_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	Creater:        "creater",
	Updater:        "updater",
}
