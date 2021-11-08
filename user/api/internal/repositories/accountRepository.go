package repositories

import (
	"admin/user/api/internal/Initialization"
	"admin/user/api/internal/config"
	"admin/user/api/internal/model"
	"context"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindOneByAccount(ctx context.Context, account string, platformId uint) (rsp model.AdminUser, err error)
	FindOneById(ctx context.Context, userId uint) (rsp model.AdminUser, err error)
	UpdatePassword(ctx context.Context, userId uint, password string) (err error)
	AddTokenToBlacklist(ctx context.Context, token string, accessExpire int64) (err error)
	CheckPermission(ctx context.Context, userId uint, httpMethod, httpPath string) (isPass bool, err error)
}

// NewAccountRepository 初始化
func NewAccountRepository(ds *Initialization.DataSources) AccountRepository {
	return &accountRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type accountRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func (a *accountRepository) FindOneByAccount(ctx context.Context, account string, platformId uint) (rsp model.AdminUser, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model.AdminUser{}).
		Where(model.AdminUserColumns.PlatformID+" = ?", platformId).
		Where(model.AdminUserColumns.Account+" = ? OR "+model.AdminUserColumns.Email+" = ? OR "+
			model.AdminUserColumns.Phone+" = ?", account, account, account).
		Find(&rsp).Error

	return
}

func (a *accountRepository) FindOneById(ctx context.Context, userId uint) (rsp model.AdminUser, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model.AdminUser{}).Where(model.AdminUserColumns.ID+" = ?", userId).Find(&rsp).Error

	return
}

func (a *accountRepository) UpdatePassword(ctx context.Context, userId uint, password string) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	err = db.Model(model.AdminUser{}).
		Where(model.AdminUserColumns.ID+" = ?", userId).
		Update(model.AdminUserColumns.Password, password).Error

	return
}

func (a *accountRepository) AddTokenToBlacklist(ctx context.Context, token string, accessExpire int64) (err error) {
	// redis 实例
	redis := a.redis.WithContext(ctx)

	redisKey := config.DefaultTokenRedisKey + ":" + time.Now().Format("2006-01-02")

	err = redis.SAdd(redisKey, token).Err()
	if err != nil {
		return
	}

	err = redis.Expire(redisKey, time.Duration(accessExpire)*time.Second).Err()

	return
}

func (a *accountRepository) CheckPermission(ctx context.Context, userId uint, httpMethod, httpPath string) (isPass bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		adminUser              model.AdminUser
		adminUserHasRole       model.AdminUserHasRole
		adminRole              model.AdminRole
		adminRoleHasPermission model.AdminRoleHasPermission
		adminPermissionHasAPI  model.AdminPermissionHasAPI
		adminAPI               model.AdminAPI
		total                  int64
	)

	// 判断是否超管
	err = db.Model(model.AdminUser{}).Select(model.AdminUserColumns.PlatformID, model.AdminUserColumns.IsSuperAdmin).
		Where(model.AdminUserColumns.ID+"=?", userId).
		Find(&adminUser).Error
	if err != nil {
		return
	}

	// 继续判断是不是对应接口平台下的超管
	if adminUser.IsSuperAdmin == uint8(1) { // 管理每个平台的超管都拥有用户微服务的权限
		err = db.Model(model.AdminAPI{}).
			Where(model.AdminAPIColumns.PlatformID+" in ? ", []uint{adminUser.PlatformID, config.InitRouteData.PlatformId}).
			Where(model.AdminAPIColumns.HTTPPath+" = ? ", httpPath).
			Where(model.AdminAPIColumns.HTTPMethod+" = ? ", httpMethod).
			Count(&total).Error
		if total > 0 {
			return true, nil
		}
	}

	// 找出该用户的所有权限  拆成两次查询(每次连表不超过三张)
	var permissionIdList []uint
	err = db.Model(model.AdminUserHasRole{}).
		Select(model.AdminRoleHasPermissionColumns.PermissionID).
		Joins(" left join "+adminRole.TableName()+" on "+
			adminUserHasRole.TableName()+"."+model.AdminUserHasRoleColumns.RoleID+"="+
			adminRole.TableName()+"."+model.AdminRoleColumns.ID).
		Joins(" left join "+adminRoleHasPermission.TableName()+" on "+
			adminRoleHasPermission.TableName()+"."+model.AdminRoleHasPermissionColumns.RoleID+"="+
			adminRole.TableName()+"."+model.AdminRoleColumns.ID).
		Where(model.AdminUserHasRoleColumns.UserID+" = ? ", userId).
		Where(model.AdminRoleHasPermissionColumns.PermissionID + " is not NULL").
		Find(&permissionIdList).Error
	if err != nil || len(permissionIdList) < 1 {
		return
	}

	err = db.Model(model.AdminPermissionHasAPI{}).
		Joins(" left join "+adminAPI.TableName()+" on "+
			adminPermissionHasAPI.TableName()+"."+model.AdminPermissionHasAPIColumns.APIID+"="+
			adminAPI.TableName()+"."+model.AdminAPIColumns.ID).
		Where(model.AdminPermissionHasAPIColumns.PermissionID+" in ?", permissionIdList).
		Where(model.AdminAPIColumns.HTTPMethod+" = ?", httpMethod).
		Where(model.AdminAPIColumns.HTTPPath+" = ?", httpPath).
		Count(&total).Error

	// 找到了对应的接口
	if total > 0 {
		isPass = true
	}

	return
}
