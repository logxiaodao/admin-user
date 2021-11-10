package handler

import (
	errorx2 "admin-user/api/common/errorx"
	logic2 "admin-user/api/internal/logic"
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"github.com/asaskevich/govalidator"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func deleteRoleHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.DeleteRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.CodeMessage[errorx2.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx2.SendParameterError(errorx2.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic2.NewRoleLogic(r.Context(), ctx)
		roleIdsValid, roleUseValid, err := l.DeleteRole(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidRoleId]))
			return
		}

		if roleUseValid {
			httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.RoleInUse]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func putRoleHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PutRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.CodeMessage[errorx2.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx2.SendParameterError(errorx2.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic2.NewRoleLogic(r.Context(), ctx)
		permissionIdsValid, roleIdsValid, err := l.PutRole(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.RoleExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidRoleId]))
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidPermissionId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func getRoleHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.GetRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.CodeMessage[errorx2.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx2.SendParameterError(errorx2.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic2.NewRoleLogic(r.Context(), ctx)
		resp, err := l.GetRole(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}

func FindRoleByIdsHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.FindRoleByIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.CodeMessage[errorx2.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx2.SendParameterError(errorx2.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic2.NewRoleLogic(r.Context(), ctx)
		resp, err := l.FindRoleByIds(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}

func postRoleHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PostRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.CodeMessage[errorx2.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx2.SendParameterError(errorx2.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic2.NewRoleLogic(r.Context(), ctx)
		permissionIdsValid, err := l.PostRole(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.RoleExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidPermissionId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}
