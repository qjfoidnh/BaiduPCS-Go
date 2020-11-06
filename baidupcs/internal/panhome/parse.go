package panhome

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"unsafe"
)

var (
	signInfoRE         = regexp.MustCompile(`"sign1":"(.*?)"[\s\S]*"sign3":"(.*?)","timestamp":(\d*?),`)
	ErrCookieInvalid   = errors.New("cookie is invalid")
	ErrUnknownLocation = errors.New("unknown location")
	ErrMatchPanHome    = errors.New("网盘首页数据匹配出错")
)

func (ph *PanHome) getSignInfo() error {
	ph.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	u := *panBaiduComURL
	u.Path = "/disk/home"
	resp, err := ph.client.Req(http.MethodGet, u.String(), nil, map[string]string{
		"User-Agent": PanHomeUserAgent,
	})
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	loc := resp.Header.Get("Location")
	switch loc {
	case "/":
		return ErrCookieInvalid
	case "":
		//pass
	default:
		locU, err := url.Parse(loc)
		if err != nil {
			return ErrUnknownLocation
		}
		if locU.Host == "passport.baidu.com" {
			return ErrCookieInvalid
		}
		return ErrUnknownLocation
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	matchRes := signInfoRE.FindSubmatch(body)
	if len(matchRes) <= 3 {
		return ErrMatchPanHome
	}

	ph.sign1 = []rune(*(*string)(unsafe.Pointer(&matchRes[1])))
	ph.sign3 = []rune(*(*string)(unsafe.Pointer(&matchRes[2])))
	ph.timestamp = *(*string)(unsafe.Pointer(&matchRes[3]))
	return nil
}
