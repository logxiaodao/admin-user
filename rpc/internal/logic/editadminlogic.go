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

type EditAdminLogic struct {
	ctx    context.Context
	svcCtx *svc2.ServiceContext
	logx.Logger
}

func NewEditAdminLogic(ctx context.Context, svcCtx *svc2.ServiceContext) *EditAdminLogic {
	return &EditAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditAdminLogic) EditAdmin(in *user2.EditAdminRequest) (*user2.EditAdminResponse, error) {
	// 把token内容信息添加到context
	token := auth2.NewToken(l.svcCtx.Config.AuthConf.AccessSecret)
	ctx, err := token.SetContext(l.ctx, true)
	if err != nil {
		return &user2.EditAdminResponse{Status: false}, err
	}

	adminIdsValid, err := l.svcCtx.AdminRepository.VerifyAdminIds(ctx, []uint64{in.Id})
	if err != nil {
		return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !adminIdsValid {
		return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidAdminId)
	}

	roleIdsValid, err := l.svcCtx.RoleRepository.VerifyRoleIds(ctx, in.RoleIdList)
	if err != nil {
		return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
	} else if !roleIdsValid {
		return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.InvalidRoleId)
	}

	err = l.svcCtx.AdminRepository.EditAdmin(ctx, in)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) && sqlErr.Number == 1062 { // mysql触发唯一索引处理
			return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.AdminExists)
		} else {
			return &user2.EditAdminResponse{Status: false}, errorx2.GetErrorByCode(errorx2.DatabaseError)
		}
	}

	return &user2.EditAdminResponse{Status: true}, nil
}
