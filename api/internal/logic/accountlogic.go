package logic

import (
	errorx2 "admin-user/api/common/errorx"
	safe2 "admin-user/api/common/safe"
	"admin-user/api/internal/config"
	"admin-user/api/internal/pkg/auth"
	util2 "admin-user/api/internal/pkg/util"
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"context"
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

func (a *AccountLogic) GetUserInfo(token string) (*types2.GetUserInfoRsp, error) {

	tk := auth.NewToken(a.svcCtx.Config.Auth.AccessSecret)
	userInfo, err := tk.GetTokenClaims(token)
	if err != nil {
		return &types2.GetUserInfoRsp{}, err
	}

	adminList, err := a.svcCtx.AdminRepository.FindAdminByIdList(a.ctx, []string{util2.InterfaceToString(userInfo["userId"])})
	if err != nil {
		return &types2.GetUserInfoRsp{}, err
	}

	return &types2.GetUserInfoRsp{
		Id:         int64(adminList[0].Id),
		Account:    adminList[0].Account,
		NickName:   adminList[0].NickName,
		Phone:      adminList[0].Phone,
		Email:      adminList[0].Email,
		PlatformID: util2.InterfaceToint64(userInfo["platformID"]),
	}, nil
}

func (a *AccountLogic) UpdatePassword(req types2.UpdatePasswordReq) error {

	// 判断密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		return errorx2.GetErrorByCode(errorx2.InconsistentPasswords)
	}

	// 获取userId
	userId := util2.InterfaceToUint(a.ctx.Value("userId"))

	result, err := a.svcCtx.AccountRepository.FindOneById(a.ctx, userId)
	if err != nil {
		return err
	}

	// 验证旧密码是否正确
	if !safe2.MatchPassword(req.OldPassword, result.Password) {
		return errorx2.GetErrorByCode(errorx2.WrongPassword)
	}

	newPassword, err := safe2.GenHashPassword(req.NewPassword)
	if err != nil {
		return errorx2.GetErrorByCode(errorx2.LogicalError)
	}

	// 修改密码
	err = a.svcCtx.AccountRepository.UpdatePassword(a.ctx, userId, newPassword)
	return err
}

func (a *AccountLogic) CheckPermission(req types2.CheckPermissionReq) (rsp types2.CheckPermissionRsp, err error) {

	// 获取userId
	userId := util2.InterfaceToUint(a.ctx.Value("userId"))

	// 判读是否公共接口
	for _, v := range config.InitRouteData.Data {
		if v.HTTPPath == req.HttpPath && v.HTTPMethod == req.HttpMethod && v.IsOpen == 1 {
			return types2.CheckPermissionRsp{IsPass: true}, nil
		}
	}

	isPass, err := a.svcCtx.AccountRepository.CheckPermission(a.ctx, userId, req.HttpMethod, req.HttpPath)

	return types2.CheckPermissionRsp{IsPass: isPass}, err
}
