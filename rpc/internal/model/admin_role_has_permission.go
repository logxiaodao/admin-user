package model

/******sql******
CREATE TABLE `admin_role_has_permission` (
  `role_id` int unsigned NOT NULL COMMENT '角色id',
  `permission_id` int unsigned NOT NULL COMMENT '权限id',
  PRIMARY KEY (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-角色与权限中间表'
******sql******/
// AdminRoleHasPermission 管理端-角色与权限中间表
type AdminRoleHasPermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id;type:int unsigned;not null" json:"roleId"`             // 角色id
	PermissionID uint `gorm:"primaryKey;column:permission_id;type:int unsigned;not null" json:"permissionId"` // 权限id
}

// TableName get sql table name.获取数据库表名
func (m *AdminRoleHasPermission) TableName() string {
	return "admin_role_has_permission"
}

// AdminRoleHasPermissionColumns get sql column name.获取数据库列名
var AdminRoleHasPermissionColumns = struct {
	RoleID       string
	PermissionID string
}{
	RoleID:       "role_id",
	PermissionID: "permission_id",
}
