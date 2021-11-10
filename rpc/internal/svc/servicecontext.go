package svc

import (
	Initialization2 "admin-user/rpc/internal/Initialization"
	config2 "admin-user/rpc/internal/config"
	repositories2 "admin-user/rpc/internal/repositories"
	"log"
)

type ServiceContext struct {
	Config               config2.Config
	AccountRepository    repositories2.AccountRepository
	ApiRepository        repositories2.ApiRepository
	AdminRepository      repositories2.AdminRepository
	RoleRepository       repositories2.RoleRepository
	PermissionRepository repositories2.PermissionRepository
}

func NewServiceContext(c config2.Config) *ServiceContext {
	ds, err := Initialization2.NewDataSources(c)
	if err != nil {
		log.Fatalf("new data sources error: %s", err.Error())
	}
	return &ServiceContext{
		Config:               c,
		AccountRepository:    repositories2.NewAccountRepository(ds),
		ApiRepository:        repositories2.NewApiRepository(ds),
		AdminRepository:      repositories2.NewAdminRepository(ds),
		RoleRepository:       repositories2.NewRoleRepository(ds),
		PermissionRepository: repositories2.NewPermissionRepository(ds),
	}
}
