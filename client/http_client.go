package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// HTTPClient http client
type HTTPClient struct {
	*http.Client
	UserAgent string
}

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
)

func NewHTTPClient() *HTTPClient {

	j, _ := cookiejar.New(nil)
	h := &HTTPClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     j,
		},
		UserAgent: UserAgent,
	}

	return h
}

// SetUserAgent 设置 UserAgent 浏览器标识
func (h *HTTPClient) SetUserAgent(ua string) {
	h.UserAgent = ua
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (h *HTTPClient) SetTimeout(t time.Duration) {
	h.Client.Timeout = t
}

func (h *HTTPClient) Fetch(webUrl string, method string, headers map[string]string, data interface{}) (respBody []byte, err error) {

	var req *http.Request
	var obody io.Reader
	if data != nil {
		switch value := data.(type) {
		case io.Reader:
			obody = value
		case map[string]string:
			query := url.Values{}
			for k := range value {
				query.Set(k, value[k])
			}
			obody = strings.NewReader(query.Encode())
		case map[string]interface{}:
			query := url.Values{}
			for k := range value {
				query.Set(k, fmt.Sprint(value[k]))
			}
			obody = strings.NewReader(query.Encode())
		case map[interface{}]interface{}:
			query := url.Values{}
			for k := range value {
				query.Set(fmt.Sprint(k), fmt.Sprint(value[k]))
			}
			obody = strings.NewReader(query.Encode())
		case string:
			obody = strings.NewReader(value)
		case []byte:
			obody = bytes.NewReader(value[:])
		default:
			return nil, fmt.Errorf("requester.Req: unknown post type: %s", value)
		}
	}

	req, err = http.NewRequest(method, webUrl, obody)
	if err != nil {
		return
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")

	resp, err := h.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//fmt.Println("res Cookies is :",resp.Cookies())
	//u,_:= url.Parse(webUrl)
	//h.Client.Jar.SetCookies(u,resp.Cookies())
	//fmt.Println("cookie jar is:",h.Client.Jar.Cookies(u))
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}
