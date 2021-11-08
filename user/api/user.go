package main

import (
	"admin/user/api/internal/Initialization"
	"admin/user/api/internal/middleware"
	"admin/user/common/errorx"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"

	"admin/user/api/internal/config"
	"admin/user/api/internal/handler"
	"admin/user/api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 关闭输出的统计日志
	logx.DisableStat()

	// 全局中间件
	periodLimit := middleware.NewPeriodLimitMiddleware() // 接口限流
	server.Use(periodLimit.Handle)

	handler.RegisterHandlers(server, ctx)

	// 自定义错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.GetData()
		default:
			return http.StatusInternalServerError, nil
		}
	})

	// 初始化mysql默认数据
	Initialization.InitializationData()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	//cron.Start() // 执行crontab脚本
	server.Start()
}
