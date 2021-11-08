package types

type LoginOutReq struct {
	AccessToken string `json:"accessToken"`
}

type UpdatePasswordReq struct {
	OldPassword     string `json:"oldPassword"  valid:"required,stringlength(8|16)" example:"12345678"`
	NewPassword     string `json:"newPassword"  valid:"required,stringlength(8|16)" example:"12345678"`
	ConfirmPassword string `json:"confirmPassword"  valid:"required,stringlength(8|16)" example:"12345678"`
}

type (
	CheckPermissionReq struct {
		HttpMethod string `json:"httpMethod" valid:"in(GET|POST|PUT|DELETE),required" example:"GET"`
		HttpPath   string `json:"httpPath" valid:"stringlength(1|255),required" example:"/api"`
	}

	CheckPermissionRsp struct {
		IsPass bool `json:"isPass"`
	}
)
