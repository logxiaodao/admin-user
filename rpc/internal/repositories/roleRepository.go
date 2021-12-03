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

type RoleRepository interface {
	PagingQuery(ctx context.Context, req *user2.GetRoleRequest) (rsp user2.GetRoleResponse, err error)
	FindRoleByIdList(ctx context.Context, idList []string) (roleList []*user2.RoleInfo, err error)
	AddRole(ctx context.Context, req *user2.AddRoleRequest) (err error)
	EditRole(ctx context.Context, req *user2.EditRoleRequest) (err error)
	DeleteRole(ctx context.Context, req *user2.DeleteRoleRequest) (err error)
	VerifyRoleIds(ctx context.Context, roleIds []uint64) (roleIdsValid bool, err error)
	VerifyRoleUse(ctx context.Context, roleIds []uint64) (roleUseValid bool, err error)
}

// NewRoleRepository 初始化
func NewRoleRepository(ds *Initialization2.DataSources) RoleRepository {
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
func (p *roleRepository) PagingQuery(ctx context.Context, req *user2.GetRoleRequest) (rsp user2.GetRoleResponse, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminPermission        model2.AdminPermission
		adminRole              model2.AdminRole
		adminRoleHasPermission model2.AdminRoleHasPermission
	)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Table(adminRole.TableName()).
		Where(adminRole.TableName()+"."+model2.AdminRoleColumns.PlatformID+" = ?", platformId)

	if len(req.Keyword) > 0 {
		db = db.Where(model2.AdminRoleColumns.RoleName + " like '" + req.Keyword + "%'")
	}

	err = db.Count(&rsp.Total).Error
	if err != nil {
		return
	}

	err = db.Select(
		"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.ID+") "+model2.AdminRoleColumns.ID,
		"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.PlatformID+") "+model2.AdminRoleColumns.PlatformID,
		"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.RoleName+") "+model2.AdminRoleColumns.RoleName,
		"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.CreatedAt+") as createdat",
		"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.UpdatedAt+") as updatedat",
		"group_concat("+model2.AdminRoleHasPermissionColumns.PermissionID+") permission_ids",
		"group_concat("+model2.AdminPermissionColumns.PermissionName+") permission_names").
		Joins("left join `" + adminRoleHasPermission.TableName() + "` on " +
			adminRole.TableName() + "." + model2.AdminRoleColumns.ID + " = " +
			adminRoleHasPermission.TableName() + "." + model2.AdminRoleHasPermissionColumns.RoleID).
		Joins("left join `" + adminPermission.TableName() + "` on " +
			adminPermission.TableName() + "." + model2.AdminPermissionColumns.ID + " = " +
			adminRoleHasPermission.TableName() + "." + model2.AdminRoleHasPermissionColumns.PermissionID).
		Group(adminRole.TableName() + "." + model2.AdminRoleColumns.ID).
		Scopes(PageDefault(int(req.CurrentPage), int(req.PageSize))).Find(&rsp.RowList).Error

	return
}

// FindRoleByIdList 通过idList 查询角色信息
func (p *roleRepository) FindRoleByIdList(ctx context.Context, idList []string) (roleList []*user2.RoleInfo, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var (
		adminPermission        model2.AdminPermission
		adminRole              model2.AdminRole
		adminRoleHasPermission model2.AdminRoleHasPermission
	)

	// sql 逻辑
	err = db.Table(adminRole.TableName()).
		Select(
			"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.ID+") "+model2.AdminRoleColumns.ID,
			"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.PlatformID+") "+model2.AdminRoleColumns.PlatformID,
			"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.RoleName+") "+model2.AdminRoleColumns.RoleName,
			"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.CreatedAt+") as createdat",
			"any_value("+adminRole.TableName()+"."+model2.AdminRoleColumns.UpdatedAt+") as updatedat",
			"group_concat("+model2.AdminRoleHasPermissionColumns.PermissionID+") permission_ids",
			"group_concat("+model2.AdminPermissionColumns.PermissionName+") permission_names").
		Joins("left join `"+adminRoleHasPermission.TableName()+"` on "+
			adminRole.TableName()+"."+model2.AdminRoleColumns.ID+" = "+
			adminRoleHasPermission.TableName()+"."+model2.AdminRoleHasPermissionColumns.RoleID).
		Joins("left join `"+adminPermission.TableName()+"` on "+
			adminPermission.TableName()+"."+model2.AdminPermissionColumns.ID+" = "+
			adminRoleHasPermission.TableName()+"."+model2.AdminRoleHasPermissionColumns.PermissionID).
		Where(adminRole.TableName()+"."+model2.AdminRoleColumns.ID+" in ?", idList).
		Group(adminRole.TableName() + "." + model2.AdminRoleColumns.ID).
		Find(&roleList).Error

	return
}

// AddRole 单条插入
func (p *roleRepository) AddRole(ctx context.Context, req *user2.AddRoleRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util2.InterfaceToUint(ctx.Value("userId"))
	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	var (
		adminRoleHasPermission []model2.AdminRoleHasPermission
		adminRole              = model2.AdminRole{
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
		adminRoleHasPermission = append(adminRoleHasPermission, model2.AdminRoleHasPermission{
			RoleID:       adminRole.ID,
			PermissionID: uint(v),
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
func (p *roleRepository) EditRole(ctx context.Context, req *user2.EditRoleRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// 获取userId
	userId := util2.InterfaceToUint(ctx.Value("userId"))

	// sql 逻辑
	err = tx.Model(model2.AdminRole{}).Omit(
		model2.AdminRoleColumns.ID,
		model2.AdminRoleColumns.PlatformID,
		model2.AdminRoleColumns.CreatedAt,
		model2.AdminRoleColumns.Creater,
	).Where("id = ?", req.Id).Updates(model2.AdminRole{
		RoleName:  req.RoleName,
		UpdatedAt: time.Now(),
		Updater:   userId,
	}).Error
	if err != nil {
		return err
	}

	// 删除该角色之前的权限
	err = tx.Where(model2.AdminRoleHasPermissionColumns.RoleID+"=?", req.Id).Delete(model2.AdminRoleHasPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var adminRoleHasPermission []model2.AdminRoleHasPermission
	for _, v := range req.PermissionIdList {
		adminRoleHasPermission = append(adminRoleHasPermission, model2.AdminRoleHasPermission{
			RoleID:       uint(req.Id),
			PermissionID: uint(v),
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
func (p *roleRepository) DeleteRole(ctx context.Context, req *user2.DeleteRoleRequest) (err error) {
	// db 实例
	tx := p.db.WithContext(ctx).Begin()

	// sql 逻辑
	err = tx.Where(model2.AdminRoleColumns.ID+" in ?", req.IdList).Delete(&model2.AdminRole{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(model2.AdminRoleHasPermissionColumns.RoleID+" in ?", req.IdList).Delete(&model2.AdminRoleHasPermission{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// VerifyRoleIds 验证列表中id的是否全部有效
func (p *roleRepository) VerifyRoleIds(ctx context.Context, roleIds []uint64) (roleIdsValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model2.AdminRole{}).
		Where(model2.AdminRoleColumns.ID+" in ?", roleIds).
		Where(model2.AdminRoleColumns.PlatformID+" = ?", platformId).
		Count(&total).Error

	if err != nil || total != int64(len(roleIds)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}

// VerifyRoleUse 验证列表中id的是否全部有效
func (p *roleRepository) VerifyRoleUse(ctx context.Context, roleIds []uint64) (roleUseValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	var total int64
	err = db.Model(model2.AdminUserHasRole{}).Where(model2.AdminUserHasRoleColumns.RoleID+" in ?", roleIds).Count(&total).Error

	if err != nil {
		return
	}

	if total > 0 {
		roleUseValid = true
	}

	return
}
