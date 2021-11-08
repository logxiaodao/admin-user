package repositories

import (
	"admin/user/rpc/internal/pkg/util"
	"context"
	"time"

	"admin/user/rpc/internal/Initialization"
	"admin/user/rpc/internal/model"
	"admin/user/rpc/user"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	PagingQuery(ctx context.Context, req *user.GetPermissionRequest) (rsp user.GetPermissionResponse, err error)
	FindPermissionByIdList(ctx context.Context, idList []string) (permissionList []*user.PermissionInfo, err error)
	AddPermission(ctx context.Context, req *user.AddPermissionRequest) (err error)
	EditPermission(ctx context.Context, req *user.EditPermissionRequest) (err error)
	DeletePermission(ctx context.Context, req *user.DeletePermissionRequest) (err error)
	VerifyPermissionIds(ctx context.Context, PermissionIds []uint64) (permissionIdsValid bool, err error)
	VerifyPermissionUse(ctx context.Context, PermissionIds []uint64) (permissionUseValid bool, err error)
}

// NewPermissionRepository 初始化
func NewPermissionRepository(ds *Initialization.DataSources) PermissionRepository {
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
func (p *permissionRepository) PagingQuery(ctx context.Context, req *user.GetPermissionRequest) (rsp user.GetPermissionResponse, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminPermission       model.AdminPermission
		adminPermissionHasAPI model.AdminPermissionHasAPI
	)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Table(adminPermission.TableName()).
		Where(adminPermission.TableName()+"."+model.AdminRoleColumns.PlatformID+" = ?", platformId)

	err = db.Count(&rsp.Total).Error
	if err != nil {
		return
	}

	err = db.Select(
		"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.ID+") "+model.AdminPermissionColumns.ID,
		"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.PlatformID+") "+model.AdminPermissionColumns.PlatformID,
		"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.PermissionName+") "+model.AdminPermissionColumns.PermissionName,
		"group_concat("+model.AdminPermissionHasAPIColumns.APIID+") api_ids").
		Joins("left join `" + adminPermissionHasAPI.TableName() + "` on " +
			adminPermission.TableName() + "." + model.AdminPermissionColumns.ID + " = " +
			adminPermissionHasAPI.TableName() + "." + model.AdminPermissionHasAPIColumns.PermissionID).
		Group(adminPermission.TableName() + "." + model.AdminPermissionColumns.ID).
		Scopes(PageDefault(int(req.CurrentPage), int(req.PageSize))).Find(&rsp.RowList).Error

	return
}

// FindPermissionByIdList 通过idList 查询权限组信息
func (p *permissionRepository) FindPermissionByIdList(ctx context.Context, idList []string) (permissionList []*user.PermissionInfo, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminPermission       model.AdminPermission
		adminPermissionHasAPI model.AdminPermissionHasAPI
	)

	// sql 逻辑
	err = db.Table(adminPermission.TableName()).
		Select(
			"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.ID+") "+model.AdminPermissionColumns.ID,
			"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.PlatformID+") "+model.AdminPermissionColumns.PlatformID,
			"any_value("+adminPermission.TableName()+"."+model.AdminPermissionColumns.PermissionName+") "+model.AdminPermissionColumns.PermissionName,
			"group_concat("+model.AdminPermissionHasAPIColumns.APIID+") api_ids").
		Joins("left join `"+adminPermissionHasAPI.TableName()+"` on "+
			adminPermission.TableName()+"."+model.AdminPermissionColumns.ID+" = "+
			adminPermissionHasAPI.TableName()+"."+model.AdminPermissionHasAPIColumns.PermissionID).
		Where(adminPermission.TableName()+"."+model.AdminPermissionColumns.ID+" in ?", idList).
		Group(adminPermission.TableName() + "." + model.AdminPermissionColumns.ID).
		Find(&permissionList).Error

	return
}

// AddPermission 单条插入
func (p *permissionRepository) AddPermission(ctx context.Context, req *user.AddPermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util.InterfaceToUint(ctx.Value("userId"))
	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// 插入平台
	var (
		adminPermissionHasApis []model.AdminPermissionHasAPI
		adminPermission        = model.AdminPermission{
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
		adminPermissionHasApis = append(adminPermissionHasApis, model.AdminPermissionHasAPI{
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
func (p *permissionRepository) EditPermission(ctx context.Context, req *user.EditPermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util.InterfaceToUint(ctx.Value("userId"))

	// sql 逻辑
	err = tx.Model(model.AdminPermission{}).Omit(
		model.AdminPermissionColumns.ID,
		model.AdminPermissionColumns.PlatformID,
		model.AdminPermissionColumns.CreatedAt,
		model.AdminPermissionColumns.Creater,
	).Where("id = ?", req.Id).Updates(model.AdminPermission{
		PermissionName: req.PermissionName,
		UpdatedAt:      time.Now(),
		Updater:        userId,
	}).Error
	if err != nil {
		return err
	}

	// 删除该权限之前的api
	err = tx.Where(model.AdminPermissionHasAPIColumns.PermissionID+"=?", req.Id).Delete(model.AdminPermissionHasAPI{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var adminPermissionHasApis []model.AdminPermissionHasAPI
	for _, v := range req.ApiIdList {
		adminPermissionHasApis = append(adminPermissionHasApis, model.AdminPermissionHasAPI{PermissionID: uint(req.Id), APIID: uint(v)})
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
func (p *permissionRepository) DeletePermission(ctx context.Context, req *user.DeletePermissionRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// sql 逻辑
	err = tx.Where(model.AdminPermissionColumns.ID+" in ?", req.IdList).Delete(&model.AdminPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(model.AdminPermissionHasAPIColumns.PermissionID+" in ?", req.IdList).Delete(&model.AdminPermissionHasAPI{}).Error
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
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model.AdminPermission{}).
		Where(model.AdminPermissionColumns.ID+" in ?", PermissionIds).
		Where(model.AdminPermissionColumns.PlatformID+"=?", platformId).
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
	err = db.Model(model.AdminRoleHasPermission{}).Where(model.AdminRoleHasPermissionColumns.PermissionID+" in ?", PermissionIds).Count(&total).Error

	if err != nil {
		return
	}

	if total > 0 {
		permissionUseValid = true
	}

	return
}
