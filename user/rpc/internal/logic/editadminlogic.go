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

type EditAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditAdminLogic {
	return &EditAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditAdminLogic) EditAdmin(in *user.EditAdminRequest) (*user.EditAdminResponse, error) {
	// 把token内容信息添加到context
	token := auth.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user.EditAdminResponse{Status: false}, err
	}

	adminIdsValid, err := l.svcCtx.AdminRepository.VerifyAdminIds(ctx, []uint64{in.Id})
	if err != nil {
		return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !adminIdsValid {
		return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidAdminId)
	}

	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.RoleIdList)
	if err != nil {
		return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
	} else if !roleIdsValid {
		return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.InvalidRoleId)
	}

	err = l.svcCtx.AdminRepository.EditAdmin(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.AdminExists)
		} else {
			return &user.EditAdminResponse{Status: false}, errorx.GetErrorByCode(errorx.DatabaseError)
		}
	}

	return &user.EditAdminResponse{Status: true}, nil
}
