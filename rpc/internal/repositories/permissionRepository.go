package repositories

import (
	Initialization2 "admin-user/rpc/internal/Initialization"
	model2 "admin-user/rpc/internal/model"
	util2 "admin-user/rpc/internal/pkg/util"
	user2 "admin-user/rpc/user"
	"context"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	PagingQuery(ctx context.Context, req *user2.GetPermissionRequest) (rsp user2.GetPermissionResponse, err error)
	FindPermissionByIdList(ctx context.Context, idList []string) (permissionList []*user2.PermissionInfo, err error)
	AddPermission(ctx context.Context, req *user2.AddPermissionRequest) (err error)
	EditPermission(ctx context.Context, req *user2.EditPermissionRequest) (err error)
	DeletePermission(ctx context.Context, req *user2.DeletePermissionRequest) (err error)
	VerifyPermissionIds(ctx context.Context, PermissionIds []uint64) (permissionIdsValid bool, err error)
	VerifyPermissionUse(ctx context.Context, PermissionIds []uint64) (permissionUseValid bool, err error)
}

// NewPermissionRepository 初始化
func NewPermissionRepository(ds *Initialization2.DataSources) PermissionRepository {
	return &permissionRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type permissionRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// PagingQuery 分页查询
func (p *permissionRepository) PagingQuery(ctx context.Context, req *user2.GetPermissionRequest) (rsp user2.GetPermissionResponse, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminAPI              model2.AdminAPI
		adminPermission       model2.AdminPermission
		adminPermissionHasAPI model2.AdminPermissionHasAPI
	)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Table(adminPermission.TableName()).
		Where(adminPermission.TableName()+"."+model2.AdminRoleColumns.PlatformID+" = ?", platformId)

	if len(req.Keyword) > 0 {
		db = db.Where(model2.AdminPermissionColumns.PermissionName + " like '" + req.Keyword + "%'")
	}

	err = db.Count(&rsp.Total).Error
	if err != nil {
		return
	}

	err = db.Select(
		"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.ID+") "+model2.AdminPermissionColumns.ID,
		"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.PlatformID+") "+model2.AdminPermissionColumns.PlatformID,
		"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.PermissionName+") "+model2.AdminPermissionColumns.PermissionName,
		"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.CreatedAt+") as createdat",
		"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.UpdatedAt+") as updatedat",
		"group_concat("+model2.AdminPermissionHasAPIColumns.APIID+") api_ids",
		"group_concat("+model2.AdminAPIColumns.Name+") api_names").
		Joins("left join `" + adminPermissionHasAPI.TableName() + "` on " +
			adminPermission.TableName() + "." + model2.AdminPermissionColumns.ID + " = " +
			adminPermissionHasAPI.TableName() + "." + model2.AdminPermissionHasAPIColumns.PermissionID).
		Joins("left join `" + adminAPI.TableName() + "` on " +
			adminAPI.TableName() + "." + model2.AdminAPIColumns.ID + " = " +
			adminPermissionHasAPI.TableName() + "." + model2.AdminPermissionHasAPIColumns.APIID).
		Group(adminPermission.TableName() + "." + model2.AdminPermissionColumns.ID).
		Scopes(PageDefault(int(req.CurrentPage), int(req.PageSize))).Find(&rsp.RowList).Error

	return
}

// FindPermissionByIdList 通过idList 查询权限组信息
func (p *permissionRepository) FindPermissionByIdList(ctx context.Context, idList []string) (permissionList []*user2.PermissionInfo, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminAPI              model2.AdminAPI
		adminPermission       model2.AdminPermission
		adminPermissionHasAPI model2.AdminPermissionHasAPI
	)

	// sql 逻辑
	err = db.Table(adminPermission.TableName()).
		Select(
			"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.ID+") "+model2.AdminPermissionColumns.ID,
			"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.PlatformID+") "+model2.AdminPermissionColumns.PlatformID,
			"any_value("+adminPermission.TableName()+"."+model2.AdminPermissionColumns.PermissionName+") "+model2.AdminPermissionColumns.PermissionName,
			"group_concat("+model2.AdminPermissionHasAPIColumns.APIID+") api_ids",
			"group_concat("+model2.AdminAPIColumns.Name+") api_names").
		Joins("left join `"+adminPermissionHasAPI.TableName()+"` on "+
			adminPermission.TableName()+"."+model2.AdminPermissionColumns.ID+" = "+
			adminPermissionHasAPI.TableName()+"."+model2.AdminPermissionHasAPIColumns.PermissionID).
		Joins("left join `"+adminAPI.TableName()+"` on "+
			adminAPI.TableName()+"."+model2.AdminAPIColumns.ID+" = "+
			adminPermissionHasAPI.TableName()+"."+model2.AdminPermissionHasAPIColumns.APIID).
		Where(adminPermission.TableName()+"."+model2.AdminPermissionColumns.ID+" in ?", idList).
		Group(adminPermission.TableName() + "." + model2.AdminPermissionColumns.ID).
		Find(&permissionList).Error

	return
}

// AddPermission 单条插入
func (p *permissionRepository) AddPermission(ctx context.Context, req *user2.AddPermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util2.InterfaceToUint(ctx.Value("userId"))
	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// 插入平台
	var (
		adminPermissionHasApis []model2.AdminPermissionHasAPI
		adminPermission        = model2.AdminPermission{
			PermissionName: req.PermissionName,
			PlatformID:     platformId,
			Creater:        userId,
		}
	)

	// sql 逻辑   插入权限表
	err = tx.Create(&adminPermission).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 组装api数据
	for _, v := range req.ApiIdList {
		adminPermissionHasApis = append(adminPermissionHasApis, model2.AdminPermissionHasAPI{
			PermissionID: adminPermission.ID,
			APIID:        uint(v),
		})
	}

	// 插入权限api中间表
	err = tx.Create(adminPermissionHasApis).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// EditPermission 单条修改
func (p *permissionRepository) EditPermission(ctx context.Context, req *user2.EditPermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util2.InterfaceToUint(ctx.Value("userId"))

	// sql 逻辑
	err = tx.Model(model2.AdminPermission{}).Omit(
		model2.AdminPermissionColumns.ID,
		model2.AdminPermissionColumns.PlatformID,
		model2.AdminPermissionColumns.CreatedAt,
		model2.AdminPermissionColumns.Creater,
	).Where("id = ?", req.Id).Updates(model2.AdminPermission{
		PermissionName: req.PermissionName,
		UpdatedAt:      time.Now(),
		Updater:        userId,
	}).Error
	if err != nil {
		return err
	}

	// 删除该权限之前的api
	err = tx.Where(model2.AdminPermissionHasAPIColumns.PermissionID+"=?", req.Id).Delete(model2.AdminPermissionHasAPI{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var adminPermissionHasApis []model2.AdminPermissionHasAPI
	for _, v := range req.ApiIdList {
		adminPermissionHasApis = append(adminPermissionHasApis, model2.AdminPermissionHasAPI{PermissionID: uint(req.Id), APIID: uint(v)})
	}

	// 给权限绑定新的api列表
	err = tx.Create(adminPermissionHasApis).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return
}

// DeletePermission 删除 支持批量
func (p *permissionRepository) DeletePermission(ctx context.Context, req *user2.DeletePermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// sql 逻辑
	err = tx.Where(model2.AdminPermissionColumns.ID+" in ?", req.IdList).Delete(&model2.AdminPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(model2.AdminPermissionHasAPIColumns.PermissionID+" in ?", req.IdList).Delete(&model2.AdminPermissionHasAPI{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// VerifyPermissionIds 验证列表中id的是否全部有效
func (p *permissionRepository) VerifyPermissionIds(ctx context.Context, PermissionIds []uint64) (permissionIdsValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model2.AdminPermission{}).
		Where(model2.AdminPermissionColumns.ID+" in ?", PermissionIds).
		Where(model2.AdminPermissionColumns.PlatformID+"=?", platformId).
		Count(&total).Error

	if err != nil || total != int64(len(PermissionIds)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}

// VerifyPermissionUse 验证列表中权限是否被角色使用
func (p *permissionRepository) VerifyPermissionUse(ctx context.Context, PermissionIds []uint64) (permissionUseValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var total int64
	err = db.Model(model2.AdminRoleHasPermission{}).Where(model2.AdminRoleHasPermissionColumns.PermissionID+" in ?", PermissionIds).Count(&total).Error

	if err != nil {
		return
	}

	if total > 0 {
		permissionUseValid = true
	}

	return
}
