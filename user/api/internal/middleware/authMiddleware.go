package middleware

import (
	"admin/user/api/internal/Initialization"
	"admin/user/api/internal/config"
	"admin/user/api/internal/pkg/util"
	"admin/user/common/errorx"
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

		dataSources := Initialization.GetDataSource()
		// 验证token是否在redis黑名单
		if dataSources.RedisClient.SIsMember(config.DefaultTokenRedisKey+":"+time.Now().Format("2006-01-02"), token).Val() {
			httpx.Error(w, errorx.SendServiceError(errorx.CodeMessage[errorx.ExpiredToken]))
			return
		}

		// 非超管不允许访问需要权限验证的接口
		isSuperAdmin := util.InterfaceToUint(r.Context().Value("isSuperAdmin"))
		if isSuperAdmin != 1 {
			httpx.Error(w, errorx.SendServiceError(errorx.CodeMessage[errorx.PermissionDenied]))
			return
		}

		next(w, r)
	}
}
