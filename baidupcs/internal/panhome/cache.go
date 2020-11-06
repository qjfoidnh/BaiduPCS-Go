package panhome

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/expires"
	"time"
)

// SetSignExpires 设置sign过期
func (ph *PanHome) SetSignExpires() {
	if ph.signExpires != nil {
		ph.signExpires.SetExpires(true)
	}
}

// CacheSignature 在有效期内返回缓存结果
func (ph *PanHome) CacheSignature() (sign SignRes, err error) {
	if ph.signExpires == nil || ph.signExpires.IsExpires() {
		// 先签名再设置有效期
		ph.signRes, err = ph.Signature()
		if err != nil { // 空指针与空接口不等价
			return nil, err
		}

		ph.signExpires = expires.NewExpires(1 * time.Hour) // 设置一小时有效期
		return ph.signRes, nil
	}

	return ph.signRes, nil
}
