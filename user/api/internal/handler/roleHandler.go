package handler

import (
	"admin/user/api/internal/logic"
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"
	"admin/user/common/errorx"
	"github.com/asaskevich/govalidator"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func deleteRoleHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.CodeMessage[errorx.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx.SendParameterError(errorx.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic.NewRoleLogic(r.Context(), ctx)
		roleIdsValid, roleUseValid, err := l.DeleteRole(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidRoleId]))
			return
		}

		if roleUseValid {
			httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.RoleInUse]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func putRoleHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.CodeMessage[errorx.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx.SendParameterError(errorx.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic.NewRoleLogic(r.Context(), ctx)
		permissionIdsValid, roleIdsValid, err := l.PutRole(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.RoleExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidRoleId]))
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidPermissionId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func getRoleHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.CodeMessage[errorx.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx.SendParameterError(errorx.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic.NewRoleLogic(r.Context(), ctx)
		resp, err := l.GetRole(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}

func FindRoleByIdsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindRoleByIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.CodeMessage[errorx.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx.SendParameterError(errorx.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic.NewRoleLogic(r.Context(), ctx)
		resp, err := l.FindRoleByIds(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}

func postRoleHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.CodeMessage[errorx.ParameterBindingFailed]))

			return
		}

		// 根据配置规则验证请求参数
		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			httpx.Error(w, errorx.SendParameterError(errorx.Msg{
				En: err.Error(),
				Zh: "参数验证不通过",
			}))
			return
		}

		l := logic.NewRoleLogic(r.Context(), ctx)
		permissionIdsValid, err := l.PostRole(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.RoleExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidPermissionId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}
