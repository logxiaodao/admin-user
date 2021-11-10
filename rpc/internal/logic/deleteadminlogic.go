package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *user2.DeleteAdminRequest) (*user2.DeleteAdminResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.DeleteAdminResponse{Status: false}, err
	}
	// 验证id有效
	adminIdsValid, err := l.svcCtx.AdminRepository.VerifyAdminIds(ctx, in.IdList)
	if err != nil || !adminIdsValid {
		return &user2.DeleteAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidAdminId)
	}

	err = l.svcCtx.AdminRepository.DeleteAdmin(ctx, in)
	if err != nil {
		return &user2.DeleteAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.DeleteAdminResponse{Status: true}, nil
}
