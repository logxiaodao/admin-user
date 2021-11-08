package logic

import (
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"
	"admin/user/common/safe"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/test/model"
	"time"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserLogic {
	return UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (a *UserLogic) Login(req types.LoginReq) (*types.LoginRes, error) {

	userInfo, err := a.svcCtx.AccountRepository.FindOneByAccount(a.ctx, req.Username, req.PlatformID)

	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("User does not exist")
	default:
		return nil, err
	}
	if !safe.MatchPassword(req.Password, userInfo.Password) {
		return nil, errors.New("Incorrect username or password")
	}

	// 使用jwt
	now := time.Now().Unix()
	accessExpire := a.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := a.getJwtToken(a.svcCtx.Config.Auth.AccessSecret, now, a.svcCtx.Config.Auth.AccessExpire, int64(userInfo.ID), int64(req.PlatformID), int64(userInfo.IsSuperAdmin))
	if err != nil {
		return nil, err
	}

	return &types.LoginRes{
		Id:           int64(userInfo.ID),
		PlatformID:   int64(req.PlatformID),
		Name:         userInfo.Account,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (a *UserLogic) getJwtToken(secretKey string, iat, seconds, userId, platformID, isSuperAdmin int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["platformID"] = platformID
	claims["isSuperAdmin"] = isSuperAdmin
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))

}
