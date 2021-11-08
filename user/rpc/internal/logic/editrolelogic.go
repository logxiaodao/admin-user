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

type EditRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditRoleLogic {
	return &EditRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditRoleLogic) EditRole(in *user.EditRoleRequest) (*user.EditRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.EditRoleResponse{Status: false}, err
	}
	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, []uint64{in.Id})
	if err != nil {
		return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !roleIdsValid {
		return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidRoleId)
	}

	permissionIdsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.PermissionIdList)
	if err != nil {
		return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !permissionIdsValid {
		return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidPermissionId)
	}

	err = l.svcCtx.RoleRepository.EditRole(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.RoleExists)
		} else {
			return &user.EditRoleResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.EditRoleResponse{Status: true}, nil
}
