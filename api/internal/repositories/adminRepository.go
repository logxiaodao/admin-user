package repositories

import (
	safe2 "admin-user/api/common/safe"
	Initialization2 "admin-user/api/internal/Initialization"
	config2 "admin-user/api/internal/config"
	model2 "admin-user/api/internal/model"
	util2 "admin-user/api/internal/pkg/util"
	types2 "admin-user/api/internal/types"
	"context"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type AdminRepository interface {
	PagingQuery(ctx context.Context, req types2.GetAdminReq) (rsp types2.GetAdminRsp, err error)
	FindAdminByIdList(ctx context.Context, idList []string) (roleList []types2.AdminList, err error)
	AddAdmin(ctx context.Context, req types2.PostAdminReq) (err error)
	EditAdmin(ctx context.Context, req types2.PutAdminReq) (err error)
	DeleteAdmin(ctx context.Context, req types2.DeleteAdminReq) (err error)
	VerifyAdminIds(ctx context.Context, adminIds []uint) (adminIdsValid bool, err error)
}

// NewAdminRepository 初始化
func NewAdminRepository(ds *Initialization2.DataSources) AdminRepository {
	return &adminRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type adminRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// PagingQuery 分页查询
func (a *adminRepository) PagingQuery(ctx context.Context, req types2.GetAdminReq) (rsp types2.GetAdminRsp, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		adminUser        model2.AdminUser
		adminUserHasRole model2.AdminUserHasRole
	)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// sql 逻辑
	db = db.Table(adminUser.TableName()).
		Where(adminUser.TableName()+"."+model2.AdminUserColumns.PlatformID+"=?", platformId)

	err = db.Count(&rsp.Total).Error
	if err != nil {
		return
	}

	err = db.Select(
		"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.ID+") "+model2.AdminUserColumns.ID,
		"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Account+") "+model2.AdminUserColumns.Account,
		"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.NickName+") "+model2.AdminUserColumns.NickName,
		"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Phone+") "+model2.AdminUserColumns.Phone,
		"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Email+") "+model2.AdminUserColumns.Email,
		"group_concat("+model2.AdminUserHasRoleColumns.RoleID+") role_ids").
		Joins("left join " + adminUserHasRole.TableName() + " on " +
			adminUser.TableName() + "." + model2.AdminUserColumns.ID + "=" +
			adminUserHasRole.TableName() + "." + model2.AdminUserHasRoleColumns.UserID).
		Group(adminUser.TableName() + "." + model2.AdminUserColumns.ID).
		Scopes(PageDefault(req.CurrentPage, req.PageSize)).Find(&rsp.RowList).Error

	return
}

// FindAdminByIdList 通过idList 查询角色信息
func (a *adminRepository) FindAdminByIdList(ctx context.Context, idList []string) (userList []types2.AdminList, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		adminUser        model2.AdminUser
		adminUserHasRole model2.AdminUserHasRole
	)

	// sql 逻辑
	err = db.Table(adminUser.TableName()).
		Select(
			"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.ID+") "+model2.AdminUserColumns.ID,
			"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Account+") "+model2.AdminUserColumns.Account,
			"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.NickName+") "+model2.AdminUserColumns.NickName,
			"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Phone+") "+model2.AdminUserColumns.Phone,
			"any_value("+adminUser.TableName()+"."+model2.AdminUserColumns.Email+") "+model2.AdminUserColumns.Email,
			"group_concat("+model2.AdminUserHasRoleColumns.RoleID+") role_ids").
		Joins("left join "+adminUserHasRole.TableName()+" on "+
			adminUser.TableName()+"."+model2.AdminUserColumns.ID+"="+
			adminUserHasRole.TableName()+"."+model2.AdminUserHasRoleColumns.UserID).
		Where(adminUser.TableName()+"."+model2.AdminUserColumns.ID+" in ?", idList).
		Group(adminUser.TableName() + "." + model2.AdminUserColumns.ID).
		Find(&userList).Error

	return
}

// AddAdmin 单条插入
func (a *adminRepository) AddAdmin(ctx context.Context, req types2.PostAdminReq) (err error) {

	// db 实例
	tx := a.db.WithContext(ctx).Begin()
	if req.Password == "" { // 给默认密码
		req.Password = config2.DefaultPassword
	}

	userId := util2.InterfaceToUint(ctx.Value("userId"))
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	// 密码哈希加盐
	pwd, err := safe2.GenHashPassword(req.Password)
	if err != nil {
		return
	}

	// 插入平台
	var (
		adminUserHasRoles []model2.AdminUserHasRole
		adminUser         = model2.AdminUser{
			Account:    req.Account,
			Password:   pwd,
			NickName:   req.NickName,
			Phone:      req.Phone,
			Email:      req.Email,
			PlatformID: platformId,
			Creater:    userId,
		}
	)

	// sql 逻辑   插入用户表
	err = tx.Create(&adminUser).Error
	if err != nil {
		tx.Rollback()
		return
	}

	// 组装角色数据
	for _, v := range req.RoleIdList {
		adminUserHasRoles = append(adminUserHasRoles, model2.AdminUserHasRole{
			UserID: adminUser.ID,
			RoleID: v,
		})
	}

	// 插入用户角色中间表
	err = tx.Create(adminUserHasRoles).Error
	if err != nil {
		return
	}

	tx.Commit()
	return
}

// EditAdmin 单条修改
func (a *adminRepository) EditAdmin(ctx context.Context, req types2.PutAdminReq) (err error) {
	// db 实例
	tx := a.db.WithContext(ctx).Begin()

	userId := util2.InterfaceToUint(ctx.Value("userId"))

	// sql 逻辑
	err = tx.Model(model2.AdminUser{}).Omit(
		model2.AdminUserColumns.ID,
		model2.AdminUserColumns.Password,
		model2.AdminUserColumns.IsSuperAdmin,
		model2.AdminUserColumns.IsBan,
		model2.AdminUserColumns.CreatedAt,
		model2.AdminUserColumns.Creater,
	).Where("id = ?", req.Id).Updates(model2.AdminUser{
		Account:   req.Account,
		NickName:  req.NickName,
		Phone:     req.Phone,
		Email:     req.Email,
		UpdatedAt: time.Now(),
		Updater:   userId,
	}).Error
	if err != nil {
		return
	}

	// 删除该用户之前的角色
	err = tx.Where(model2.AdminUserHasRoleColumns.UserID+"=?", req.Id).Delete(model2.AdminUserHasRole{}).Error
	if err != nil {
		tx.Rollback()
		return
	}

	var adminUserHasRoles []model2.AdminUserHasRole
	for _, v := range req.RoleIdList {
		adminUserHasRoles = append(adminUserHasRoles, model2.AdminUserHasRole{UserID: req.Id, RoleID: v})
	}

	// 给用户绑定新的角色列表
	err = tx.Create(adminUserHasRoles).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

// DeleteAdmin 删除 支持批量
func (a *adminRepository) DeleteAdmin(ctx context.Context, req types2.DeleteAdminReq) (err error) {
	// db 实例
	tx := a.db.WithContext(ctx).Begin()

	// 删除用户
	err = tx.Where(model2.AdminUserColumns.ID+" in ?", req.IdList).
		Delete(&model2.AdminUser{}).Error
	if err != nil {
		tx.Rollback()
		return
	}

	// 解除用户绑定的角色
	err = tx.Where(model2.AdminUserHasRoleColumns.UserID+" in ?", req.IdList).Delete(&model2.AdminUserHasRole{}).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

// VerifyAdminIds 验证列表中id的是否全部有效
func (p *adminRepository) VerifyAdminIds(ctx context.Context, adminIds []uint) (adminIdsValid bool, err error) {
	// db 实例
	db := p.db.WithContext(ctx)

	// 获取jwt中的平台id
	platformId := util2.InterfaceToUint(ctx.Value("platformID"))

	var total int64
	err = db.Model(model2.AdminUser{}).
		Where(model2.AdminUserColumns.ID+" in ?", adminIds).
		Where(model2.AdminUserColumns.PlatformID+"=?", platformId).
		Count(&total).Error

	if err != nil || total != int64(len(adminIds)) { // 严格限制id列表全部为有效id
		return false, err
	}

	return true, nil
}
