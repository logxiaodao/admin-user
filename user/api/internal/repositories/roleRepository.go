package repositories

import (
	"admin/user/api/internal/pkg/util"
	"context"
	"time"

	"admin/user/api/internal/Initialization"
	"admin/user/api/internal/model"
	"admin/user/api/internal/types"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type RoleRepository interface {
	PagingQuery(ctx context.Context, req types.GetRoleReq) (rsp types.GetRoleRsp, err error)
	FindRoleByIdList(ctx context.Context, idList []string) (roleList []types.RoleList, err error)
	AddRole(ctx context.Context, req types.PostRoleReq) (err error)
	EditRole(ctx context.Context, req types.PutRoleReq) (err error)
	DeleteRole(ctx context.Context, req types.DeleteRoleReq) (err error)
	VerifyRoleIds(ctx context.Context, roleIds []uint) (roleIdsValid bool, err error)
	VerifyRoleUse(ctx context.Context, roleIds []uint) (roleUseValid bool, err error)
}

// NewRoleRepository 初始化
func NewRoleRepository(ds *Initialization.DataSources) RoleRepository {
	return &roleRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type roleRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// PagingQuery 分页查询
func (p *roleRepository) PagingQuery(ctx context.Context, req types.GetRoleReq) (rsp types.GetRoleRsp, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminRole              model.AdminRole
		adminRoleHasPermission model.AdminRoleHasPermission
	)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Table(adminRole.TableName()).
		Where(adminRole.TableName()+"."+model.AdminRoleColumns.PlatformID+" = ?", platformId)

	err = db.Count(&rsp.Total).Error
	if err != nil {
		return
	}

	err = db.Select(
		"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.ID+") "+model.AdminRoleColumns.ID,
		"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.PlatformID+") "+model.AdminRoleColumns.PlatformID,
		"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.RoleName+") "+model.AdminRoleColumns.RoleName,
		"group_concat("+model.AdminRoleHasPermissionColumns.PermissionID+") permission_ids").
		Joins("left join `" + adminRoleHasPermission.TableName() + "` on " +
			adminRole.TableName() + "." + model.AdminRoleColumns.ID + " = " +
			adminRoleHasPermission.TableName() + "." + model.AdminRoleHasPermissionColumns.RoleID).
		Group(adminRole.TableName() + "." + model.AdminRoleColumns.ID).
		Scopes(PageDefault(req.CurrentPage, req.PageSize)).Find(&rsp.RowList).Error

	return
}

// FindRoleByIdList 通过idList 查询角色信息
func (p *roleRepository) FindRoleByIdList(ctx context.Context, idList []string) (roleList []types.RoleList, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminRole              model.AdminRole
		adminRoleHasPermission model.AdminRoleHasPermission
	)

	// sql 逻辑
	err = db.Table(adminRole.TableName()).
		Select(
			"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.ID+") "+model.AdminRoleColumns.ID,
			"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.PlatformID+") "+model.AdminRoleColumns.PlatformID,
			"any_value("+adminRole.TableName()+"."+model.AdminRoleColumns.RoleName+") "+model.AdminRoleColumns.RoleName,
			"group_concat("+model.AdminRoleHasPermissionColumns.PermissionID+") permission_ids").
		Joins("left join `"+adminRoleHasPermission.TableName()+"` on "+
			adminRole.TableName()+"."+model.AdminRoleColumns.ID+" = "+
			adminRoleHasPermission.TableName()+"."+model.AdminRoleHasPermissionColumns.RoleID).
		Where(adminRole.TableName()+"."+model.AdminRoleColumns.ID+" in ?", idList).
		Group(adminRole.TableName() + "." + model.AdminRoleColumns.ID).
		Find(&roleList).Error

	return
}

// AddRole 单条插入
func (p *roleRepository) AddRole(ctx context.Context, req types.PostRoleReq) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util.InterfaceToUint(ctx.Value("userId"))
	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	var (
		adminRoleHasPermission []model.AdminRoleHasPermission
		adminRole              = model.AdminRole{
			RoleName:   req.RoleName,
			PlatformID: platformId,
			Creater:    userId,
		}
	)

	err = tx.Create(&adminRole).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 组装权限数据
	for _, v := range req.PermissionIdList {
		adminRoleHasPermission = append(adminRoleHasPermission, model.AdminRoleHasPermission{
			RoleID:       adminRole.ID,
			PermissionID: v,
		})
	}

	// 插入角色权限中间表
	err = tx.Create(adminRoleHasPermission).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// EditRole 单条修改
func (p *roleRepository) EditRole(ctx context.Context, req types.PutRoleReq) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util.InterfaceToUint(ctx.Value("userId"))

	// sql 逻辑
	err = tx.Model(model.AdminRole{}).Omit(
		model.AdminRoleColumns.ID,
		model.AdminRoleColumns.PlatformID,
		model.AdminRoleColumns.CreatedAt,
		model.AdminRoleColumns.Creater,
	).Where("id = ?", req.Id).Updates(model.AdminRole{
		RoleName:  req.RoleName,
		UpdatedAt: time.Now(),
		Updater:   userId,
	}).Error
	if err != nil {
		return err
	}

	// 删除该角色之前的权限
	err = tx.Where(model.AdminRoleHasPermissionColumns.RoleID+"=?", req.Id).Delete(model.AdminRoleHasPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var adminRoleHasPermission []model.AdminRoleHasPermission
	for _, v := range req.PermissionIdList {
		adminRoleHasPermission = append(adminRoleHasPermission, model.AdminRoleHasPermission{
			RoleID:       req.Id,
			PermissionID: v,
		})
	}

	// 给用户绑定新的角色列表
	err = tx.Create(adminRoleHasPermission).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return
}

// DeleteRole 删除 支持批量
func (p *roleRepository) DeleteRole(ctx context.Context, req types.DeleteRoleReq) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// sql 逻辑
	err = tx.Where(model.AdminRoleColumns.ID+" in ?", req.IdList).Delete(&model.AdminRole{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(model.AdminRoleHasPermissionColumns.RoleID+" in ?", req.IdList).Delete(&model.AdminRoleHasPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// VerifyRoleIds 验证列表中id的是否全部有效
func (p *roleRepository) VerifyRoleIds(ctx context.Context, roleIds []uint) (roleIdsValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model.AdminRole{}).
		Where(model.AdminRoleColumns.ID+" in ?", roleIds).
		Where(model.AdminRoleColumns.PlatformID+" = ?", platformId).
		Count(&total).Error

	if err != nil || total != int64(len(roleIds)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}

// VerifyRoleUse 验证列表中id的是否全部有效
func (p *roleRepository) VerifyRoleUse(ctx context.Context, roleIds []uint) (roleUseValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var total int64
	err = db.Model(model.AdminUserHasRole{}).Where(model.AdminUserHasRoleColumns.RoleID+" in ?", roleIds).Count(&total).Error

	if err != nil {
		return
	}

	if total > 0 {
		roleUseValid = true
	}

	return
}
