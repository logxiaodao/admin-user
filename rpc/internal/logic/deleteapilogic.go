package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteApiLogic) DeleteApi(in *user2.DeleteApiRequest) (*user2.DeleteApiResponse, error) {

	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.DeleteApiResponse{Status: false}, err
	}

	// 验证id有效
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, in.IdList)
	if err != nil {
		return &user2.DeleteApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !apiIdsValid {
		return &user2.DeleteApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidId)
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	apiUseValid, err := l.svcCtx.ApiRepository.VerifyApiUse(ctx, in.IdList)
	if err != nil {
		return &user2.DeleteApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if apiUseValid {
		return &user2.DeleteApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.ApiInUse)
	}

	err = l.svcCtx.ApiRepository.DeleteApi(ctx, in)
	if err != nil {
		return &user2.DeleteApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	}

	return &user2.DeleteApiResponse{Status: true}, nil
}
