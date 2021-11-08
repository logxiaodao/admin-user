package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  role
func (l *GetRoleLogic) GetRole(in *user.GetRoleRequest) (*user.GetRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.GetRoleResponse{Total: 0, RowList: nil}, err
	}
	data, err := l.svcCtx.RoleRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user.GetRoleResponse{RowList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &data, nil
}
