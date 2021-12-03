package main

import (
	Initialization2 "admin-user/api/internal/Initialization"
	config2 "admin-user/api/internal/config"
	handler2 "admin-user/api/internal/handler"
	middleware2 "admin-user/api/internal/middleware"
	svc2 "admin-user/api/internal/svc"
	errorx2 "admin-user/rpc/common/errorx"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
)

var configFile = flag.String("f", "etc/dev.yaml", "the config file")

func main() {
	flag.Parse()

	var c config2.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc2.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 关闭输出的统计日志
	logx.DisableStat()

	// 全局中间件
	periodLimit := middleware2.NewPeriodLimitMiddleware() // 接口限流
	server.Use(periodLimit.Handle)

	handler2.RegisterHandlers(server, ctx)

	// 自定义错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx2.CodeError:
			return http.StatusOK, e.GetData()
		default:
			return http.StatusInternalServerError, nil
		}
	})

	// 初始化mysql默认数据
	Initialization2.InitializationData()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	//cron.Start() // 执行crontab脚本
	server.Start()
}
