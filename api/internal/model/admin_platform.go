package model

/******sql******
CREATE TABLE `admin_platform` (
  `id` int unsigned NOT NULL,
  `platform_en` varchar(128) NOT NULL COMMENT '平台-英文',
  `platform_zh` varchar(255) NOT NULL COMMENT '平台-中文',
  PRIMARY KEY (`id`),
  UNIQUE KEY `platform` (`platform_en`,`platform_zh`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
******sql******/
// AdminPlatform [...]
type AdminPlatform struct {
	ID         uint   `gorm:"primaryKey;column:id;type:int unsigned;not null" json:"id"`
	PlatformEn string `gorm:"uniqueIndex:platform;column:platform_en;type:varchar(128);not null" json:"platformEn"` // 平台-英文
	PlatformZh string `gorm:"uniqueIndex:platform;column:platform_zh;type:varchar(255);not null" json:"platformZh"` // 平台-中文
}

// TableName get sql table name.获取数据库表名
func (m *AdminPlatform) TableName() string {
	return "admin_platform"
}

// AdminPlatformColumns get sql column name.获取数据库列名
var AdminPlatformColumns = struct {
	ID         string
	PlatformEn string
	PlatformZh string
}{
	ID:         "id",
	PlatformEn: "platform_en",
	PlatformZh: "platform_zh",
}
