package types

type LoginReq struct {
	PlatformID uint   `json:"platformID"  valid:"range(1|2147483647),required" example:"1"`
	Username   string `json:"username"  valid:"required" example:"root"`
	Password   string `json:"password"  valid:"required" example:"123456"`
}

type LoginRes struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	PlatformID   int64  `json:"platformID"`
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	RefreshAfter int64  `json:"refreshAfter"`
}
