package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddRoleLogic) AddRole(in *user.AddRoleRequest) (*user.AddRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.AddRoleResponse{Status: false}, err
	}
	permissionIdsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.PermissionIdList)

	if err != nil {
		return &user.AddRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !permissionIdsValid {
		return &user.AddRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidPermissionId)
	}

	err = l.svcCtx.RoleRepository.AddRole(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.RoleExists)
		} else {
			return &user.AddRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.AddRoleResponse{Status: true}, nil
}
