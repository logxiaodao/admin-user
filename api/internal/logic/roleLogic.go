package logic

import (
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"context"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type RoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc2.ServiceContext
}

func NewRoleLogic(ctx context.Context, svcCtx *svc2.ServiceContext) RoleLogic {
	return RoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleLogic) DeleteRole(req types2.DeleteRoleReq) (roleIdsValid, roleUseValid bool, err error) {

	// 验证id有效
	roleIdsValid, err = l.svcCtx.RoleRepository.VerifyRoleIds(l.ctx, req.IdList)
	if err != nil || !roleIdsValid {
		return
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	roleUseValid, err = l.svcCtx.RoleRepository.VerifyRoleUse(l.ctx, req.IdList)
	if err != nil || roleUseValid {
		return
	}

	err = l.svcCtx.RoleRepository.DeleteRole(l.ctx, req)

	return
}

func (l *RoleLogic) PutRole(req types2.PutRoleReq) (permissionIdsValid, roleIdsValid bool, err error) {

	roleIdsValid, err = l.svcCtx.RoleRepository.VerifyRoleIds(l.ctx, []uint{req.Id})
	if err != nil || !roleIdsValid {
		return
	}

	permissionIdsValid, err = l.svcCtx.PermissionRepository.VerifyPermissionIds(l.ctx, req.PermissionIdList)
	if err != nil || !permissionIdsValid {
		return
	}

	err = l.svcCtx.RoleRepository.EditRole(l.ctx, req)

	return
}

func (l *RoleLogic) PostRole(req types2.PostRoleReq) (permissionIdsValid bool, err error) {

	permissionIdsValid, err = l.svcCtx.PermissionRepository.VerifyPermissionIds(l.ctx, req.PermissionIdList)
	if err != nil || !permissionIdsValid {
		return
	}

	err = l.svcCtx.RoleRepository.AddRole(l.ctx, req)

	return
}

func (l *RoleLogic) GetRole(req types2.GetRoleReq) (*types2.GetRoleRsp, error) {

	result, err := l.svcCtx.RoleRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *RoleLogic) FindRoleByIds(req types2.FindRoleByIdsReq) (*types2.FindRoleByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	apiList, err := l.svcCtx.RoleRepository.FindRoleByIdList(l.ctx, idList)
	return &types2.FindRoleByIdsRsp{
		DataList: apiList,
	}, err
}
