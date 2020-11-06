package bdcrypto

import (
	"encoding/hex"
)

// RSAEncryptOfWapBaidu 针对 WAP 登录百度的 RSA 加密
func RSAEncryptOfWapBaidu(rsaPublicKeyModulus string, origData []byte) (string, error) {
	ciphertext, err := RSAEncryptNoPadding(rsaPublicKeyModulus, DefaultRSAPublicKeyExponent, BytesReverse(origData))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}
