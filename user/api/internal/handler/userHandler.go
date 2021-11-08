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

func loginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
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

		l := logic.NewUserLogic(r.Context(), ctx)
		resp, err := l.Login(req)
		if err != nil {
			httpx.Error(w, errorx.SendPermissionError(errorx.Msg{
				En: err.Error(),
				Zh: "权限错误",
			}))
		} else {
			httpx.OkJson(w, errorx.SendSuccess(resp))
		}
	}
}
