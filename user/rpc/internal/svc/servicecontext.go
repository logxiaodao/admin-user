package svc

import (
	"admin/user/rpc/internal/Initialization"
	"admin/user/rpc/internal/config"
	"admin/user/rpc/internal/repositories"
	"log"
)

type ServiceContext struct {
	Config               config.Config
	AccountRepository    repositories.AccountRepository
	ApiRepository        repositories.ApiRepository
	AdminRepository      repositories.AdminRepository
	RoleRepository       repositories.RoleRepository
	PermissionRepository repositories.PermissionRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	ds, err := Initialization.NewDataSources(c)
	if err != nil {
		log.Fatalf("new data sources error: %s", err.Error())
	}
	return &ServiceContext{
		Config:               c,
		AccountRepository:    repositories.NewAccountRepository(ds),
		ApiRepository:        repositories.NewApiRepository(ds),
		AdminRepository:      repositories.NewAdminRepository(ds),
		RoleRepository:       repositories.NewRoleRepository(ds),
		PermissionRepository: repositories.NewPermissionRepository(ds),
	}
}
