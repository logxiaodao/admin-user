package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"errors"
	"strings"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type FindApiByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindApiByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindApiByIdsLogic {
	return &FindApiByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindApiByIdsLogic) FindApiByIds(in *user.FindApiByIdsRequest) (*user.FindApiByIdsResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.FindApiByIdsResponse{DataList: nil}, err
	}
	databaseError := errorx.DefaultCodeMessage[errorx.DatabaseError]
	idList := strings.Split(in.Ids, ",")
	apiList, err := l.svcCtx.ApiRepository.FindApiByIdList(ctx, idList)
	if err != nil {
		return &user.FindApiByIdsResponse{DataList: nil}, errors.New(databaseError.Error())
	}

	return &user.FindApiByIdsResponse{DataList: apiList}, nil
}
