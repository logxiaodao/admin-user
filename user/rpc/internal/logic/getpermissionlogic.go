package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  permission
func (l *GetPermissionLogic) GetPermission(in *user.GetPermissionRequest) (*user.GetPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.GetPermissionResponse{Total: 0, RowList: nil}, err
	}
	data, err := l.svcCtx.PermissionRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user.GetPermissionResponse{RowList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &data, nil
}
