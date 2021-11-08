package config

import (
	"admin/user/api/internal/model"
	"net/http"
)

type InitDataConf struct {
	InitAdminConf
	InitPlatformConf
	InitRouteConf
}

// 定义数据格式
type (
	// InitRouteConf 路由初始化
	InitRouteConf struct {
		PlatformId uint
		Data       []model.AdminAPI
	}
	// InitPlatformConf 平台初始化
	InitPlatformConf struct {
		Data []model.AdminPlatform
	}
	// InitAdminConf 超级管理员初始化
	InitAdminConf struct {
		Data []AdminList
	}
	AdminList struct {
		PlatformIds []int
		Account     string
		Password    string
		NickName    string
		Phone       string
		Email       string
	}
)

const (
	DefaultPassword      = "12345678"
	AdminPlatformId      = 999999999
	DefaultTokenRedisKey = "system:token:blacklist"
)

// 数据初始化
var (
	// InitAdminData 超级管理员数据初始化(管理后台)
	InitAdminData = InitAdminConf{
		Data: []AdminList{
			{
				Account:  "lxj",
				Password: "12345678",
				NickName: "lxj",
				Phone:    "18598814566",
				Email:    "longxiangjun@laihua.com",
			},
		},
	}

	// InitPlatformData 平台数据初始化(管理后台)
	InitPlatformData = InitPlatformConf{
		Data: []model.AdminPlatform{
			{ // 系统自身,如果改动这里则必须改动下面的路由数据
				ID:         1,
				PlatformEn: "LaiHuaUserMicroservice",
				PlatformZh: "来画用户微服务",
			},
			{
				ID:         2,
				PlatformEn: "LaiHuaDemo",
				PlatformZh: "来画演示管理后台",
			},
			{
				ID:         3,
				PlatformEn: "LaiHuaDesign",
				PlatformZh: "来画设计管理后台",
			},
			{
				ID:         4,
				PlatformEn: "LaiHuaAnimation",
				PlatformZh: "来画动画管理后台",
			},
		},
	}

	// InitRouteData 路由数据初始化(这里只初始化自身的路由,其他平台需调用批量添加api接口)
	InitRouteData = InitRouteConf{
		PlatformId: 1,
		Data: []model.AdminAPI{
			{
				Name:       "分页查询api列表",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/api",
			},
			{
				Name:       "使用id列表来查询api信息",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/findApiByIds",
			},
			{
				Name:       "添加api",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/auth/api",
			},
			{
				Name:       "批量添加api",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/auth/batchApi",
			},
			{
				Name:       "修改api",
				HTTPMethod: http.MethodPut,
				HTTPPath:   "/v1/auth/api",
			},
			{
				Name:       "删除api(包含批量)",
				HTTPMethod: http.MethodDelete,
				HTTPPath:   "/v1/auth/api",
			},
			{
				Name:       "分页获取权限",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/permission",
			},
			{
				Name:       "通过id查询permission信息",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/FindPermissionByIds",
			},
			{
				Name:       "添加权限",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/auth/permission",
			},
			{
				Name:       "修改权限",
				HTTPMethod: http.MethodPut,
				HTTPPath:   "/v1/auth/permission",
			},
			{
				Name:       "删除权限(包含批量)",
				HTTPMethod: http.MethodDelete,
				HTTPPath:   "/v1/auth/permission",
			},
			{
				Name:       "分页获取角色",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/role",
			},
			{
				Name:       "使用id列表来查询角色信息",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/FindRoleByIds",
			},
			{
				Name:       "添加角色",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/auth/role",
			},
			{
				Name:       "修改角色",
				HTTPMethod: http.MethodPut,
				HTTPPath:   "/v1/auth/role",
			},
			{
				Name:       "删除角色(包含批量)",
				HTTPMethod: http.MethodDelete,
				HTTPPath:   "/v1/auth/role",
			},
			{
				Name:       "分页获取管理员用户列表",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/admin",
			},
			{
				Name:       "使用id列表来查询用户信息",
				HTTPMethod: http.MethodGet,
				HTTPPath:   "/v1/auth/FindAdminByIds",
			},
			{
				Name:       "添加管理员用户",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/auth/admin",
			},
			{
				Name:       "修改管理员用户",
				HTTPMethod: http.MethodPut,
				HTTPPath:   "/v1/auth/admin",
			},
			{
				Name:       "删除管理员用户(包含批量)",
				HTTPMethod: http.MethodDelete,
				HTTPPath:   "/v1/auth/admin",
			},
			{
				Name:       "登陆",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/user/login",
			},
			{
				Name:       "修改用户密码",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/account/updatePassword",
			},
			{
				Name:       "退出登陆",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/account/loginOut",
			},
			{
				Name:       "校验api接口权限",
				HTTPMethod: http.MethodPost,
				HTTPPath:   "/v1/account/checkPermission",
			},
		},
	}
)
