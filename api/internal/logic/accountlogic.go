package logic

import (
	errorx2 "admin-user/api/common/errorx"
	safe2 "admin-user/api/common/safe"
	util2 "admin-user/api/internal/pkg/util"
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"context"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
)

type AccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc2.ServiceContext
}

func NewAccountLogic(ctx context.Context, svcCtx *svc2.ServiceContext) AccountLogic {
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

func (a *AccountLogic) UpdatePassword(req types2.UpdatePasswordReq) error {

	// 判断密码是否一致
	inconsistentPasswords := errorx2.CodeMessage[errorx2.InconsistentPasswords]
	if req.NewPassword != req.ConfirmPassword {
		return errors.New(inconsistentPasswords.Error())
	}

	// 获取userId
	userId := util2.InterfaceToUint(a.ctx.Value("userId"))

	result, err := a.svcCtx.AccountRepository.FindOneById(a.ctx, userId)
	if err != nil {
		return err
	}

	// 验证旧密码是否正确
	wrongPassword := errorx2.CodeMessage[errorx2.WrongPassword]
	if !safe2.MatchPassword(req.OldPassword, result.Password) {
		return errors.New(wrongPassword.Error())
	}

	newPassword, err := safe2.GenHashPassword(req.NewPassword)
	logicalError := errorx2.CodeMessage[errorx2.LogicalError]
	if err != nil {
		return errors.New(logicalError.Error())
	}

	// 修改密码
	err = a.svcCtx.AccountRepository.UpdatePassword(a.ctx, userId, newPassword)
	return err
}

func (a *AccountLogic) CheckPermission(req types2.CheckPermissionReq) (rsp types2.CheckPermissionRsp, err error) {

	// 获取userId
	userId := util2.InterfaceToUint(a.ctx.Value("userId"))

	isPass, err := a.svcCtx.AccountRepository.CheckPermission(a.ctx, userId, req.HttpMethod, req.HttpPath)

	return types2.CheckPermissionRsp{IsPass: isPass}, err
}
