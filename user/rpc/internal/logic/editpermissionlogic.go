package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type EditPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditPermissionLogic {
	return &EditPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditPermissionLogic) EditPermission(in *user.EditPermissionRequest) (*user.EditPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.EditPermissionResponse{Status: false}, err
	}

	permissionIdsValid, err := l.svcCtx.PermissionRepository.VerifyPermissionIds(ctx, []uint64{in.Id})
	if err != nil {
		return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !permissionIdsValid {
		return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidId)
	}

	// 验证 apiIds 的有效性
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, in.ApiIdList)
	if err != nil {
		return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !apiIdsValid {
		return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidApiId)
	}

	err = l.svcCtx.PermissionRepository.EditPermission(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.PermissionExists)
		} else {
			return &user.EditPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.EditPermissionResponse{Status: true}, nil
}
