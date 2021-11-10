package errorx

import "github.com/pkg/errors"

//400000  业务逻辑错误
//410000  权限错误
//420000  参数错误
//430000  限制错误
//500000  服务错误
//510000  数据库错误
//520000  io错误
//530000  第三方调用错误
//540000  网络错误

const (
	Success = 0

	LogicalError = 400000 // 业务逻辑错误

	PermissionError       = 410000 // 权限错误
	ExpiredToken          = 410001
	IncorrectToken        = 410002
	InconsistentPasswords = 410003
	WrongPassword         = 410004
	PermissionDenied      = 410005

	ParameterError         = 420000 // 参数错误
	InvalidId              = 420001
	InvalidRoleId          = 420005
	InvalidAdminId         = 420002
	InvalidApiId           = 420003
	InvalidPermissionId    = 420004
	PermissionInUse        = 420006
	ApiInUse               = 420007
	RoleInUse              = 420008
	ApiExists              = 420009
	RoleExists             = 420010
	PermissionExists       = 420011
	AdminExists            = 420012
	ParameterBindingFailed = 420013

	LimitError = 430000 // 限制错误

	ServiceError = 500000 // 服务错误

	DatabaseError = 510000 // 数据库错误

	IoError = 520000 // io错误

	ThirdPartyError = 530000 // 第三方调用错误

	NetworkError = 540000 // 网络错误

)

var (
	DefaultCodeMessage = map[int]Msg{
		Success: {
			En: "Success",
			Zh: "成功",
		},
		LogicalError: {
			En: "Logical Error",
			Zh: "业务逻辑错误",
		},
		PermissionError: {
			En: "Success",
			Zh: "权限错误",
		},
		ParameterError: {
			En: "Parameter Error",
			Zh: "参数错误",
		},
		LimitError: {
			En: "Limit Error",
			Zh: "限制错误",
		},
		ServiceError: {
			En: "Service Error",
			Zh: "服务器错误",
		},
		DatabaseError: {
			En: "Database Error",
			Zh: "数据库错误",
		},
		IoError: {
			En: "Io Error",
			Zh: "输入输出错误",
		},
		ThirdPartyError: {
			En: "Third Party Error",
			Zh: "第三方调用错误",
		},
		NetworkError: {
			En: "Network Error",
			Zh: "网络错误",
		},
	}

	CodeMessage = map[int]Msg{
		ExpiredToken: {
			En: "token has expired",
			Zh: "token已失效",
		},
		IncorrectToken: {
			En: "Incorrect token",
			Zh: "不正确的token",
		},
		PermissionDenied: {
			En: "Permission denied",
			Zh: "没有权限",
		},
		InconsistentPasswords: {
			En: "Inconsistent passwords",
			Zh: "两次密码不一致",
		},
		WrongPassword: {
			En: "Wrong password",
			Zh: "密码错误",
		},
		ParameterBindingFailed: {
			En: "Parameter binding failed",
			Zh: "参数绑定失败",
		},
		InvalidId: {
			En: "There is an invalid id",
			Zh: "存在无效id",
		},
		InvalidAdminId: {
			En: "There is an invalid adminId",
			Zh: "存在无效adminId",
		},
		InvalidApiId: {
			En: "There is an invalid apiId",
			Zh: "存在无效apiId",
		},
		InvalidRoleId: {
			En: "There is an invalid roleId",
			Zh: "存在无效roleId",
		},
		InvalidPermissionId: {
			En: "There is an invalid permissionId",
			Zh: "存在无效permissionId",
		},
		PermissionInUse: {
			En: "The permission group is in use, please unbind the permission group and the role first",
			Zh: "权限组正在使用中，请先解除该权限组和角色的绑定",
		},
		RoleInUse: {
			En: "The role is in use, please remove the binding between the role and the user first",
			Zh: "角色正在使用中，请先接除该角色和用户的绑定",
		},
		ApiInUse: {
			En: "The api is in use, please unbind the api and the permission group first",
			Zh: "api接口正在使用中，请先解除该api和权限组的绑定",
		},
		ApiExists: {
			En: "Api Exists",
			Zh: "api接口已存在",
		},
		PermissionExists: {
			En: "Permission Exists",
			Zh: "权限已存在",
		},
		RoleExists: {
			En: "Role Exists",
			Zh: "角色已存在",
		},
		AdminExists: {
			En: "Admin Exists",
			Zh: "用户已存在",
		},
	}
)

// GetErrorByCode 通过code生成错误
func GetErrorByCode(code int) (err error) {

	if _, ok := DefaultCodeMessage[code]; ok {
		msg := DefaultCodeMessage[code]
		err = errors.New(msg.Error())
	}

	if _, ok := CodeMessage[code]; ok {
		msg := CodeMessage[code]
		err = errors.New(msg.Error())
	}

	return
}
