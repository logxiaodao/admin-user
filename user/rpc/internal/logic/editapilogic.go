package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type EditApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditApiLogic {
	return &EditApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditApiLogic) EditApi(in *user.EditApiRequest) (*user.EditApiResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.EditApiResponse{Status: false}, err
	}

	// 验证id有效
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, []uint64{in.Id})
	if err != nil || !apiIdsValid {
		return &user.EditApiResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidId)
	}

	err = l.svcCtx.ApiRepository.EditApi(ctx, in)
	if err != nil {
		return &user.EditApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.EditApiResponse{Status: true}, nil
}
