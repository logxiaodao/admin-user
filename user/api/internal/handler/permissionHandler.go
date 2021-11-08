package handler

import (
	"admin/user/common/errorx"
	"github.com/asaskevich/govalidator"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"net/http"

	"admin/user/api/internal/logic"
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func deletePermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeletePermissionReq
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

		l := logic.NewPermissionLogic(r.Context(), ctx)
		idIsValid, permissionUseValid, err := l.DeletePermission(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		if !idIsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidId]))
			return
		}

		if permissionUseValid {
			httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.PermissionInUse]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))

	}
}

func putPermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutPermissionReq
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

		l := logic.NewPermissionLogic(r.Context(), ctx)
		apiIdsValid, permissionIdsValid, err := l.PutPermission(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.PermissionExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidId]))
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidApiId]))
			return
		}
		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func postPermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostPermissionReq
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

		l := logic.NewPermissionLogic(r.Context(), ctx)
		apiIdsValid, err := l.PostPermission(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.PermissionExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidApiId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func getPermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPermissionReq
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

		l := logic.NewPermissionLogic(r.Context(), ctx)
		resp, err := l.GetPermission(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}

func FindPermissionByIdsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindPermissionByIdsReq
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

		l := logic.NewPermissionLogic(r.Context(), ctx)
		resp, err := l.FindPermissionByIds(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}
