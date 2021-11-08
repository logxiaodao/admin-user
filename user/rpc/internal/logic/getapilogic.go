package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiLogic {
	return &GetApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  api
func (l *GetApiLogic) GetApi(in *user.GetApiRequest) (*user.GetApiResponse, error) {

	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.GetApiResponse{Total: 0, RowList: nil}, err
	}

	data, err := l.svcCtx.ApiRepository.PagingQuery(ctx, in)
	if err != nil {
		return &user.GetApiResponse{Total: 0, RowList: nil}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &data, err
}
