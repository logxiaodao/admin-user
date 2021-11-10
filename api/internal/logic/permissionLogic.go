package logic

import (
	svc2 "admin-user/api/internal/svc"
	types2 "admin-user/api/internal/types"
	"context"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type PermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc2.ServiceContext
}

func NewPermissionLogic(ctx context.Context, svcCtx *svc2.ServiceContext) PermissionLogic {
	return PermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionLogic) GetPermission(req types2.GetPermissionReq) (*types2.GetPermissionRsp, error) {

	result, err := l.svcCtx.PermissionRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *PermissionLogic) DeletePermission(req types2.DeletePermissionReq) (idIsValid, permissionUseValid bool, err error) {

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

func (l *PermissionLogic) PutPermission(req types2.PutPermissionReq) (apiIdsValid, permissionIdsValid bool, err error) {

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

func (l *PermissionLogic) PostPermission(req types2.PostPermissionReq) (apiIdsValid bool, err error) {

	// 验证 apiIds 的有效性
	apiIdsValid, err = l.svcCtx.ApiRepository.VerifyApiIds(l.ctx, req.ApiIdList)
	if err != nil || !apiIdsValid {
		return
	}

	err = l.svcCtx.PermissionRepository.AddPermission(l.ctx, req)

	return
}

func (l *PermissionLogic) FindPermissionByIds(req types2.FindPermissionByIdsReq) (*types2.FindPermissionByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	result, err := l.svcCtx.PermissionRepository.FindPermissionByIdList(l.ctx, idList)
	if err != nil {
		return nil, err
	}

	return &types2.FindPermissionByIdsRsp{DataList: result}, nil
}
