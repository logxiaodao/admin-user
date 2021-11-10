package logic

import (
	errorx2 "admin-user/rpc/common/errorx"
	auth2 "admin-user/rpc/internal/pkg/auth"
	svc2 "admin-user/rpc/internal/svc"
	user2 "admin-user/rpc/user"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddAdminLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewAddAdminLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *AddAdminLogic {
	return &AddAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAdminLogic) AddAdmin(in *user2.AddAdminRequest) (*user2.AddAdminResponse, error) {

	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.AddAdminResponse{Status: false}, err
	}

	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.RoleIdList)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.AddAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.AdminExists)
		} else {
			return &user2.AddAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	if !roleIdsValid {
		return &user2.AddAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidRoleId)
	}

	err = l.svcCtx.AdminRepository.AddAdmin(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.AddAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.AdminExists)
		} else {
			return &user2.AddAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.AddAdminResponse{Status: true}, nil
}
