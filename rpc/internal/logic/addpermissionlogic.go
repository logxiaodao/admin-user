package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewAddPermissionLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *AddPermissionLogic {
	return &AddPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPermissionLogic) AddPermission(in *user2.AddPermissionRequest) (*user2.AddPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.AddPermissionResponse{Status: false}, err
	}
	// 验证 apiIds 的有效性
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, in.ApiIdList)
	if err != nil {
		return &user2.AddPermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !apiIdsValid {
		return &user2.AddPermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidApiId)
	}

	err = l.svcCtx.PermissionRepository.AddPermission(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.AddPermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.PermissionExists)
		} else {
			return &user2.AddPermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.AddPermissionResponse{Status: true}, nil
}
