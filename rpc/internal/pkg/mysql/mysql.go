package mysql

import (
	config2 "admin-user/rpc/internal/config"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New 返回mysql 实例
func New(mysqlConf config2.MySQLConfig) (*gorm.DB, error) {
	db, err := newMysql(mysqlConf)
	return db, err
}

func newMysql(mysqlConf config2.MySQLConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		mysqlConf.UserName,
		mysqlConf.Password,
		mysqlConf.Addr,
		mysqlConf.Name,
		true,
		"Local")
	gormConfig := &gorm.Config{}
	// 根据配置文件 是否打印sql日志
	if mysqlConf.ShowLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(dns), gormConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open mysql")
	}
	sql, err := db.DB()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to config mysql")
	}
	sql.SetMaxIdleConns(mysqlConf.MaxIdleConn)
	sql.SetMaxOpenConns(mysqlConf.MaxOpenConn)
	sql.SetConnMaxLifetime(time.Minute * time.Duration(mysqlConf.ConnMaxLifeTime))

	return db, err

}
