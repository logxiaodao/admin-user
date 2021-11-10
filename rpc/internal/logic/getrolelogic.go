package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  role
func (l *GetRoleLogic) GetRole(in *user2.GetRoleRequest) (*user2.GetRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.GetRoleResponse{Total: 0, RowList: nil}, err
	}
	data, err := l.svcCtx.RoleRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user2.GetRoleResponse{RowList: nil}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &data, nil
}
