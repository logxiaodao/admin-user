package model

/******sql******
CREATE TABLE `admin_permission_has_api` (
  `permission_id` int unsigned NOT NULL COMMENT '权限id',
  `api_id` int unsigned NOT NULL COMMENT 'api_id',
  PRIMARY KEY (`permission_id`,`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理端-权限与api中间表'
******sql******/
// AdminPermissionHasAPI 管理端-权限与api中间表
type AdminPermissionHasAPI struct {
	PermissionID uint `gorm:"primaryKey;column:permission_id;type:int unsigned;not null" json:"permissionId"` // 权限id
	APIID        uint `gorm:"primaryKey;column:api_id;type:int unsigned;not null" json:"apiId"`               // api_id
}

// TableName get sql table name.获取数据库表名
func (m *AdminPermissionHasAPI) TableName() string {
	return "admin_permission_has_api"
}

// AdminPermissionHasAPIColumns get sql column name.获取数据库列名
var AdminPermissionHasAPIColumns = struct {
	PermissionID string
	APIID        string
}{
	PermissionID: "permission_id",
	APIID:        "api_id",
}
