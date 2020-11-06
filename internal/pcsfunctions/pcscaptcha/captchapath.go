package pcscaptcha

import (
	"os"
	"path/filepath"
)

// CaptchaPath 返回验证码存放路径
func CaptchaPath() string {
	return filepath.Join(os.TempDir(), CaptchaName)
}
