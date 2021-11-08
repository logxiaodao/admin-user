package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"strings"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type FindPermissionByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindPermissionByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPermissionByIdsLogic {
	return &FindPermissionByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindPermissionByIdsLogic) FindPermissionByIds(in *user.FindPermissionByIdsRequest) (*user.FindPermissionByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.FindPermissionByIdsResponse{DataList: nil}, err
	}
	idList := strings.Split(in.Ids, ",")
	permissionList, err := l.svcCtx.PermissionRepository.FindPermissionByIdList(ctx, idList)
	if err != nil {
		return &user.FindPermissionByIdsResponse{DataList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.FindPermissionByIdsResponse{DataList: permissionList}, nil
}
