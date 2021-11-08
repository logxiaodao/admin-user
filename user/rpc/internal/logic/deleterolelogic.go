package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *user.DeleteRoleRequest) (*user.DeleteRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.DeleteRoleResponse{Status: false}, err
	}

	// 验证id有效
	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.IdList)
	if err != nil {
		return &user.DeleteRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !roleIdsValid {
		return &user.DeleteRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidId)
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	roleUseValid, err := l.svcCtx.RoleRepository.VerifyRoleUse(ctx, in.IdList)
	if err != nil {
		return &user.DeleteRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if roleUseValid {
		return &user.DeleteRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.ApiInUse)
	}

	err = l.svcCtx.RoleRepository.DeleteRole(ctx, in)
	if err != nil {
		return &user.DeleteRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.DeleteRoleResponse{Status: true}, nil
}
