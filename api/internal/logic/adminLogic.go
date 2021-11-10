package logic

import (
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"context"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc2.ServiceContext
}

func NewAdminLogic(ctx context.Context, svcCtx *svc2.ServiceContext) AdminLogic {
	return AdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogic) DeleteAdmin(req types2.DeleteAdminReq) (adminIdsValid bool, err error) {
	// 验证id有效
	adminIdsValid, err = l.svcCtx.AdminRepository.VerifyAdminIds(l.ctx, req.IdList)
	if err != nil || !adminIdsValid {
		return
	}

	err = l.svcCtx.AdminRepository.DeleteAdmin(l.ctx, req)

	return
}

func (l *AdminLogic) GetAdmin(req types2.GetAdminReq) (*types2.GetAdminRsp, error) {

	result, err := l.svcCtx.AdminRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *AdminLogic) FindAdminByIds(req types2.FindAdminByIdsReq) (*types2.FindAdminByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	result, err := l.svcCtx.AdminRepository.FindAdminByIdList(l.ctx, idList)
	if err != nil {
		return nil, err
	}

	return &types2.FindAdminByIdsRsp{DataList: result}, nil
}

func (l *AdminLogic) PostAdmin(req types2.PostAdminReq) (roleIdsValid bool, err error) {

	roleIdsValid, err = l.svcCtx.RoleRepository.VerifyRoleIds(l.ctx, req.RoleIdList)
	if err != nil || !roleIdsValid {
		return
	}

	err = l.svcCtx.AdminRepository.AddAdmin(l.ctx, req)

	return
}

func (l *AdminLogic) PutAdmin(req types2.PutAdminReq) (roleIdsValid, adminIdsValid bool, err error) {

	adminIdsValid, err = l.svcCtx.AdminRepository.VerifyAdminIds(l.ctx, []uint{req.Id})
	if err != nil || !adminIdsValid {
		return
	}

	roleIdsValid, err = l.svcCtx.RoleRepository.VerifyRoleIds(l.ctx, req.RoleIdList)
	if err != nil || !roleIdsValid {
		return
	}

	err = l.svcCtx.AdminRepository.EditAdmin(l.ctx, req)

	return
}
