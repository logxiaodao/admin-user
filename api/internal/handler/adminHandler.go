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

func getAdminHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.GetAdminReq
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

		l := logic2.NewAdminLogic(r.Context(), ctx)
		resp, err := l.GetAdmin(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}

func FindAdminByIdsHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.FindAdminByIdsReq
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

		l := logic2.NewAdminLogic(r.Context(), ctx)
		resp, err := l.FindAdminByIds(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}

func postAdminHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PostAdminReq
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

		l := logic2.NewAdminLogic(r.Context(), ctx)
		roleIdsValid, err := l.PostAdmin(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.AdminExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidRoleId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func putAdminHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PutAdminReq
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

		l := logic2.NewAdminLogic(r.Context(), ctx)
		roleIdsValid, adminIdsValid, err := l.PutAdmin(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.AdminExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !adminIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidId]))
			return
		}

		if !roleIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidRoleId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func deleteAdminHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.DeleteAdminReq
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

		l := logic2.NewAdminLogic(r.Context(), ctx)
		adminIdsValid, err := l.DeleteAdmin(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		if !adminIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidAdminId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))

	}
}
