package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	"admin-user/rpc/internal/pkg/util"
	"context"

	"admin-user/rpc/internal/svc"
	"admin-user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  account
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {

	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, false)
	if err != nil {
		return &user.GetUserInfoResponse{}, err
	}

	roleList, err := l.svcCtx.AdminRepository.FindAdminByIdList(ctx, []string{util.InterfaceToString(ctx.Value("userId"))})
	if err != nil {
		return &user.GetUserInfoResponse{}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user.GetUserInfoResponse{
		Id:         roleList[0].Id,
		Account:    roleList[0].Account,
		NickName:   roleList[0].NickName,
		Phone:      roleList[0].Phone,
		Email:      roleList[0].Email,
		Createdat:  roleList[0].Createdat,
		Updatedat:  roleList[0].Updatedat,
		PlatformId: util.InterfaceToString(ctx.Value("platformID"))}, nil
}
