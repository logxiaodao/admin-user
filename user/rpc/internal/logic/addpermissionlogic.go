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

type AddPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPermissionLogic {
	return &AddPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPermissionLogic) AddPermission(in *user.AddPermissionRequest) (*user.AddPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.AddPermissionResponse{Status: false}, err
	}
	// 验证 apiIds 的有效性
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, in.ApiIdList)
	if err != nil {
		return &user.AddPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !apiIdsValid {
		return &user.AddPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidApiId)
	}

	err = l.svcCtx.PermissionRepository.AddPermission(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.PermissionExists)
		} else {
			return &user.AddPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.AddPermissionResponse{Status: true}, nil
}
