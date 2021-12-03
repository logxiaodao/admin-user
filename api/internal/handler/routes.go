// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"admin-user/api/internal/svc"

	"github.com/tal-tech/go-zero/rest"
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/v1/user/login",
				Handler: loginHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/account/loginOut",
					Handler: loginOutHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/account/getUserInfo",
					Handler: getUserInfoHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/account/updatePassword",
					Handler: updatePasswordHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/account/checkPermission",
					Handler: checkPermissionHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/api",
					Handler: getApiHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/findApiByIds",
					Handler: findApiByIdListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/auth/api",
					Handler: postApiHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/auth/batchApi",
					Handler: batchPostApiReqHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/auth/api",
					Handler: putApiHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/auth/api",
					Handler: deleteApiHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/admin",
					Handler: getAdminHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/findAdminByIds",
					Handler: FindAdminByIdsHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/auth/admin",
					Handler: postAdminHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/auth/admin",
					Handler: putAdminHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/auth/admin",
					Handler: deleteAdminHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/role",
					Handler: getRoleHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/findRoleByIds",
					Handler: FindRoleByIdsHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/auth/role",
					Handler: postRoleHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/auth/role",
					Handler: putRoleHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/auth/role",
					Handler: deleteRoleHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/permission",
					Handler: getPermissionHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/auth/FindPermissionByIds",
					Handler: FindPermissionByIdsHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/auth/permission",
					Handler: postPermissionHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/auth/permission",
					Handler: putPermissionHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/auth/permission",
					Handler: deletePermissionHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
