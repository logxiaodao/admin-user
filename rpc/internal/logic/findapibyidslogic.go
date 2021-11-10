package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"errors"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type FindApiByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewFindApiByIdsLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *FindApiByIdsLogic {
	return &FindApiByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindApiByIdsLogic) FindApiByIds(in *user2.FindApiByIdsRequest) (*user2.FindApiByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.FindApiByIdsResponse{DataList: nil}, err
	}
	databaseError := errorx2.DefaultCodeMessage[errorx2.DatabaseError]
	idList := strings.Split(in.Ids, ",")
	apiList, err := l.svcCtx.ApiRepository.FindApiByIdList(ctx, idList)
	if err != nil {
		return &user2.FindApiByIdsResponse{DataList: nil}, errors.New(databaseError.Error())
	}

	return &user2.FindApiByIdsResponse{DataList: apiList}, nil
}
