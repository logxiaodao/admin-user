package auth

import (
	errorx2 "admin-user/rpc/common/errorx"
	Initialization2 "admin-user/rpc/internal/Initialization"
	config2 "admin-user/rpc/internal/config"
	util2 "admin-user/rpc/internal/pkg/util"
	"context"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	SuperAdmin   = true
	NoSuperAdmin = false
)

type Token interface {
	GetToken(ctx context.Context) (string, error)
	ParseHS256Token(tokenString string) (jwt.MapClaims, error)
	GetTokenClaims(ctx context.Context) (jwt.MapClaims, error)
	SetContext(ctx context.Context, isSuperAdmin bool) (newCtx context.Context, err error)
}

type token struct {
	AccessSecret string
}

func NewToken(accessSecret string) Token {
	return &token{AccessSecret: accessSecret}
}

// GetToken 获取token
func (t *token) GetToken(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	tokenList, ok := md["authorization"]
	if !ok || len(tokenList) == 0 {
		return "", errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	if tokenList[0][:7] == "Bearer " { // 去掉头部 Bearer
		tokenList[0] = tokenList[0][7:]
	}

	dataSources := Initialization2.GetDataSource()
	// 验证token是否在redis黑名单
	if dataSources.RedisClient.SIsMember(config2.DefaultTokenRedisKey+":"+time.Now().Format("2006-01-02"), tokenList[0]).Val() {
		return "", errorx2.GetErrorByCode(errorx2.ExpiredToken)
	}

	return tokenList[0], nil
}

// ParseHS256Token 解密 jwt
func (t *token) ParseHS256Token(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.AccessSecret), nil
	})
	if err != nil {
		return nil, errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	return claims, nil
}

// GetTokenClaims 获取token 内容
func (t *token) GetTokenClaims(ctx context.Context) (jwt.MapClaims, error) {

	tokenString, err := t.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	userInfo, err := t.ParseHS256Token(tokenString)
	if err != nil {
		return nil, err
	}

	if _, ok := userInfo["exp"]; ok {
		if int64(util2.InterfaceToUint(userInfo["exp"])) < time.Now().Unix() {
			return nil, errorx2.GetErrorByCode(errorx2.ExpiredToken)
		}
	}

	return userInfo, nil
}

// SetContext 将用户信息加入上下文
func (t *token) SetContext(ctx context.Context, isSuperAdmin bool) (newCtx context.Context, err error) {
	userInfo, err := t.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}

	if isSuperAdmin {
		// 拦截非超管
		if _, ok := userInfo["isSuperAdmin"]; !ok {
			return nil, errorx2.GetErrorByCode(errorx2.PermissionDenied)
		}
		if util2.InterfaceToUint(userInfo["isSuperAdmin"]) != uint(1) {
			return nil, errorx2.GetErrorByCode(errorx2.PermissionDenied)
		}
	}

	// 将用户信息存入上下文
	for k, v := range userInfo {
		ctx = context.WithValue(ctx, k, v)
	}
	newCtx = ctx

	return
}
