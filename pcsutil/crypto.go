package pcsutil

import (
	"fmt"
	"github.com/iikira/Baidu-Login/bdcrypto"
	"io"
	"os"
	"strings"
)

// CryptoMethodSupport 检测是否支持加密解密方法
func CryptoMethodSupport(method string) bool {
	switch method {
	case "aes-128-ctr", "aes-192-ctr", "aes-256-ctr", "aes-128-cfb", "aes-192-cfb", "aes-256-cfb", "aes-128-ofb", "aes-192-ofb", "aes-256-ofb":
		return true
	}

	return false
}

// EncryptFile 加密本地文件
func EncryptFile(method string, key []byte, filePath string, isGzip bool) (encryptedFilePath string, err error) {
	if !CryptoMethodSupport(method) {
		return "", fmt.Errorf("unknown encrypt method: %s", method)
	}

	if isGzip {
		err = bdcrypto.GZIPCompressFile(filePath)
		if err != nil {
			return
		}
	}

	plainFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return
	}

	defer plainFile.Close()

	var cipherReader io.Reader
	switch method {
	case "aes-128-ctr":
		cipherReader, err = bdcrypto.Aes128CTREncrypt(bdcrypto.Convert16bytes(key), plainFile)
	case "aes-192-ctr":
		cipherReader, err = bdcrypto.Aes192CTREncrypt(bdcrypto.Convert24bytes(key), plainFile)
	case "aes-256-ctr":
		cipherReader, err = bdcrypto.Aes256CTREncrypt(bdcrypto.Convert32bytes(key), plainFile)
	case "aes-128-cfb":
		cipherReader, err = bdcrypto.Aes128CFBEncrypt(bdcrypto.Convert16bytes(key), plainFile)
	case "aes-192-cfb":
		cipherReader, err = bdcrypto.Aes192CFBEncrypt(bdcrypto.Convert24bytes(key), plainFile)
	case "aes-256-cfb":
		cipherReader, err = bdcrypto.Aes256CFBEncrypt(bdcrypto.Convert32bytes(key), plainFile)
	case "aes-128-ofb":
		cipherReader, err = bdcrypto.Aes128OFBEncrypt(bdcrypto.Convert16bytes(key), plainFile)
	case "aes-192-ofb":
		cipherReader, err = bdcrypto.Aes192OFBEncrypt(bdcrypto.Convert24bytes(key), plainFile)
	case "aes-256-ofb":
		cipherReader, err = bdcrypto.Aes256OFBEncrypt(bdcrypto.Convert32bytes(key), plainFile)
	default:
		return "", fmt.Errorf("unknown encrypt method: %s", method)
	}

	if err != nil {
		return
	}

	plainFileInfo, err := plainFile.Stat()
	if err != nil {
		return
	}

	encryptedFilePath = filePath + ".encrypt"
	encryptedFile, err := os.OpenFile(encryptedFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, plainFileInfo.Mode())
	if err != nil {
		return
	}

	defer encryptedFile.Close()

	_, err = io.Copy(encryptedFile, cipherReader)
	if err != nil {
		return
	}

	os.Remove(filePath)

	return encryptedFilePath, nil
}

// DecryptFile 加密本地文件
func DecryptFile(method string, key []byte, filePath string, isGzip bool) (decryptedFilePath string, err error) {
	if !CryptoMethodSupport(method) {
		return "", fmt.Errorf("unknown decrypt method: %s", method)
	}

	cipherFile, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	defer cipherFile.Close()

	var plainReader io.Reader
	switch method {
	case "aes-128-ctr":
		plainReader, err = bdcrypto.Aes128CTRDecrypt(bdcrypto.Convert16bytes(key), cipherFile)
	case "aes-192-ctr":
		plainReader, err = bdcrypto.Aes192CTRDecrypt(bdcrypto.Convert24bytes(key), cipherFile)
	case "aes-256-ctr":
		plainReader, err = bdcrypto.Aes256CTRDecrypt(bdcrypto.Convert32bytes(key), cipherFile)
	case "aes-128-cfb":
		plainReader, err = bdcrypto.Aes128CFBDecrypt(bdcrypto.Convert16bytes(key), cipherFile)
	case "aes-192-cfb":
		plainReader, err = bdcrypto.Aes192CFBDecrypt(bdcrypto.Convert24bytes(key), cipherFile)
	case "aes-256-cfb":
		plainReader, err = bdcrypto.Aes256CFBDecrypt(bdcrypto.Convert32bytes(key), cipherFile)
	case "aes-128-ofb":
		plainReader, err = bdcrypto.Aes128OFBDecrypt(bdcrypto.Convert16bytes(key), cipherFile)
	case "aes-192-ofb":
		plainReader, err = bdcrypto.Aes192OFBDecrypt(bdcrypto.Convert24bytes(key), cipherFile)
	case "aes-256-ofb":
		plainReader, err = bdcrypto.Aes256OFBDecrypt(bdcrypto.Convert32bytes(key), cipherFile)
	default:
		return "", fmt.Errorf("unknown decrypt method: %s", method)
	}

	if err != nil {
		return
	}

	cipherFileInfo, err := cipherFile.Stat()
	if err != nil {
		return
	}

	decryptedFilePath = strings.TrimSuffix(filePath, ".encrypt")
	decryptedTmpFilePath := decryptedFilePath + ".decrypted"
	decryptedTmpFile, err := os.OpenFile(decryptedTmpFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, cipherFileInfo.Mode())
	if err != nil {
		return
	}

	_, err = io.Copy(decryptedTmpFile, plainReader)
	if err != nil {
		return
	}

	defer decryptedTmpFile.Close()

	if isGzip {
		err = bdcrypto.GZIPUnompressFile(decryptedTmpFilePath)
		if err != nil {
			os.Remove(decryptedTmpFilePath)
			return
		}

		// 删除已加密的文件
		os.Remove(filePath)
	}

	if filePath != decryptedFilePath {
		os.Rename(decryptedTmpFilePath, decryptedFilePath)
	} else {
		decryptedFilePath = decryptedTmpFilePath
	}

	return decryptedFilePath, nil
}
