package repositories

import (
	"admin/user/rpc/internal/pkg/util"
	"admin/user/rpc/user"
	"context"
	"strings"
	"time"

	"admin/user/rpc/internal/Initialization"
	"admin/user/rpc/internal/model"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ApiRepository interface {
	PagingQuery(ctx context.Context, req *user.GetApiRequest) (rsp user.GetApiResponse, err error)
	FindApiByIdList(ctx context.Context, idList []string) (apiList []*user.ApiInfo, err error)
	AddApi(ctx context.Context, req *user.AddApiRequest) (err error)
	BatchAddApi(ctx context.Context, req *user.AddBatchApiRequest) (err error)
	EditApi(ctx context.Context, req *user.EditApiRequest) (err error)
	DeleteApi(ctx context.Context, req *user.DeleteApiRequest) (err error)
	VerifyApiIds(ctx context.Context, apiIds []uint64) (apiIdsValid bool, err error)
	VerifyApiUse(ctx context.Context, apiIdList []uint64) (apiUseValid bool, err error)
}

// NewApiRepository 初始化
func NewApiRepository(ds *Initialization.DataSources) ApiRepository {
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
func (a *apiRepository) PagingQuery(ctx context.Context, req *user.GetApiRequest) (rsp user.GetApiResponse, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Model(model.AdminAPI{}).Select(
		model.AdminAPIColumns.ID,
		model.AdminAPIColumns.Name,
		model.AdminAPIColumns.HTTPMethod,
		model.AdminAPIColumns.HTTPPath).
		Where(model.AdminAPIColumns.PlatformID+" = ?", platformId)
	db.Count(&rsp.Total) // 统计条数
	err = db.Scopes(PageDefault(int(req.CurrentPage), int(req.PageSize))).Find(&rsp.RowList).Error

	return
}

// FindApiByIdList 通过idList 查询api信息
func (a *apiRepository) FindApiByIdList(ctx context.Context, idList []string) (apiList []*user.ApiInfo, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model.AdminAPI{}).Select(model.AdminAPIColumns.ID, model.AdminAPIColumns.Name,
		model.AdminAPIColumns.HTTPPath, model.AdminAPIColumns.HTTPMethod).
		Where("id in ?", idList).Limit(DefaultLimit).
		Find(&apiList).Error

	return
}

// AddApi 单条插入
func (a *apiRepository) AddApi(ctx context.Context, req *user.AddApiRequest) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	err = db.Create(&model.AdminAPI{
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
func (a *apiRepository) BatchAddApi(ctx context.Context, req *user.AddBatchApiRequest) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		total = 500            // 每次插入的条数
		data  []model.AdminAPI // 插入的数据
	)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// 拼装数据
	for _, v := range req.ItemList {
		data = append(data, model.AdminAPI{
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
func (a *apiRepository) EditApi(ctx context.Context, req *user.EditApiRequest) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model.AdminAPI{}).Omit(
		model.AdminAPIColumns.ID,
		model.AdminAPIColumns.CreatedAt,
		model.AdminAPIColumns.PlatformID,
	).Where("id = ?", req.Id).Updates(model.AdminAPI{
		Name:       req.Name,
		HTTPMethod: strings.ToUpper(req.HttpMethod),
		HTTPPath:   req.HttpPath,
		UpdatedAt:  time.Now(),
	}).Error

	return
}

// DeleteApi 删除 支持批量
func (a *apiRepository) DeleteApi(ctx context.Context, req *user.DeleteApiRequest) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Delete(&model.AdminAPI{}, req.IdList).Limit(len(req.IdList)).Error

	return
}

// VerifyApiIds 验证列表中id的是否全部有效
func (a *apiRepository) VerifyApiIds(ctx context.Context, apiIdList []uint64) (apiIdsValid bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model.AdminAPI{}).Where(model.AdminRoleColumns.ID+" in ?", apiIdList).
		Where(model.AdminRoleColumns.PlatformID+"=?", platformId).Count(&total).Error

	if err != nil || total != int64(len(apiIdList)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}

// VerifyApiUse 验证api有没有被使用
func (a *apiRepository) VerifyApiUse(ctx context.Context, apiIdList []uint64) (apiUseValid bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var total int64
	err = db.Model(model.AdminPermissionHasAPI{}).
		Where(model.AdminPermissionHasAPIColumns.APIID+" in ?", apiIdList).
		Count(&total).Error

	if err != nil {
		return apiUseValid, err
	}

	if total > 0 { // api 正在使用中
		apiUseValid = true
	}

	return apiUseValid, nil
}
