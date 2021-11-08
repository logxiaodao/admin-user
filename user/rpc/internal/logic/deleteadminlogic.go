package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *user.DeleteAdminRequest) (*user.DeleteAdminResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.DeleteAdminResponse{Status: false}, err
	}
	// 验证id有效
	adminIdsValid, err := l.svcCtx.AdminRepository.VerifyAdminIds(ctx, in.IdList)
	if err != nil || !adminIdsValid {
		return &user.DeleteAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidAdminId)
	}

	err = l.svcCtx.AdminRepository.DeleteAdmin(ctx, in)
	if err != nil {
		return &user.DeleteAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.DeleteAdminResponse{Status: true}, nil
}
