
### 1. "登陆"

1. 路由定义

- Url: /user/login
- Method: POST
- Request: `LoginReq`
- Response: `LoginRes`

2. 请求定义


```golang
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```


3. 返回定义


```golang
type LoginRes struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	AccessToken string `json:"accessToken"`
	AccessExpire int64 `json:"accessExpire"`
	RefreshAfter int64 `json:"refreshAfter"`
}
```
  

