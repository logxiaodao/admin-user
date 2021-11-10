package handler

import (
	errorx2 "admin-user/api/common/errorx"
	logic2 "admin-user/api/internal/logic"
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"github.com/asaskevich/govalidator"
	"github.com/go-sql-driver/mysql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func getApiHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.GetApiReq
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		resp, err := l.GetApi(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(resp))
	}
}

func findApiByIdListHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.FindApiByIdsReq
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		resp, err := l.FindApiByIdList(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(resp))
	}
}

func deleteApiHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.DeleteApiReq
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		apiIdsValid, apiUseValid, err := l.DeleteApi(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidId]))
			return
		}

		if apiUseValid {
			httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.ApiInUse]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func postApiHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PostApiReq
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		err = l.PostApi(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.ApiExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func batchPostApiReqHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.BatchPostApiReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		err = l.BatchPostApiReq(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.ApiExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}

func putApiHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.PutApiReq
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

		l := logic2.NewApiLogic(r.Context(), ctx)
		apiIdsValid, err := l.PutApi(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx2.SendLogicalError(errorx2.CodeMessage[errorx2.ApiExists]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[errorx2.InvalidId]))
			return
		}

		httpx.OkJson(w, errorx2.SendSuccess(nil))
	}
}
