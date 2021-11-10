package repositories

import (
	Initialization2 "admin-user/rpc/internal/Initialization"
	config2 "admin-user/rpc/internal/config"
	model2 "admin-user/rpc/internal/model"
	"context"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindOneByAccount(ctx context.Context, account string, platformId int64) (rsp model2.AdminUser, err error)
	FindOneById(ctx context.Context, userId int64) (rsp model2.AdminUser, err error)
	UpdatePassword(ctx context.Context, userId int64, password string) (err error)
	AddTokenToBlacklist(ctx context.Context, token string, accessExpire int64) (err error)
	CheckPermission(ctx context.Context, userId int64, httpMethod, httpPath string) (isPass bool, err error)
}

// NewAccountRepository 初始化
func NewAccountRepository(ds *Initialization2.DataSources) AccountRepository {
	return &accountRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type accountRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func (a *accountRepository) FindOneByAccount(ctx context.Context, account string, platformId int64) (rsp model2.AdminUser, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model2.AdminUser{}).
		Where(model2.AdminUserColumns.PlatformID+" = ?", platformId).
		Where(model2.AdminUserColumns.Account+" = ? OR "+model2.AdminUserColumns.Email+" = ? OR "+
			model2.AdminUserColumns.Phone+" = ?", account, account, account).
		Find(&rsp).Error

	return
}

func (a *accountRepository) FindOneById(ctx context.Context, userId int64) (rsp model2.AdminUser, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	// sql 逻辑
	err = db.Model(model2.AdminUser{}).Where(model2.AdminUserColumns.ID+" = ?", userId).Find(&rsp).Error

	return
}

func (a *accountRepository) UpdatePassword(ctx context.Context, userId int64, password string) (err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	err = db.Model(model2.AdminUser{}).
		Where(model2.AdminUserColumns.ID+" = ?", userId).
		Update(model2.AdminUserColumns.Password, password).Error

	return
}

func (a *accountRepository) AddTokenToBlacklist(ctx context.Context, token string, accessExpire int64) (err error) {
	// redis 实例
	redis := a.redis.WithContext(ctx)

	redisKey := config2.DefaultTokenRedisKey + ":" + time.Now().Format("2006-01-02")

	err = redis.SAdd(redisKey, token).Err()
	if err != nil {
		return
	}

	err = redis.Expire(redisKey, time.Duration(accessExpire)*time.Second).Err()

	return
}

func (a *accountRepository) CheckPermission(ctx context.Context, userId int64, httpMethod, httpPath string) (isPass bool, err error) {
	// db 实例
	db := a.db.WithContext(ctx)

	var (
		adminUser              model2.AdminUser
		adminUserHasRole       model2.AdminUserHasRole
		adminRole              model2.AdminRole
		adminRoleHasPermission model2.AdminRoleHasPermission
		adminPermissionHasAPI  model2.AdminPermissionHasAPI
		adminAPI               model2.AdminAPI
		total                  int64
	)

	// 判断是否超管
	err = db.Model(model2.AdminUser{}).Select(model2.AdminUserColumns.PlatformID, model2.AdminUserColumns.IsSuperAdmin).
		Where(model2.AdminUserColumns.ID+"=?", userId).
		Find(&adminUser).Error
	if err != nil {
		return
	}

	// 继续判断是不是对应接口平台下的超管
	if adminUser.IsSuperAdmin == uint8(1) { // 管理每个平台的超管都拥有用户微服务的权限
		err = db.Model(model2.AdminAPI{}).
			Where(model2.AdminAPIColumns.PlatformID+" in ? ", []int64{adminUser.PlatformID, config2.InitRouteData.PlatformId}).
			Where(model2.AdminAPIColumns.HTTPPath+" = ? ", httpPath).
			Where(model2.AdminAPIColumns.HTTPMethod+" = ? ", httpMethod).
			Count(&total).Error
		if total > 0 {
			return true, nil
		}
	}

	// 找出该用户的所有权限  拆成两次查询(每次连表不超过三张)
	var permissionIdList []uint
	err = db.Model(model2.AdminUserHasRole{}).
		Select(model2.AdminRoleHasPermissionColumns.PermissionID).
		Joins(" left join "+adminRole.TableName()+" on "+
			adminUserHasRole.TableName()+"."+model2.AdminUserHasRoleColumns.RoleID+"="+
			adminRole.TableName()+"."+model2.AdminRoleColumns.ID).
		Joins(" left join "+adminRoleHasPermission.TableName()+" on "+
			adminRoleHasPermission.TableName()+"."+model2.AdminRoleHasPermissionColumns.RoleID+"="+
			adminRole.TableName()+"."+model2.AdminRoleColumns.ID).
		Where(model2.AdminUserHasRoleColumns.UserID+" = ? ", userId).
		Where(model2.AdminRoleHasPermissionColumns.PermissionID + " is not NULL").
		Find(&permissionIdList).Error
	if err != nil || len(permissionIdList) < 1 {
		return
	}

	err = db.Model(model2.AdminPermissionHasAPI{}).
		Joins(" left join "+adminAPI.TableName()+" on "+
			adminPermissionHasAPI.TableName()+"."+model2.AdminPermissionHasAPIColumns.APIID+"="+
			adminAPI.TableName()+"."+model2.AdminAPIColumns.ID).
		Where(model2.AdminPermissionHasAPIColumns.PermissionID+" in ?", permissionIdList).
		Where(model2.AdminAPIColumns.HTTPMethod+" = ?", httpMethod).
		Where(model2.AdminAPIColumns.HTTPPath+" = ?", httpPath).
		Count(&total).Error

	// 找到了对应的接口
	if total > 0 {
		isPass = true
	}

	return
}
