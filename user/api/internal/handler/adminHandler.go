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

func getAdminHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAdminReq
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

		l := logic.NewAdminLogic(r.Context(), ctx)
		resp, err := l.GetAdmin(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}

func FindAdminByIdsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindAdminByIdsReq
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

		l := logic.NewAdminLogic(r.Context(), ctx)
		resp, err := l.FindAdminByIds(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}

func postAdminHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostAdminReq
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

		l := logic.NewAdminLogic(r.Context(), ctx)
		roleIdsValid, err := l.PostAdmin(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.AdminExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidRoleId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func putAdminHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutAdminReq
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

		l := logic.NewAdminLogic(r.Context(), ctx)
		roleIdsValid, adminIdsValid, err := l.PutAdmin(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.AdminExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !adminIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidId]))
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidRoleId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func deleteAdminHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAdminReq
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

		l := logic.NewAdminLogic(r.Context(), ctx)
		adminIdsValid, err := l.DeleteAdmin(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		if !adminIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidAdminId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))

	}
}
