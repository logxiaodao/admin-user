package svc

import (
	Initialization2 "admin-user/api/internal/Initialization"
	config2 "admin-user/api/internal/config"
	middleware2 "admin-user/api/internal/middleware"
	repositories2 "admin-user/api/internal/repositories"
	"log"

	"github.com/tal-tech/go-zero/rest"
)

type ServiceContext struct {
	Config               config2.Config
	Auth                 rest.Middleware
	AccountRepository    repositories2.AccountRepository
	ApiRepository        repositories2.ApiRepository
	PermissionRepository repositories2.PermissionRepository
	RoleRepository       repositories2.RoleRepository
	AdminRepository      repositories2.AdminRepository
}

func NewServiceContext(c config2.Config) *ServiceContext {
	ds, err := Initialization2.NewDataSources(c)
	if err != nil {
		log.Fatalf("new data sources error: %s", err.Error())
	}
	return &ServiceContext{
		Config:               c,
		Auth:                 middleware2.NewAuthMiddleware().Handle,
		AccountRepository:    repositories2.NewAccountRepository(ds),
		ApiRepository:        repositories2.NewApiRepository(ds),
		PermissionRepository: repositories2.NewPermissionRepository(ds),
		RoleRepository:       repositories2.NewRoleRepository(ds),
		AdminRepository:      repositories2.NewAdminRepository(ds),
	}
}
