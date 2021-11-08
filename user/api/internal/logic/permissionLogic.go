package logic

import (
	"context"
	"strings"

	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) PermissionLogic {
	return PermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionLogic) GetPermission(req types.GetPermissionReq) (*types.GetPermissionRsp, error) {

	result, err := l.svcCtx.PermissionRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *PermissionLogic) DeletePermission(req types.DeletePermissionReq) (idIsValid, permissionUseValid bool, err error) {

	// 验证ids的有效性
	idIsValid, err = l.svcCtx.PermissionRepository.VerifyPermissionIds(l.ctx, req.IdList)
	if err != nil || !idIsValid {
		return
	}

	// 验证权限是否在被角色使用，如果被使用不允许删除
	permissionUseValid, err = l.svcCtx.PermissionRepository.VerifyPermissionUse(l.ctx, req.IdList)
	if err != nil || permissionUseValid {
		return
	}

	err = l.svcCtx.PermissionRepository.DeletePermission(l.ctx, req)

	return
}

func (l *PermissionLogic) PutPermission(req types.PutPermissionReq) (apiIdsValid, permissionIdsValid bool, err error) {

	permissionIdsValid, err = l.svcCtx.PermissionRepository.VerifyPermissionIds(l.ctx, []uint{req.Id})
	if err != nil || !permissionIdsValid {
		return
	}

	// 验证 apiIds 的有效性
	apiIdsValid, err = l.svcCtx.ApiRepository.VerifyApiIds(l.ctx, req.ApiIdList)
	if err != nil || !apiIdsValid {
		return
	}

	err = l.svcCtx.PermissionRepository.EditPermission(l.ctx, req)

	return
}

func (l *PermissionLogic) PostPermission(req types.PostPermissionReq) (apiIdsValid bool, err error) {

	// 验证 apiIds 的有效性
	apiIdsValid, err = l.svcCtx.ApiRepository.VerifyApiIds(l.ctx, req.ApiIdList)
	if err != nil || !apiIdsValid {
		return
	}

	err = l.svcCtx.PermissionRepository.AddPermission(l.ctx, req)

	return
}

func (l *PermissionLogic) FindPermissionByIds(req types.FindPermissionByIdsReq) (*types.FindPermissionByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	result, err := l.svcCtx.PermissionRepository.FindPermissionByIdList(l.ctx, idList)
	if err != nil {
		return nil, err
	}

	return &types.FindPermissionByIdsRsp{DataList: result}, nil
}
