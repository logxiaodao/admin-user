package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddApiLogic {
	return &AddApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddApiLogic) AddApi(in *user.AddApiRequest) (*user.AddApiResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.AddApiResponse{Status: false}, err
	}

	err = l.svcCtx.ApiRepository.AddApi(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddApiResponse{Status: false}, errorx.GetErrorByCode(errorx.ApiExists)
		} else {
			return &user.AddApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.AddApiResponse{Status: true}, nil
}
