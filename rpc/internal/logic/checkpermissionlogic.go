package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	"admin-user/rpc/internal/config"
	auth2 "admin-user/rpc/internal/pkg/auth"
	util2 "admin-user/rpc/internal/pkg/util"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckPermissionLogic) CheckPermission(in *user2.CheckPermissionRequest) (*user2.CheckPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, false)
	if err != nil {
		return &user2.CheckPermissionResponse{Status: false}, err
	}

	// 获取userid
	userId := util2.InterfaceToInt64(ctx.Value("userId"))
	if userId == 0 {
		return &user2.CheckPermissionResponse{Status: false}, errorx2.GetErrorByCode(errorx2.IncorrectToken)
	}

	// 判读是否公共接口
	for _, v := range config.SecurityApiData {
		if v.HTTPPath == in.HttpPath && v.HTTPMethod == in.HttpMethod {
			return &user2.CheckPermissionResponse{Status: true}, nil
		}
	}

	isPass, err := l.svcCtx.AccountRepository.CheckPermission(ctx, userId, in.HttpMethod, in.HttpPath)
	if err != nil {
		return &user2.CheckPermissionResponse{Status: isPass}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.CheckPermissionResponse{Status: isPass}, nil
}
