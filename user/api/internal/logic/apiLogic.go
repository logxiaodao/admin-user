package logic

import (
	"context"
	"strings"

	"admin/user/api/internal/svc"
	"admin/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) ApiLogic {
	return ApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetApi 分页查询
func (l *ApiLogic) GetApi(req types.GetApiReq) (*types.GetApiRsp, error) {

	result, err := l.svcCtx.ApiRepository.PagingQuery(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *ApiLogic) FindApiByIdList(req types.FindApiByIdsReq) (*types.FindApiByIdsRsp, error) {

	idList := strings.Split(req.Ids, ",")
	apiList, err := l.svcCtx.ApiRepository.FindApiByIdList(l.ctx, idList)
	return &types.FindApiByIdsRsp{
		DataList: apiList,
	}, err
}

// PutApi 单条api修改
func (l *ApiLogic) PutApi(req types.PutApiReq) (apiIdsValid bool, err error) {
	// 验证id有效
	apiIdsValid, err = l.svcCtx.ApiRepository.VerifyApiIds(l.ctx, []uint{req.Id})
	if err != nil || !apiIdsValid {
		return
	}

	err = l.svcCtx.ApiRepository.EditApi(l.ctx, req)

	return
}

// BatchPostApiReq 修改api 支持批量
func (l *ApiLogic) BatchPostApiReq(req types.BatchPostApiReq) (err error) {

	err = l.svcCtx.ApiRepository.BatchAddApi(l.ctx, req)

	return
}

// PostApi 添加api 支持批量
func (l *ApiLogic) PostApi(req types.PostApiReq) (err error) {

	err = l.svcCtx.ApiRepository.AddApi(l.ctx, req)

	return
}

// DeleteApi api删除 支持批量
func (l *ApiLogic) DeleteApi(req types.DeleteApiReq) (apiIdsValid, apiUseValid bool, err error) {

	// 验证id有效
	apiIdsValid, err = l.svcCtx.ApiRepository.VerifyApiIds(l.ctx, req.IdList)
	if err != nil || !apiIdsValid {
		return
	}

	// 验证api是否再被权限使用，如果被使用不允许删除
	apiUseValid, err = l.svcCtx.ApiRepository.VerifyApiUse(l.ctx, req.IdList)
	if err != nil || apiUseValid {
		return
	}

	err = l.svcCtx.ApiRepository.DeleteApi(l.ctx, req)

	return
}
