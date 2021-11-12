package middleware

import (
	errorx2 "admin-user/api/common/errorx"
	Initialization2 "admin-user/api/internal/Initialization"
	config2 "admin-user/api/internal/config"
	util2 "admin-user/api/internal/pkg/util"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"time"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		dataSources := Initialization2.GetDataSource()
		// 验证token是否在redis黑名单
		if dataSources.RedisClient.SIsMember(config2.DefaultTokenRedisKey+":"+time.Now().Format("2006-01-02"), token).Val() {
			httpx.Error(w, errorx2.SendServiceError(errorx2.CodeMessage[errorx2.ExpiredToken]))
			return
		}

		// 判读是否公共接口
		for _, v := range config2.SecurityApiData {
			if v.HTTPPath == r.RequestURI && v.HTTPMethod == r.Method {
				next(w, r)
			}
		}

		// 非超管不允许访问需要权限验证的接口
		isSuperAdmin := util2.InterfaceToUint(r.Context().Value("isSuperAdmin"))
		if isSuperAdmin != 1 {
			httpx.Error(w, errorx2.SendServiceError(errorx2.CodeMessage[errorx2.PermissionDenied]))
			return
		}

		next(w, r)
	}
}
