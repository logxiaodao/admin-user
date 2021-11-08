package logic

import (
	"context"
	"strings"

	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminLogic {
	return AdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogic) DeleteAdmin(req types.DeleteAdminReq) (adminIdsValid bool, err error) {
	// 验证id有效
	adminIdsValid, err = l.svcCtx.AdminRepository.VerifyAdminIds(l.ctx, req.IdList)
	if err != nil || !adminIdsValid {
		return
	}

	err = l.svcCtx.AdminRepository.DeleteAdmin(l.ctx, req)

	return
}

func (l *AdminLogic) GetAdmin(req types.GetAdminReq) (*types.GetAdminRsp, error) {

	result, err := l.svcCtx.AdminRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *AdminLogic) FindAdminByIds(req types.FindAdminByIdsReq) (*types.FindAdminByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	result, err := l.svcCtx.AdminRepository.FindAdminByIdList(l.ctx, idList)
	if err != nil {
		return nil, err
	}

	return &types.FindAdminByIdsRsp{DataList: result}, nil
}

func (l *AdminLogic) PostAdmin(req types.PostAdminReq) (roleIdsValid bool, err error) {

	roleIdsValid, err = l.svcCtx.RoleRepository.VerifyRoleIds(l.ctx, req.RoleIdList)
	if err != nil || !roleIdsValid {
		return
	}

	err = l.svcCtx.AdminRepository.AddAdmin(l.ctx, req)

	return
}

func (l *AdminLogic) PutAdmin(req types.PutAdminReq) (roleIdsValid, adminIdsValid bool, err error) {

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
