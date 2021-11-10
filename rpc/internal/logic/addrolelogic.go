package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddRoleLogic) AddRole(in *user2.AddRoleRequest) (*user2.AddRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.AddRoleResponse{Status: false}, err
	}
	permissionIdsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.PermissionIdList)

	if err != nil {
		return &user2.AddRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !permissionIdsValid {
		return &user2.AddRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidPermissionId)
	}

	err = l.svcCtx.RoleRepository.AddRole(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.AddRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.RoleExists)
		} else {
			return &user2.AddRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.AddRoleResponse{Status: true}, nil
}
