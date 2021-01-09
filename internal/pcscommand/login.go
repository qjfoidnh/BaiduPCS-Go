package pcscommand

import (
	"bytes"
	"fmt"
	baidulogin "github.com/qjfoidnh/Baidu-Login"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcscaptcha"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsliner"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"image/png"
	"io/ioutil"
	"strings"
)

// handleVerifyImg 处理验证码, 下载到本地
func handleVerifyImg(imgURL string) (savePath string, err error) {
	imgContents, err := requester.Fetch("GET", imgURL, nil, nil)
	if err != nil {
		return "", fmt.Errorf("获取验证码失败, 错误: %s", err)
	}

	_, err = png.Decode(bytes.NewReader(imgContents))
	if err != nil {
		return "", fmt.Errorf("验证码解析错误: %s", err)
	}

	savePath = pcscaptcha.CaptchaPath()

	return savePath, ioutil.WriteFile(savePath, imgContents, 0777)
}

// RunLogin 登录百度帐号
func RunLogin(username, password string) (bduss, ptoken, stoken string, cookies string, err error) {
	line := pcsliner.NewLiner()
	defer line.Close()

	bc := baidulogin.NewBaiduClinet()

	if username == "" {
		username, err = line.State.Prompt("请输入百度用户名(手机号/邮箱/用户名), 回车键提交 > ")
		if err != nil {
			return
		}
	}

	if password == "" {
		// liner 的 PasswordPrompt 不安全, 拆行之后密码就会显示出来了
		fmt.Printf("请输入密码(输入的密码无回显, 确认输入完成, 回车提交即可) > ")
		password, err = line.State.PasswordPrompt("")
		if err != nil {
			return
		}
	}

	var vcode_raw, vcode, vcodestr string
	// 移除验证码文件
	defer func() {
		pcscaptcha.RemoveCaptchaPath()
		pcscaptcha.RemoveOldCaptchaPath()
	}()

for_1:
	for i := 0; i < 10; i++ {
	BEGIN:
		lj := bc.BaiduLogin(username, password, vcode_raw, vcodestr)
		switch lj.ErrInfo.No {
		case "0": // 登录成功, 退出循环
			return lj.Data.BDUSS, lj.Data.PToken, lj.Data.SToken, lj.Data.CookieString, nil
		case "400023", "400101": // 需要验证手机或邮箱
			fmt.Printf("\n需要验证手机或邮箱才能登录\n选择一种验证方式\n")
			fmt.Printf("1: 手机: %s\n", lj.Data.Phone)
			fmt.Printf("2: 邮箱: %s\n", lj.Data.Email)
			fmt.Printf("\n")

			var verifyType string
			for et := 0; et < 3; et++ {
				verifyType, err = line.State.Prompt("请输入验证方式 (1 或 2) > ")
				if err != nil {
					return
				}

				switch verifyType {
				case "1":
					verifyType = "mobile"
				case "2":
					verifyType = "email"
				default:
					fmt.Printf("[%d/3] 验证方式不合法\n", et+1)
					continue
				}
				break
			}
			if verifyType != "mobile" && verifyType != "email" {
				err = fmt.Errorf("验证方式不合法")
				return
			}
			msg := ""
			if lj.Data.AuthID != "" {
				msg = bc.SendCodeToUser(verifyType, lj.Data.VerifyURL, lj.Data.AuthID) // 发送验证码
			} else {
				msg = bc.SendCodeToUser2(verifyType, lj.Data.Token)
			}
			fmt.Printf("消息: %s\n\n", msg)
			if strings.Contains(msg, "系统出错") {
				return
			}
			for et := 0; et < 3; et++ {
				vcode, err = line.State.Prompt("请输入接收到的验证码 > ")
				if err != nil {
					return
				}
				nlj := &baidulogin.LoginJSON{}
				if lj.Data.AuthID != "" {
					// 此处 BDUSS 等信息尚未获取到, 仅仅完成了邮箱/电话验证
					nlj = bc.VerifyCode(vcode, verifyType, lj.Data.VerifyURL, lj.Data.AuthID, lj.Data.LoginProxy, lj.Data.AuthSID)
				} else {
					// 此处 BDUSS 等信息已在请求中返回
					nlj = bc.VerifyCode2(verifyType, lj.Data.Token, vcode, lj.Data.U)
				}
				if nlj.ErrInfo.No != "0" {
					fmt.Printf("[%d/3] 错误消息: %s\n\n", et+1, nlj.ErrInfo.Msg)
					if nlj.ErrInfo.No == "-2" { // 需要重发验证码
						return
					}
					continue
				} else {
					vcode_raw = ""
					vcodestr = ""
						goto BEGIN
				}
				// 登录成功
				return nlj.Data.BDUSS, nlj.Data.PToken, nlj.Data.SToken, nlj.Data.CookieString, nil
			}
			break for_1
		case "500001", "500002": // 验证码
			fmt.Printf("\n%s\n", lj.ErrInfo.Msg)
			vcodestr = lj.Data.CodeString
			if vcodestr == "" {
				err = fmt.Errorf("未找到codeString")
				return
			}

			// 图片验证码
			var (
				verifyImgURL = "https://wappass.baidu.com/cgi-bin/genimage?" + vcodestr
				savePath     string
			)

			savePath, err = handleVerifyImg(verifyImgURL)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("打开以下路径, 以查看验证码\n%s\n\n", savePath)
			}

			fmt.Printf("或者打开以下的网址, 以查看验证码\n")
			fmt.Printf("%s\n\n", verifyImgURL)

			vcode_raw, err = line.State.Prompt("请输入验证码 > ")
			if err != nil {
				return
			}
			continue
		default:
			err = fmt.Errorf("错误代码: %s, 消息: %s", lj.ErrInfo.No, lj.ErrInfo.Msg)
			return
		}
	}
	return
}
