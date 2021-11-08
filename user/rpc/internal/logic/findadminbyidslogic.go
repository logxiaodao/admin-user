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

type FindAdminByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAdminByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAdminByIdsLogic {
	return &FindAdminByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAdminByIdsLogic) FindAdminByIds(in *user.FindAdminByIdsRequest) (*user.FindAdminByIdsResponse, error) {

	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.FindAdminByIdsResponse{DataList: nil}, err
	}

	idList := strings.Split(in.Ids, ",")
	roleList, err := l.svcCtx.AdminRepository.FindAdminByIdList(ctx, idList)
	if err != nil {
		return &user.FindAdminByIdsResponse{DataList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.FindAdminByIdsResponse{DataList: roleList}, nil
}
