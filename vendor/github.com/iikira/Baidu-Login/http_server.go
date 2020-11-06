package baidulogin

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/astaxie/beego/session"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/json-iterator/go"
	"log"
	"net"
	"net/http"
	"net/url"
)

// StartServer 启动 http 服务
func StartServer(port string) {
	templateFilesBox = rice.MustFindBox("http-files/template")
	libFilesBox = rice.MustFindBox("http-files/static")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/index.html", indexPage)
	http.HandleFunc("/favicon.ico", favicon)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(libFilesBox.HTTPBox())))
	http.HandleFunc("/cgi-bin/baidu/login", execBaiduLogin)
	http.HandleFunc("/cgi-bin/baidu/sendcode", sendCode)
	http.HandleFunc("/cgi-bin/baidu/verifylogin", execVerify)

	fmt.Println("Server is starting...")

	// 初始化 session 管理器
	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
	})

	go globalSessions.GC()

	// Print available URLs.
	for _, address := range pcsutil.ListAddresses() {
		fmt.Printf(
			"URL: %s\n",
			(&url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(address, port),
				Path:   "/",
			}).String(),
		)
	}
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":"+port, nil))
}

// rootHandler 根目录处理
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// 跳转到 /index.html
		w.Header().Set("Location", "/index.html")
		http.Error(w, "", 301)
	} else {
		http.Error(w, "404 Not Found", 404)
	}
}

// execBaiduLogin 发送 百度登录 请求
func execBaiduLogin(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	registerBaiduClient(&sess)   // 如果没有 baiduClinet , 就添加
	defer sess.SessionRelease(w) // 更新 session 储存

	r.ParseForm()                          // 解析 url 传递的参数
	username := r.Form.Get("username")     // 百度 用户名
	password := r.Form.Get("password")     // 密码
	verifycode := r.Form.Get("verifycode") // 图片验证码
	vcodestr := r.Form.Get("vcodestr")     // 与 图片验证码 相对应的字串

	bc, err := getBaiduClient(sess.SessionID()) // 查找该 sessionID 下是否存在 cookiejar
	if err != nil {
		log.Println(err)
		return
	}

	lj := bc.BaiduLogin(username, password, verifycode, vcodestr) //发送登录请求

	// 输出 json 编码
	byteBody, _ := jsoniter.MarshalIndent(&lj, "", " ")
	w.Write(byteBody)
}

// sendCode 发送 获取验证码 请求
func sendCode(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	bc, err := getBaiduClient(sess.SessionID())
	if err != nil {
		log.Println(err)
		return
	}
	r.ParseForm()
	verifyType := r.Form.Get("type")
	token := r.Form.Get("token")
	if token == "" {
		w.Write([]byte(`{"msg":"Token is null."}`))
		return
	}

	msg := bc.SendCodeToUser(verifyType, token)
	w.Write([]byte(`{"msg": "` + msg + `"}`))
}

// execVerifiy 发送 提交验证码 请求
func execVerify(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	bc, err := getBaiduClient(sess.SessionID())

	if err != nil {
		log.Println(err)
		return
	}
	r.ParseForm()                    // 解析 url 传递的参数
	verifyType := r.Form.Get("type") // email/mobile
	token := r.Form.Get("token")     // token 不可或缺
	vcode := r.Form.Get("vcode")     // email/mobile 收到的验证码
	u := r.Form.Get("u")

	lj := bc.VerifyCode(verifyType, token, vcode, u)

	// 输出 json 编码
	byteBody, _ := jsoniter.MarshalIndent(&lj, "", " ")
	w.Write(byteBody)
}
