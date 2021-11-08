package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePermissionLogic) DeletePermission(in *user.DeletePermissionRequest) (*user.DeletePermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.DeletePermissionResponse{Status: false}, err
	}

	// 验证ids的有效性
	idIsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.IdList)
	if err != nil {
		return &user.DeletePermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !idIsValid {
		return &user.DeletePermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidId)
	}

	// 验证权限是否在被角色使用，如果被使用不允许删除
	permissionUseValid, err := l.svcCtx.PermissionRepository.VerifyPermissionUse(ctx, in.IdList)
	if err != nil {
		return &user.DeletePermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if permissionUseValid {
		return &user.DeletePermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.PermissionInUse)
	}

	err = l.svcCtx.PermissionRepository.DeletePermission(ctx, in)
	if err != nil {
		return &user.DeletePermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.DeletePermissionResponse{Status: true}, nil
}
