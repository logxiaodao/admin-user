package types

type GetAdminReq struct {
	CurrentPage int `form:"currentPage" valid:"range(1|2147483647),required" example:"1"`
	PageSize    int `form:"pageSize" valid:"in(1|10|20|50|100),required" example:"10"`
}

type GetAdminRsp struct {
	Total   int64       `json:"total"`
	RowList []AdminList `json:"rowList"`
}

type AdminList struct {
	Id       uint   `json:"id"`
	Account  string `json:"account"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	RoleIds  string `json:"roleIds"`
}

type FindAdminByIdsReq struct {
	Ids string `form:"ids"  valid:"required" example:"1,2,3"`
}

type FindAdminByIdsRsp struct {
	DataList []AdminList `json:"dataList"`
}

type PostAdminReq struct {
	Account    string `json:"account"  valid:"stringlength(1|255),required" example:"root"`
	Password   string `json:"password"  valid:"stringlength(8|16),optional" example:"12345678"`
	NickName   string `json:"nickName" valid:"stringlength(1|255),required" example:"root"`
	Phone      string `json:"phone" valid:"stringlength(1|20),required" example:"18598814577"`
	Email      string `json:"email" valid:"stringlength(1|255),email,required" example:"1212231@qq.com"`
	RoleIdList []uint `json:"roleIdList" valid:"required" example:"[1,2,3]"`
}

type PutAdminReq struct {
	Id         uint   `json:"id" valid:"range(1|2147483647),required" example:"1"`
	Account    string `json:"account"  valid:"stringlength(1|255),required" example:"root"`
	NickName   string `json:"nickName" valid:"stringlength(1|255),required" example:"root"`
	Phone      string `json:"phone" valid:"stringlength(1|20),required" example:"18598814577"`
	Email      string `json:"email" valid:"email,required" example:"1212231@qq.com"`
	RoleIdList []uint `json:"roleIdList" valid:"required" example:"[1,2,3]"`
}

type DeleteAdminReq struct {
	IdList []uint `json:"idList" valid:"required" example:"1"`
}
