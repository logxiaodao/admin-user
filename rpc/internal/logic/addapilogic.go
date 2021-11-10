package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddApiLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewAddApiLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *AddApiLogic {
	return &AddApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddApiLogic) AddApi(in *user2.AddApiRequest) (*user2.AddApiResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.AddApiResponse{Status: false}, err
	}

	err = l.svcCtx.ApiRepository.AddApi(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.AddApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.ApiExists)
		} else {
			return &user2.AddApiResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.AddApiResponse{Status: true}, nil
}
