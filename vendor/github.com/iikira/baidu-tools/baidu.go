package baidu

// Baidu 百度帐号详细情况
type Baidu struct {
	UID      uint64  // 百度ID对应的uid
	Name     string  // 真实ID
	NameShow string  // 显示的用户名(昵称)
	Sex      string  // 性别
	Age      float64 // 帐号年龄
	Auth     *Auth
}

// Auth 百度验证
type Auth struct {
	BDUSS  string // 百度BDUSS
	PTOKEN string
	STOKEN string
}
