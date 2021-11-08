package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"admin/user/rpc/internal/pkg/util"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckPermissionLogic) CheckPermission(in *user.CheckPermissionRequest) (*user.CheckPermissionResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, false)
	if err != nil {
		return &user.CheckPermissionResponse{Status: false}, err
	}
	// 获取userid
	userId := util.InterfaceToInt64(ctx.Value("userId"))
	if userId == 0 {
		return &user.CheckPermissionResponse{Status: false}, errorx.GetErrorByCode(errorx.IncorrectToken)
	}

	isPass, err := l.svcCtx.AccountRepository.CheckPermission(ctx, userId, in.HttpMethod, in.HttpPath)
	if err != nil {
		return &user.CheckPermissionResponse{Status: isPass}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.CheckPermissionResponse{Status: isPass}, nil
}
