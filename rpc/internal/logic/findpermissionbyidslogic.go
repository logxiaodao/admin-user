package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type FindPermissionByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewFindPermissionByIdsLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *FindPermissionByIdsLogic {
	return &FindPermissionByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindPermissionByIdsLogic) FindPermissionByIds(in *user2.FindPermissionByIdsRequest) (*user2.FindPermissionByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.FindPermissionByIdsResponse{DataList: nil}, err
	}
	idList := strings.Split(in.Ids, ",")
	permissionList, err := l.svcCtx.PermissionRepository.FindPermissionByIdList(ctx, idList)
	if err != nil {
		return &user2.FindPermissionByIdsResponse{DataList: nil}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.FindPermissionByIdsResponse{DataList: permissionList}, nil
}
