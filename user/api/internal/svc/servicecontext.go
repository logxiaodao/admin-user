package svc

import (
	"log"

	"admin/user/api/internal/Initialization"
	"admin/user/api/internal/config"
	"admin/user/api/internal/middleware"
	"admin/user/api/internal/repositories"

	"github.com/tal-tech/go-zero/rest"
)

type ServiceContext struct {
	Config               config.Config
	Auth                 rest.Middleware
	AccountRepository    repositories.AccountRepository
	ApiRepository        repositories.ApiRepository
	PermissionRepository repositories.PermissionRepository
	RoleRepository       repositories.RoleRepository
	AdminRepository      repositories.AdminRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	ds, err := Initialization.NewDataSources(c)
	if err != nil {
		log.Fatalf("new data sources error: %s", err.Error())
	}
	return &ServiceContext{
		Config:               c,
		Auth:                 middleware.NewAuthMiddleware().Handle,
		AccountRepository:    repositories.NewAccountRepository(ds),
		ApiRepository:        repositories.NewApiRepository(ds),
		PermissionRepository: repositories.NewPermissionRepository(ds),
		RoleRepository:       repositories.NewRoleRepository(ds),
		AdminRepository:      repositories.NewAdminRepository(ds),
	}
}
