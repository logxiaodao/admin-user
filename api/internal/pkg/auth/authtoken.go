package auth

import (
	errorx2 "admin-user/api/common/errorx"
	util2 "admin-user/api/internal/pkg/util"
	"context"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	SuperAdmin   = true
	NoSuperAdmin = false
)

type Token interface {
	ParseHS256Token(tokenString string) (jwt.MapClaims, error)
	GetTokenClaims(ctx context.Context) (jwt.MapClaims, error)
}

type token struct {
	AccessSecret string
}

func NewToken(accessSecret string) *token {
	return &token{AccessSecret: accessSecret}
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
func (t *token) GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	if tokenString[:7] == "Bearer " { // 去掉头部 Bearer
		tokenString = tokenString[7:]
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
