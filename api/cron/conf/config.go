package conf

import (
	app2 "admin-user/api/cron/app"
)

type CronConfig struct {
	CronList []CronItem
}

type CronItem struct {
	TimeFormat string
	Function   func()
}

/*
 *  与 linux crontab 设置规则是一样的
 * 	具体规则   https://en.wikipedia.org/wiki/Cron
 */

var Crontab = CronConfig{CronList: []CronItem{
	// 任务模块配置例子
	{
		TimeFormat: "0 0 */15 * *",     // 每15天清理一次消息表数据
		Function:   app2.ClearMessages, // testB
	},
}}
