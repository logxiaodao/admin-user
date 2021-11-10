package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"errors"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginOutLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewLoginOutLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *LoginOutLogic {
	return &LoginOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LoginOut 退出登陆
func (l *LoginOutLogic) LoginOut(in *user2.LoginOutRequest) (*user2.LoginOutResponse, error) {

	// 获取要使用的的错误码
	databaseError := errorx2.CodeMessage[errorx2.DatabaseError]
	accessSecret, accessExpire := l.svcCtx.Config.AuthConf.AccessSecret, l.svcCtx.Config.AuthConf.AccessExpire
	token := auth2.NewToken(accessSecret)

	tokenString, err := token.GetToken(l.ctx)
	if err != nil {
		return &user2.LoginOutResponse{Status: false}, err
	}

	err = l.svcCtx.AccountRepository.AddTokenToBlacklist(l.ctx, tokenString, accessExpire)
	if err != nil {
		return &user2.LoginOutResponse{Status: false}, errors.New(databaseError.Error())
	}

	return &user2.LoginOutResponse{Status: true}, nil
}
