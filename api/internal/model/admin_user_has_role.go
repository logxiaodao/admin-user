package model

/******sql******
CREATE TABLE `admin_user_has_role` (
  `user_id` int unsigned NOT NULL COMMENT '用户id',
  `role_id` int unsigned NOT NULL COMMENT '角色id',
  PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-用户与角色中间表'
******sql******/
// AdminUserHasRole 管理端-用户与角色中间表
type AdminUserHasRole struct {
	UserID uint `gorm:"primaryKey;column:user_id;type:int unsigned;not null" json:"userId"` // 用户id
	RoleID uint `gorm:"primaryKey;column:role_id;type:int unsigned;not null" json:"roleId"` // 角色id
}

// TableName get sql table name.获取数据库表名
func (m *AdminUserHasRole) TableName() string {
	return "admin_user_has_role"
}

// AdminUserHasRoleColumns get sql column name.获取数据库列名
var AdminUserHasRoleColumns = struct {
	UserID string
	RoleID string
}{
	UserID: "user_id",
	RoleID: "role_id",
}
