package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type EditApiLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewEditApiLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *EditApiLogic {
	return &EditApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditApiLogic) EditApi(in *user2.EditApiRequest) (*user2.EditApiResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.EditApiResponse{Status: false}, err
	}

	// 验证id有效
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, []uint64{in.Id})
	if err != nil || !apiIdsValid {
		return &user2.EditApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidId)
	}

	err = l.svcCtx.ApiRepository.EditApi(ctx, in)
	if err != nil {
		return &user2.EditApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.EditApiResponse{Status: true}, nil
}
