package logic

import (
	"admin/user/common/errorx"
	"admin/user/rpc/internal/pkg/auth"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"

	"admin/user/rpc/internal/svc"
	"admin/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminLogic {
	return &AddAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAdminLogic) AddAdmin(in *user.AddAdminRequest) (*user.AddAdminResponse, error) {

	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.AddAdminResponse{Status: false}, err
	}

	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.RoleIdList)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.AdminExists)
		} else {
			return &user.AddAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	if !roleIdsValid {
		return &user.AddAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidRoleId)
	}

	err = l.svcCtx.AdminRepository.AddAdmin(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.AddAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.AdminExists)
		} else {
			return &user.AddAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.AddAdminResponse{Status: true}, nil
}
