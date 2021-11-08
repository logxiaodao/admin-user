package logic

import (
	"admin/user/common/errorx"
	"admin/user/common/safe"
	"admin/user/rpc/internal/pkg/auth"
	"admin/user/rpc/internal/pkg/util"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdatePassword 修改密码
func (l *UpdatePasswordLogic) UpdatePassword(in *user.UpdatePasswordRequest) (*user.UpdatePasswordResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, false)
	if err != nil {
		return &user.UpdatePasswordResponse{Status: false}, err
	}

	// 判断密码是否一致
	if in.NewPassword != in.ConfirmPassword {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.InconsistentPasswords)
	}

	userId := util.InterfaceToInt64(ctx.Value("userId"))
	if userId == 0 {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.IncorrectToken)
	}

	// 查询用户信息
	result, err := l.svcCtx.AccountRepository.FindOneById(ctx, userId)
	if err != nil {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	// 验证旧密码是否正确
	if !safe.MatchPassword(in.OldPassword, result.Password) {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.WrongPassword)
	}

	// 生成加密的新密码
	newPassword, err := safe.GenHashPassword(in.NewPassword)
	if err != nil {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.LogicalError)
	}

	// 修改密码
	err = l.svcCtx.AccountRepository.UpdatePassword(ctx, userId, newPassword)
	if err != nil {
		return &user.UpdatePasswordResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.UpdatePasswordResponse{Status: true}, nil
}
