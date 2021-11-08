package logic

import (
	"admin/user/common/safe"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/test/model"
	"time"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {

	userInfo, err := l.svcCtx.AccountRepository.FindOneByAccount(l.ctx, in.Username, in.PlatformID)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("User does not exist")
	default:
		return nil, err
	}
	if !safe.MatchPassword(in.Password, userInfo.Password) {
		return nil, errors.New("Incorrect username or password")
	}

	// 使用jwt
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.AuthConf.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.AuthConf.AccessSecret, now, accessExpire, int64(userInfo.ID), in.PlatformID, int64(userInfo.IsSuperAdmin))
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Id:           int64(userInfo.ID),
		PlatformID:   in.PlatformID,
		Name:         userInfo.Account,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId, platformID, isSuperAdmin int64) (string, error) {
	claims := make(jwt.MapClaims)
	// 这里的参数是写死的 改动需要改其他使用的地方
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["platformID"] = platformID
	claims["isSuperAdmin"] = isSuperAdmin
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
