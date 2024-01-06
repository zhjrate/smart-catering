package model

// User 用户信息
type User struct {
	//ID       int64  //用户id
	UserName string //用户名
	Password string //用户密码
	Token    string
	Expiry   int64
	Times    int64 //判断用户第几次登录
}
