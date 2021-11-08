package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"errors"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginOutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginOutLogic {
	return &LoginOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LoginOut 退出登陆
func (l *LoginOutLogic) LoginOut(in *user.LoginOutRequest) (*user.LoginOutResponse, error) {

	// 获取要使用的的错误码
	databaseError := errorx.CodeMessage[errorx.DatabaseError]
	accessSecret, accessExpire := l.svcCtx.Config.AuthConf.AccessSecret, l.svcCtx.Config.AuthConf.AccessExpire
	token := auth.NewToken(accessSecret)

	tokenString, err := token.GetToken(l.ctx)
	if err != nil {
		return &user.LoginOutResponse{Status: false}, err
	}

	err = l.svcCtx.AccountRepository.AddTokenToBlacklist(l.ctx, tokenString, accessExpire)
	if err != nil {
		return &user.LoginOutResponse{Status: false}, errors.New(databaseError.Error())
	}

	return &user.LoginOutResponse{Status: true}, nil
}
