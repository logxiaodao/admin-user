package types

type GetRoleReq struct {
	CurrentPage int `form:"currentPage" valid:"range(1|2147483647),required" example:"1"`
	PageSize    int `form:"pageSize" valid:"in(1|10|20|50|100),required" example:"10"`
}

type GetRoleRsp struct {
	Total   int64      `json:"total"`
	RowList []RoleList `json:"rowList"`
}

type RoleList struct {
	Id            uint   `json:"id"`
	RoleName      string `json:"roleName"`
	PermissionIds string `json:"permissionIds"`
}

type FindRoleByIdsReq struct {
	Ids string `form:"ids"  valid:"required" example:"1,2,3"`
}

type FindRoleByIdsRsp struct {
	DataList []RoleList `json:"dataList"`
}

type PostRoleReq struct {
	RoleName         string `json:"roleName"   valid:"stringlength(1|50),required" example:"产品"`
	PermissionIdList []uint `json:"permissionIdList" valid:"required" example:"[1,2,3]"`
}

type PutRoleReq struct {
	Id               uint   `json:"id" valid:"range(1|2147483647),required" example:"1"`
	RoleName         string `json:"roleName" valid:"stringlength(1|50),required" example:"产品"`
	PermissionIdList []uint `json:"permissionIdList" valid:"required" example:"[1,2,3]"`
}

type DeleteRoleReq struct {
	IdList []uint `json:"idList" valid:"required" example:"[1,2,3]"`
}
