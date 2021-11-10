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

type EditRoleLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewEditRoleLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *EditRoleLogic {
	return &EditRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditRoleLogic) EditRole(in *user2.EditRoleRequest) (*user2.EditRoleResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.EditRoleResponse{Status: false}, err
	}
	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, []uint64{in.Id})
	if err != nil {
		return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !roleIdsValid {
		return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidRoleId)
	}

	permissionIdsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, in.PermissionIdList)
	if err != nil {
		return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !permissionIdsValid {
		return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidPermissionId)
	}

	err = l.svcCtx.RoleRepository.EditRole(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.RoleExists)
		} else {
			return &user2.EditRoleResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.EditRoleResponse{Status: true}, nil
}
