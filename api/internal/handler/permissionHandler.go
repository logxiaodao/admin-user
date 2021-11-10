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

func deletePermissionHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.DeletePermissionReq
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

		l := logic2.NewPermissionLogic(r.Context(), ctx)
		idIsValid, permissionUseValid, err := l.DeletePermission(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		if !idIsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidId]))
			return
		}

		if permissionUseValid {
			httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.PermissionInUse]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))

	}
}

func putPermissionHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PutPermissionReq
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

		l := logic2.NewPermissionLogic(r.Context(), ctx)
		apiIdsValid, permissionIdsValid, err := l.PutPermission(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.PermissionExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !permissionIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidId]))
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidApiId]))
			return
		}
		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func postPermissionHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PostPermissionReq
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

		l := logic2.NewPermissionLogic(r.Context(), ctx)
		apiIdsValid, err := l.PostPermission(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.PermissionExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidApiId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func getPermissionHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.GetPermissionReq
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

		l := logic2.NewPermissionLogic(r.Context(), ctx)
		resp, err := l.GetPermission(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}

func FindPermissionByIdsHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.FindPermissionByIdsReq
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

		l := logic2.NewPermissionLogic(r.Context(), ctx)
		resp, err := l.FindPermissionByIds(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}
