package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteApiLogic) DeleteApi(in *user.DeleteApiRequest) (*user.DeleteApiResponse, error) {

	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.DeleteApiResponse{Status: false}, err
	}

	// 验证id有效
	apiIdsValid, err := l.svcCtx.ApiRepository.VerifyApiIds(ctx, in.IdList)
	if err != nil {
		return &user.DeleteApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !apiIdsValid {
		return &user.DeleteApiResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidId)
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	apiUseValid, err := l.svcCtx.ApiRepository.VerifyApiUse(ctx, in.IdList)
	if err != nil {
		return &user.DeleteApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if apiUseValid {
		return &user.DeleteApiResponse{Status: false}, errorx.GetErrorByCode(errorx.ApiInUse)
	}

	err = l.svcCtx.ApiRepository.DeleteApi(ctx, in)
	if err != nil {
		return &user.DeleteApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	}

	return &user.DeleteApiResponse{Status: true}, nil
}
