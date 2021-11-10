package repositories

import (
	Initialization2 "admin-user/api/internal/Initialization"
	model2 "admin-user/api/internal/model"
	util2 "admin-user/api/internal/pkg/util"
	types2 "admin-user/api/internal/types"
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ApiRepository interface {
	PagingQuery(ctx context.Context, req types2.GetApiReq) (rsp types2.GetApiRsp, err error)
	FindApiByIdList(ctx context.Context, idList []string) (apiList []types2.ApiList, err error)
	AddApi(ctx context.Context, req types2.PostApiReq) (err error)
	BatchAddApi(ctx context.Context, req types2.BatchPostApiReq) (err error)
	EditApi(ctx context.Context, req types2.PutApiReq) (err error)
	DeleteApi(ctx context.Context, req types2.DeleteApiReq) (err error)
	VerifyApiIds(ctx context.Context, apiIds []uint) (apiIdsValid bool, err error)
	VerifyApiUse(ctx context.Context, apiIdList []uint) (apiUseValid bool, err error)
}

// NewApiRepository 初始化
func NewApiRepository(ds *Initialization2.DataSources) ApiRepository {
	return &apiRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type apiRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// PagingQuery 分页查询
func (a *apiRepository) PagingQuery(ctx context.Context, req types2.GetApiReq) (rsp types2.GetApiRsp, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Model(model2.AdminAPI{}).Select(
		model2.AdminAPIColumns.ID,
		model2.AdminAPIColumns.Name,
		model2.AdminAPIColumns.HTTPMethod,
		model2.AdminAPIColumns.HTTPPath).
		Where(model2.AdminAPIColumns.PlatformID+" = ?", platformId)
	db.Count(&rsp.Total) // 统计条数
	err = db.Scopes(PageDefault(req.CurrentPage, req.PageSize)).Find(&rsp.RowList).Error

	return
}

// FindApiByIdList 通过idList 查询api信息
func (a *apiRepository) FindApiByIdList(ctx context.Context, idList []string) (apiList []types2.ApiList, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model2.AdminAPI{}).Select(model2.AdminAPIColumns.ID, model2.AdminAPIColumns.Name,
		model2.AdminAPIColumns.HTTPPath, model2.AdminAPIColumns.HTTPMethod).
		Where("id in ?", idList).Limit(DefaultLimit).
		Find(&apiList).Error

	return
}

// AddApi 单条插入
func (a *apiRepository) AddApi(ctx context.Context, req types2.PostApiReq) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	err = db.Create(&model2.AdminAPI{
		Name:       req.Name,
		HTTPMethod: strings.ToUpper(req.HttpMethod),
		HTTPPath:   req.HttpPath,
		PlatformID: platformId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}).Error

	return
}

// BatchAddApi 批量插入
func (a *apiRepository) BatchAddApi(ctx context.Context, req types2.BatchPostApiReq) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		total = 500             // 每次插入的条数
		data  []model2.AdminAPI // 插入的数据
	)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// 拼装数据
	for _, v := range req.ItemList {
		data = append(data, model2.AdminAPI{
			Name:       v.Name,
			HTTPMethod: strings.ToUpper(v.HttpMethod),
			HTTPPath:   v.HttpPath,
			PlatformID: platformId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	// sql 逻辑
	err = db.CreateInBatches(data, total).Error

	return
}

// EditApi 单条修改
func (a *apiRepository) EditApi(ctx context.Context, req types2.PutApiReq) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model2.AdminAPI{}).Omit(
		model2.AdminAPIColumns.ID,
		model2.AdminAPIColumns.CreatedAt,
		model2.AdminAPIColumns.PlatformID,
	).Where("id = ?", req.Id).Updates(model2.AdminAPI{
		Name:       req.Name,
		HTTPMethod: strings.ToUpper(req.HttpMethod),
		HTTPPath:   req.HttpPath,
		UpdatedAt:  time.Now(),
	}).Error

	return
}

// DeleteApi 删除 支持批量
func (a *apiRepository) DeleteApi(ctx context.Context, req types2.DeleteApiReq) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Delete(&model2.AdminAPI{}, req.IdList).Limit(len(req.IdList)).Error

	return
}

// VerifyApiIds 验证列表中id的是否全部有效
func (a *apiRepository) VerifyApiIds(ctx context.Context, apiIdList []uint) (apiIdsValid bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model2.AdminAPI{}).Where(model2.AdminRoleColumns.ID+" in ?", apiIdList).
		Where(model2.AdminRoleColumns.PlatformID+"=?", platformId).Count(&total).Error

	if err != nil || total != int64(len(apiIdList)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}

// VerifyApiUse 验证api有没有被使用
func (a *apiRepository) VerifyApiUse(ctx context.Context, apiIdList []uint) (apiUseValid bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var total int64
	err = db.Model(model2.AdminPermissionHasAPI{}).
		Where(model2.AdminPermissionHasAPIColumns.APIID+" in ?", apiIdList).
		Count(&total).Error

	if err != nil {
		return apiUseValid, err
	}

	if total > 0 { // api 正在使用中
		apiUseValid = true
	}

	return apiUseValid, nil
}
