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

type BatchApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchApiLogic {
	return &BatchApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchApiLogic) BatchApi(in *user.AddBatchApiRequest) (*user.AddBatchApiResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.AddBatchApiResponse{Status: false}, err
	}
	err = l.svcCtx.ApiRepository.BatchAddApi(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddBatchApiResponse{Status: false}, errorx.GetErrorByCode(errorx.ApiExists)
		} else {
			return &user.AddBatchApiResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.AddBatchApiResponse{Status: true}, nil
}
