package types

type GetApiReq struct {
	CurrentPage int `form:"currentPage" valid:"range(1|2147483647),required" example:"1"`
	PageSize    int `form:"pageSize" valid:"in(1|10|20|50|100),required" example:"10"`
}

type GetApiRsp struct {
	Total   int64     `json:"total"`
	RowList []ApiList `json:"rowList"`
}

type ApiList struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	HttpMethod string `json:"httpMethod"`
	HttpPath   string `json:"httpPath"`
}

type FindApiByIdsReq struct {
	Ids string `form:"ids"  valid:"required" example:"1,2,3"`
}

type FindApiByIdsRsp struct {
	DataList []ApiList `json:"dataList"`
}

type PostApiReq struct {
	PostApiList
}

type BatchPostApiReq struct {
	ItemList []PostApiList `json:"itemList"`
}

type PostApiList struct {
	Name       string `json:"name" valid:"stringlength(1|255),required" example:"root"`
	HttpMethod string `json:"httpMethod" valid:"in(GET|POST|PUT|DELETE),required" example:"get"`
	HttpPath   string `json:"httpPath" valid:"stringlength(1|255),required" example:"/api"`
}

type PutApiReq struct {
	Id         uint   `json:"id" valid:"required,range(1|2147483647)" example:"1"`
	Name       string `json:"name" valid:"stringlength(1|255),required" example:"root"`
	HttpMethod string `json:"httpMethod" valid:"in(GET|POST|PUT|DELETE),required" example:"get"`
	HttpPath   string `json:"httpPath" valid:"stringlength(1|255),required" example:"/api"`
}

type DeleteApiReq struct {
	IdList []uint `json:"idList" valid:"required" example:"1"`
}
