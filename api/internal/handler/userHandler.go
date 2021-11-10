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

func loginHandler(ctx *svc2.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types2.LoginReq
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

		l := logic2.NewUserLogic(r.Context(), ctx)
		resp, err := l.Login(req)
		if err != nil {
			httpx.Error(w, errorx2.SendPermissionError(errorx2.Msg{
				En: err.Error(),
				Zh: "权限错误",
			}))
		} else {
			httpx.OkJson(w, errorx2.SendSuccess(resp))
		}
	}
}
