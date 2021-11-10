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

type FindAdminByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewFindAdminByIdsLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *FindAdminByIdsLogic {
	return &FindAdminByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAdminByIdsLogic) FindAdminByIds(in *user2.FindAdminByIdsRequest) (*user2.FindAdminByIdsResponse, error) {

	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.FindAdminByIdsResponse{DataList: nil}, err
	}

	idList := strings.Split(in.Ids, ",")
	roleList, err := l.svcCtx.AdminRepository.FindAdminByIdList(ctx, idList)
	if err != nil {
		return &user2.FindAdminByIdsResponse{DataList: nil}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.FindAdminByIdsResponse{DataList: roleList}, nil
}
