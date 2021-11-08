package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminLogic {
	return &GetAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  admin
func (l *GetAdminLogic) GetAdmin(in *user.GetAdminRequest) (*user.GetAdminResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.GetAdminResponse{Total: 0, RowList: nil}, err
	}

	data, err := l.svcCtx.AdminRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user.GetAdminResponse{RowList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &data, nil
}
