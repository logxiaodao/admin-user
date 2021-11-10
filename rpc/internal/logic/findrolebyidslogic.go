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

type FindRoleByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewFindRoleByIdsLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *FindRoleByIdsLogic {
	return &FindRoleByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindRoleByIdsLogic) FindRoleByIds(in *user2.FindRoleByIdsRequest) (*user2.FindRoleByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.FindRoleByIdsResponse{DataList: nil}, err
	}
	idList := strings.Split(in.Ids, ",")
	permissionList, err := l.svcCtx.RoleRepository.FindRoleByIdList(ctx, idList)
	if err != nil {
		return &user2.FindRoleByIdsResponse{DataList: nil}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.FindRoleByIdsResponse{DataList: permissionList}, nil
}
