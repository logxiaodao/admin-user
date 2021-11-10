package handler

import (
	errorx2 "admin-user/api/common/errorx"
	logic2 "admin-user/api/internal/logic"
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"github.com/asaskevich/govalidator"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func loginOutHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic2.NewAccountLogic(r.Context(), ctx)
		err := l.LoginOut(r.Header.Get("Authorization"))
		if err != nil {
			httpx.Error(w, errorx2.SendServiceError(errorx2.DefaultCodeMessage[errorx2.ServiceError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(nil))
		}
	}
}

func updatePasswordHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.UpdatePasswordReq
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

		l := logic2.NewAccountLogic(r.Context(), ctx)
		err = l.UpdatePassword(req)
		if err != nil {
			code := errorx2.FindCodeByMsg(err.Error())
			if code != -1 {
				httpx.Error(w, errorx2.SendParameterError(errorx2.CodeMessage[code]))
			} else {
				httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
			}
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(nil))
		}
	}
}

func checkPermissionHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.CheckPermissionReq
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

		l := logic2.NewAccountLogic(r.Context(), ctx)
		rsp, err := l.CheckPermission(req)
		if err != nil {
			httpx.Error(w, errorx2.SendDatabaseError(errorx2.DefaultCodeMessage[errorx2.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(rsp))
		}
	}
}
