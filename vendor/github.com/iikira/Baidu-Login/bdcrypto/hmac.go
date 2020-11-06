package bdcrypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// HmacSHA1 HMAC-SHA-1签名认证
func HmacSHA1(key, origData []byte) (sum []byte) {
	mac := hmac.New(sha1.New, key)
	mac.Write(origData)
	return mac.Sum(nil)
}

// HmacSHA256 HMAC-SHA-256签名认证
func HmacSHA256(key, origData []byte) (sum []byte) {
	mac := hmac.New(sha256.New, key)
	mac.Write(origData)
	return mac.Sum(nil)
}

// HmacSHA512 HMAC-SHA-512签名认证
func HmacSHA512(key, origData []byte) (sum []byte) {
	mac := hmac.New(sha512.New, key)
	mac.Write(origData)
	return mac.Sum(nil)
}

// HmacMD5 HMAC-SHA512-签名认证
func HmacMD5(key, origData []byte) (sum []byte) {
	mac := hmac.New(md5.New, key)
	mac.Write(origData)
	return mac.Sum(nil)
}
