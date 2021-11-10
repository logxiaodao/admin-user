package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePermissionLogic) DeletePermission(in *user2.DeletePermissionRequest) (*user2.DeletePermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.DeletePermissionResponse{Status: false}, err
	}

	// 验证ids的有效性
	idIsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.IdList)
	if err != nil {
		return &user2.DeletePermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !idIsValid {
		return &user2.DeletePermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidId)
	}

	// 验证权限是否在被角色使用，如果被使用不允许删除
	permissionUseValid, err := l.svcCtx.PermissionRepository.VerifyPermissionUse(ctx, in.IdList)
	if err != nil {
		return &user2.DeletePermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if permissionUseValid {
		return &user2.DeletePermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.PermissionInUse)
	}

	err = l.svcCtx.PermissionRepository.DeletePermission(ctx, in)
	if err != nil {
		return &user2.DeletePermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.DeletePermissionResponse{Status: true}, nil
}
