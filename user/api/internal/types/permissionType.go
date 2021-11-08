package types

type GetPermissionReq struct {
	CurrentPage int `form:"currentPage" valid:"range(1|2147483647),required" example:"1"`
	PageSize    int `form:"pageSize" valid:"in(1|10|20|50|100),required" example:"10"`
}

type GetPermissionRsp struct {
	Total   int64            `json:"total"`
	RowList []PermissionList `json:"rowList"`
}

type PermissionList struct {
	Id             uint   `json:"id"`
	PermissionName string `json:"permissionName"`
	ApiIds         string `json:"apiIds"`
}

type FindPermissionByIdsReq struct {
	Ids string `form:"ids"  valid:"required" example:"1,2,3"`
}

type FindPermissionByIdsRsp struct {
	DataList []PermissionList `json:"dataList"`
}

type PostPermissionReq struct {
	PermissionName string `json:"permissionName" valid:"stringlength(1|255),required" example:"产品"`
	ApiIdList      []uint `json:"apiIdList" valid:"required" example:"[1,2]"`
}

type PutPermissionReq struct {
	Id             uint   `json:"id" valid:"range(1|2147483647),required" example:"1"`
	PermissionName string `json:"permissionName" valid:"stringlength(1|255),required" example:"产品"`
	ApiIdList      []uint `json:"apiIdList" valid:"required" example:"[1,2]"`
}

type DeletePermissionReq struct {
	IdList []uint `json:"idList" valid:"required" example:"[1,2]"`
}
