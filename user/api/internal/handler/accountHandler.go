package handler

import (
	"admin/user/common/errorx"
	"github.com/asaskevich/govalidator"
	"net/http"

	"admin/user/api/internal/logic"
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func loginOutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewAccountLogic(r.Context(), ctx)
		err := l.LoginOut(r.Header.Get("Authorization"))
		if err != nil {
			httpx.Error(w, errorx.SendServiceError(errorx.DefaultCodeMessage[errorx.ServiceError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(nil))
		}
	}
}

func updatePasswordHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePasswordReq
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

		l := logic.NewAccountLogic(r.Context(), ctx)
		err = l.UpdatePassword(req)
		if err != nil {
			code := errorx.FindCodeByMsg(err.Error())
			if code != -1 {
				httpx.Error(w, errorx.SendParameterError(errorx.CodeMessage[code]))
			} else {
				httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
			}
		} else {
			httpx.OkJson(w, errorx.SendSuccess(nil))
		}
	}
}

func checkPermissionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckPermissionReq
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

		l := logic.NewAccountLogic(r.Context(), ctx)
		rsp, err := l.CheckPermission(req)
		if err != nil {
			httpx.Error(w, errorx.SendDatabaseError(errorx.DefaultCodeMessage[errorx.DatabaseError]))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(rsp))
		}
	}
}
