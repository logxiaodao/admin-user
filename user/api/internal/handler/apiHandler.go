package handler

import (
	"admin/user/common/errorx"
	"github.com/asaskevich/govalidator"
	"github.com/go-sql-driver/mysql"
	"net/http"

	"admin/user/api/internal/logic"
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func getApiHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetApiReq
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

		l := logic.NewApiLogic(r.Context(), ctx)
		resp, err := l.GetApi(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(resp))
	}
}

func findApiByIdListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindApiByIdsReq
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

		l := logic.NewApiLogic(r.Context(), ctx)
		resp, err := l.FindApiByIdList(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(resp))
	}
}

func deleteApiHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteApiReq
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

		l := logic.NewApiLogic(r.Context(), ctx)
		apiIdsValid, apiUseValid, err := l.DeleteApi(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidId]))
			return
		}

		if apiUseValid {
			httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.ApiInUse]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func postApiHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostApiReq
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

		l := logic.NewApiLogic(r.Context(), ctx)
		err = l.PostApi(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.ApiExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func batchPostApiReqHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchPostApiReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
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

		l := logic.NewApiLogic(r.Context(), ctx)
		err = l.BatchPostApiReq(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.ApiExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}

func putApiHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutApiReq
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

		l := logic.NewApiLogic(r.Context(), ctx)
		apiIdsValid, err := l.PutApi(req)
		if err != nil {
			var sqlErr *mysql.MySQLError
			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
				httpx.Error(w, errorx.SendLogicalError(errorx.CodeMessage[errorx.ApiExists]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
			return
		}

		if !apiIdsValid {
			httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[errorx.InvalidId]))
			return
		}

		httpx.OkJson(w, errorx.SendSuccess(nil))
	}
}
