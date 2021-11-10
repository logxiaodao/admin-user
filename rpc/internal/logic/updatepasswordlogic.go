package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	safe2 "admin-user/rpc/common/safe"
	auth2 "admin-user/rpc/internal/pkg/auth"
	util2 "admin-user/rpc/internal/pkg/util"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdatePassword 修改密码
func (l *UpdatePasswordLogic) UpdatePassword(in *user2.UpdatePasswordRequest) (*user2.UpdatePasswordResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, false)
	if err != nil {
		return &user2.UpdatePasswordResponse{Status: false}, err
	}

	// 判断密码是否一致
	if in.NewPassword != in.ConfirmPassword {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InconsistentPasswords)
	}

	userId := util2.InterfaceToInt64(ctx.Value("userId"))
	if userId == 0 {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	// 查询用户信息
	result, err := l.svcCtx.AccountRepository.FindOneById(ctx, userId)
	if err != nil {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	// 验证旧密码是否正确
	if !safe2.MatchPassword(in.OldPassword, result.Password) {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.WrongPassword)
	}

	// 生成加密的新密码
	newPassword, err := safe2.GenHashPassword(in.NewPassword)
	if err != nil {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.LogicalError)
	}

	// 修改密码
	err = l.svcCtx.AccountRepository.UpdatePassword(ctx, userId, newPassword)
	if err != nil {
		return &user2.UpdatePasswordResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.UpdatePasswordResponse{Status: true}, nil
}
