package message

type LoginReq struct {
	Username string `json:"username" binding:"required,alphanum"`      //用户名称
	Password string `json:"password"  binding:"required,min=6,max=12"` //用户密码
}
