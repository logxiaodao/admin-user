## 1.背景说明
> 本服使用go-zero框架,go-zero 是一个集成了各种工程实践的 web 和 rpc 框架。通过弹性设计保障了大并发服务端的稳定性，经受了充分的实战检验。
### 1.1 安装依赖
```
$ 预先需安装服务
$ mysql(关系型数据库)、redis(key-value数据库[缓存])、etcd(key-value数据库[服务发现])
$
$ 建议使用homebrew或者docker管理服务
$ brew install --build-from-source etcd
$ brew install --build-from-source mysql    这里mysql的数据默认会存放在reds里面
$ brew install --build-from-source redis
$
$ 选择性安装prometheus(监控告警工具),不用则需要注释配置里的Prometheus配置
$ brew install --build-from-source prometheus
$ 
$ homebrew 管理服务示例
$ brew services list|start|stop|restart|run|  etcd|other
$
$ 安装 go mod 依赖
$ go env -w GOPROXY=https://goproxy.cn,direct
$ go get -u
$ 上面的命令会自动安装 protoc(编译工具)、goctl(代码生成工具)
```
> 记得把 $GOPATH/bin 添加到环境变量, 以便于使用protoc(编译工具)、goctl(代码生成工具)
写代码之前建议先学习以下goctl的使用，使用命令行工具goctl生成代码可以节约大量的开发时间[自动生成api、model、dockerFile等……]
goctl文档: https://go-zero.dev/cn/goctl-commands.html

### 1.2 启动日志收集
> filebeat 采集服务器日志丢进 kafka  ->  go-stash 收集 kafka 日志丢进es  ->  我们在 kibana 上看服务器日志
##### 建议使用homebrew或者docker管理服务
>  服务依赖安装: kafka(分布式发布订阅消息系统)、elasticsearch(全文搜索引擎)、kibana(es的可视化工具)、filebeat(轻量型日志采集器)
```
$ brew install elastic/tap/elasticsearch-full
$ brew install elastic/tap/kibana-full
$ brew install elastic/tap/filebeat-full  需配置好filebeat.yml
$ brew install --build-from-source zookeeper  kafka依赖zookeeper
$ brew install --build-from-source kafka
```

##### homebrew 管理服务示例
```
$ brew services kafka|other list|start|stop|restart|run|
```

##### go-stash 使用

> 项目里面以及给大家在bin目录生成了可执行文件 stash
并在user/api/etc/other目录初始化了config.yaml(go-stash使用)、filebeat.yaml(filebeat的配置)  

如果想自己生成 stash 可执行文件，可按下面步骤执行
```
$ go clone https://github.com/tal-tech/go-stash
$ cd stash && go build stash.go
$ ./stash -f etc/config.yaml
```
> 好了，可以去 kibana 查看收集到的日志信息了

### 1.3 启动对应的服务
> 服务调用方式: rpc服务供api服务、内部其他服务调用 api服务供外部调用
```
$ go run user/api/user.go -f user/api/etc/user.yaml
$ 
$ go run user/rpc/user.go -f user/rpc/etc/user.yaml
$ 
$ go run order/api/order.go -f order/api/etc/order.yaml
```
> 启动前先改好上面的配置文件

### 1.4 homebrew 镜像设置
```
$ 设置homebrew镜像
$ cd /opt/homebrew/Library/Taps/homebrew/homebrew-core
$ 中科院
$ git remote set-url origin https://mirrors.ustc.edu.cn/homebrew-core.git
$ 阿里
$ git remote set-url origin https://mirrors.aliyun.com/homebrew/homebrew-core.git
$ 默认官方镜像
$ git remote set-url origin https://github.com/Homebrew/brew.git 
```

## 2. API文档

### 2.1 文档类型

2.1.1 go-zero 框架[go 框架]
```
go-zero 文档：https://go-zero.dev/cn/
github  地址：https://github.com/tal-tech/go-zero
```

2.1.2 govalidator[请求参数校验器]
```
参数验证规则参考文档：https://github.com/asaskevich/govalidator/blob/master/README.md
```

2.1.6 cron 文档[crontab 定时脚本管理]
```
cron 文档： https://github.com/robfig/cron
```

## 3. 表结构

### 3.1 表结构存放

表结构放在common/model/doc下，
>需注意：  有几个别的项目也在共用这个库
