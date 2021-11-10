package Initialization

import (
	safe2 "admin-user/rpc/common/safe"
	config2 "admin-user/rpc/internal/config"
	model2 "admin-user/rpc/internal/model"
	"fmt"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	wg sync.WaitGroup
	ds *DataSources
)

// InitializationData 数据初始化
func InitializationData() {
	wg, ds = sync.WaitGroup{}, GetDataSource()
	wg.Add(2)
	go putRouterIntoDB(wg, ds, config2.InitRouteData)
	go putAdminIntoDB(wg, ds, config2.InitAdminData)
	go putPlatformIntoDB(wg, ds, config2.InitPlatformData)
}

// putRouterIntoDB api 初始化
func putRouterIntoDB(wg sync.WaitGroup, ds *DataSources, route config2.InitRouteConf) {
	defer wg.Done()

	for k, _ := range route.Data { // 初始化平台
		route.Data[k].PlatformID = uint(route.PlatformId)
	}

	// sql 逻辑  触发唯一索引限制则不插入
	err := ds.DB.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(route.Data, 500).Error
	if err != nil {
		fmt.Printf("[debug] 初始化路由失败: %s \n", err.Error())
	} else {
		fmt.Printf("[debug] 路由初始化入库成功 \n")
	}
}

// putAdminIntoDB 初始化系统管理员
func putAdminIntoDB(wg sync.WaitGroup, ds *DataSources, admin config2.InitAdminConf) {
	defer wg.Done()

	tx := ds.DB.Begin()
	var (
		AdminUser model2.AdminUser
	)

	for _, v := range admin.Data {
		pwd, err := safe2.GenHashPassword(v.Password)
		if err != nil {
			tx.Rollback()
			fmt.Printf("[debug] 初始化管理员用户%s设置的密码不规范: %s \n", v.Account, err.Error())
			return
		}

		AdminUser = model2.AdminUser{
			Account:      v.Account,
			Password:     pwd,
			NickName:     v.NickName,
			Phone:        v.Phone,
			Email:        v.Email,
			IsSuperAdmin: 1,
			PlatformID:   config2.AdminPlatformId,
		}

		// sql 逻辑  插入用户表
		err = tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&AdminUser).Error
		if err != nil {
			tx.Rollback()
			fmt.Printf("[debug] 管理员用户数据初始化入库失败: %s \n", err.Error())
			return
		}
	}

	tx.Commit()
	fmt.Printf("[debug] 管理员初始化入库成功 \n")
}

// putPlatformIntoDB 初始化系统平台
func putPlatformIntoDB(wg sync.WaitGroup, ds *DataSources, platform config2.InitPlatformConf) {
	defer wg.Done()

	// sql 逻辑  触发唯一索引限制则不插入
	err := ds.DB.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(platform.Data, 500).Error
	if err != nil {
		fmt.Printf("[debug] 平台初始化入库失败: %s \n", err.Error())
	} else {
		fmt.Printf("[debug] 平台初始化入库成功 \n")
	}
}
