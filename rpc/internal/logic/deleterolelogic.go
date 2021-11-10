package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *user2.DeleteRoleRequest) (*user2.DeleteRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.DeleteRoleResponse{Status: false}, err
	}

	// 验证id有效
	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.IdList)
	if err != nil {
		return &user2.DeleteRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !roleIdsValid {
		return &user2.DeleteRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidId)
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	roleUseValid, err := l.svcCtx.RoleRepository.VerifyRoleUse(ctx, in.IdList)
	if err != nil {
		return &user2.DeleteRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if roleUseValid {
		return &user2.DeleteRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.ApiInUse)
	}

	err = l.svcCtx.RoleRepository.DeleteRole(ctx, in)
	if err != nil {
		return &user2.DeleteRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.DeleteRoleResponse{Status: true}, nil
}
