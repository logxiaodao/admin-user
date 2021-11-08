package logic

import (
	"admin/user/api/internal/pkg/util"
	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"
	"admin/user/common/errorx"
	"admin/user/common/safe"
	"context"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
)

type AccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) AccountLogic {
	return AccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// LoginOut 维护一个redis token黑名单
func (a *AccountLogic) LoginOut(token string) (err error) {

	accessExpire := a.svcCtx.Config.Auth.AccessExpire

	err = a.svcCtx.AccountRepository.AddTokenToBlacklist(a.ctx, token, accessExpire)

	return
}

func (a *AccountLogic) UpdatePassword(req types.UpdatePasswordReq) error {

	// 判断密码是否一致
	inconsistentPasswords := errorx.CodeMessage[errorx.InconsistentPasswords]
	if req.NewPassword != req.ConfirmPassword {
		return errors.New(inconsistentPasswords.Error())
	}

	// 获取userId
	userId := util.InterfaceToUint(a.ctx.Value("userId"))

	result, err := a.svcCtx.AccountRepository.FindOneById(a.ctx, userId)
	if err != nil {
		return err
	}

	// 验证旧密码是否正确
	wrongPassword := errorx.CodeMessage[errorx.WrongPassword]
	if !safe.MatchPassword(req.OldPassword, result.Password) {
		return errors.New(wrongPassword.Error())
	}

	newPassword, err := safe.GenHashPassword(req.NewPassword)
	logicalError := errorx.CodeMessage[errorx.LogicalError]
	if err != nil {
		return errors.New(logicalError.Error())
	}

	// 修改密码
	err = a.svcCtx.AccountRepository.UpdatePassword(a.ctx, userId, newPassword)
	return err
}

func (a *AccountLogic) CheckPermission(req types.CheckPermissionReq) (rsp types.CheckPermissionRsp, err error) {

	// 获取userId
	userId := util.InterfaceToUint(a.ctx.Value("userId"))

	isPass, err := a.svcCtx.AccountRepository.CheckPermission(a.ctx, userId, req.HttpMethod, req.HttpPath)

	return types.CheckPermissionRsp{IsPass: isPass}, err
}
