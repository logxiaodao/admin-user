package cron

import (
	conf2 "admin-user/api/cron/conf"
	"github.com/robfig/cron/v3"
)

func Start() {

	c := cron.New()

	// 读取配置添加定时任务
	for _, data := range conf2.Crontab.CronList {
		c.AddFunc(data.TimeFormat, data.Function)
	}

	// 启动执行计划
	c.Start()
}
