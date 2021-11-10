package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  permission
func (l *GetPermissionLogic) GetPermission(in *user2.GetPermissionRequest) (*user2.GetPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.GetPermissionResponse{Total: 0, RowList: nil}, err
	}
	data, err := l.svcCtx.PermissionRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user2.GetPermissionResponse{RowList: nil}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &data, nil
}
