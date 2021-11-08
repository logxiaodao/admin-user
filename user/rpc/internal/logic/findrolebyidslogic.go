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

type FindRoleByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindRoleByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRoleByIdsLogic {
	return &FindRoleByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindRoleByIdsLogic) FindRoleByIds(in *user.FindRoleByIdsRequest) (*user.FindRoleByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.FindRoleByIdsResponse{DataList: nil}, err
	}
	idList := strings.Split(in.Ids, ",")
	permissionList, err := l.svcCtx.RoleRepository.FindRoleByIdList(ctx, idList)
	if err != nil {
		return &user.FindRoleByIdsResponse{DataList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.FindRoleByIdsResponse{DataList: permissionList}, nil
}
