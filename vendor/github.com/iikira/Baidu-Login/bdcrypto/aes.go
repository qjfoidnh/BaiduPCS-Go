package bdcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"fmt"
	"github.com/iikira/Baidu-Login/bdcrypto/ecb"
	"io"
)

// AesMode AES 工作模式
type AesMode int

const (
	// AesECB ecb 模式
	AesECB AesMode = iota
	// AesCBC cbc 模式
	AesCBC
	// AesCTR ctr 模式
	AesCTR
	// AesCFB cfb 模式
	AesCFB
	// AesOFB ofb 模式
	AesOFB
)

// Aes128ECBEncrypt aes-128-ecb 加密
func Aes128ECBEncrypt(key [16]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesECBEnc(key[:], plaintext)
}

// Aes128ECBDecrypt aes-128-ecb 解密
func Aes128ECBDecrypt(key [16]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesECBDec(key[:], ciphertext)
}

// Aes192ECBEncrypt aes-192-ecb 加密
func Aes192ECBEncrypt(key [24]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesECBEnc(key[:], plaintext)
}

// Aes192ECBDecrypt aes-192-ecb 解密
func Aes192ECBDecrypt(key [24]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesECBDec(key[:], ciphertext)
}

// Aes256ECBEncrypt aes-256-ecb 加密
func Aes256ECBEncrypt(key [32]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesECBEnc(key[:], plaintext)
}

// Aes256ECBDecrypt aes-256-ecb 解密
func Aes256ECBDecrypt(key [32]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesECBDec(key[:], ciphertext)
}

// Aes128CBCEncrypt aes-128-cbc 加密
func Aes128CBCEncrypt(key [16]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesCBCEnc(key[:], plaintext)
}

// Aes128CBCDecrypt aes-128-cbc 解密
func Aes128CBCDecrypt(key [16]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesCBCDec(key[:], ciphertext)
}

// Aes192CBCEncrypt aes-192-cbc 加密
func Aes192CBCEncrypt(key [24]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesCBCEnc(key[:], plaintext)
}

// Aes192CBCDecrypt aes-192-cbc 解密
func Aes192CBCDecrypt(key [24]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesCBCDec(key[:], ciphertext)
}

// Aes256CBCEncrypt aes-256-cbc 加密
func Aes256CBCEncrypt(key [32]byte, plaintext []byte) (ciphertext []byte, err error) {
	return aesCBCEnc(key[:], plaintext)
}

// Aes256CBCDecrypt aes-256-cbc 解密
func Aes256CBCDecrypt(key [32]byte, ciphertext []byte) (plaintext []byte, err error) {
	return aesCBCDec(key[:], ciphertext)
}

// Aes128CTREncrypt aes-128-ctr 加密
func Aes128CTREncrypt(key [16]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCTREnc(key[:], plainReader)
}

// Aes128CTRDecrypt aes-128-ctr 解密
func Aes128CTRDecrypt(key [16]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCTRDec(key[:], cipherReader)
}

// Aes192CTREncrypt aes-192-ctr 加密
func Aes192CTREncrypt(key [24]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCTREnc(key[:], plainReader)
}

// Aes192CTRDecrypt aes-192-ctr 解密
func Aes192CTRDecrypt(key [24]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCTRDec(key[:], cipherReader)
}

// Aes256CTREncrypt aes-256-ctr 加密
func Aes256CTREncrypt(key [32]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCTREnc(key[:], plainReader)
}

// Aes256CTRDecrypt aes-256-ctr 解密
func Aes256CTRDecrypt(key [32]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCTRDec(key[:], cipherReader)
}

// Aes128CFBEncrypt aes-128-cfb 加密
func Aes128CFBEncrypt(key [16]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCFBEnc(key[:], plainReader)
}

// Aes128CFBDecrypt aes-128-cfb 解密
func Aes128CFBDecrypt(key [16]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCFBDec(key[:], cipherReader)
}

// Aes192CFBEncrypt aes-192-cfb 加密
func Aes192CFBEncrypt(key [24]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCFBEnc(key[:], plainReader)
}

// Aes192CFBDecrypt aes-192-cfb 解密
func Aes192CFBDecrypt(key [24]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCFBDec(key[:], cipherReader)
}

// Aes256CFBEncrypt aes-256-cfb 加密
func Aes256CFBEncrypt(key [32]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesCFBEnc(key[:], plainReader)
}

// Aes256CFBDecrypt aes-256-cfb 解密
func Aes256CFBDecrypt(key [32]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesCFBDec(key[:], cipherReader)
}

// Aes128OFBEncrypt aes-128-ofb 加密
func Aes128OFBEncrypt(key [16]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesOFBEnc(key[:], plainReader)
}

// Aes128OFBDecrypt aes-128-ofb 解密
func Aes128OFBDecrypt(key [16]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesOFBDec(key[:], cipherReader)
}

// Aes192OFBEncrypt aes-192-ofb 加密
func Aes192OFBEncrypt(key [24]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesOFBEnc(key[:], plainReader)
}

// Aes192OFBDecrypt aes-192-ofb 解密
func Aes192OFBDecrypt(key [24]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesOFBDec(key[:], cipherReader)
}

// Aes256OFBEncrypt aes-256-ofb 加密
func Aes256OFBEncrypt(key [32]byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return aesOFBEnc(key[:], plainReader)
}

// Aes256OFBDecrypt aes-256-ofb 解密
func Aes256OFBDecrypt(key [32]byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return aesOFBDec(key[:], cipherReader)
}

// Convert16bytes 将 []byte 转为 [16]byte
func Convert16bytes(b []byte) (b16 [16]byte) {
	copy(b16[:], b)
	return
}

// Convert24bytes 将 []byte 转为 [24]byte
func Convert24bytes(b []byte) (b24 [24]byte) {
	copy(b24[:], b)
	return
}

// Convert32bytes 将 []byte 转为 [32]byte
func Convert32bytes(b []byte) (b32 [32]byte) {
	copy(b32[:], b)
	return
}

func aesECBEnc(key []byte, plaintext []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}

	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	blockModel := ecb.NewECBEncrypter(block)

	ciphertext = plaintext

	blockModel.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func aesECBDec(key []byte, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockModel := ecb.NewECBDecrypter(block)

	plaintext = ciphertext
	blockModel.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS5UnPadding(plaintext)

	return plaintext, nil
}

func aesCBCEnc(key []byte, plaintext []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}

	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	_, err = cryptorand.Read(iv[:])
	if err != nil {
		return nil, err
	}

	blockModel := cipher.NewCBCEncrypter(block, iv[:])

	blockModel.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func aesCBCDec(key []byte, ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)
	plaintext = ciphertext
	blockModel.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS5UnPadding(plaintext)
	return plaintext, nil
}

func aesCTREnc(key []byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return streamEnc(AesCTR, key, plainReader)
}

func aesCTRDec(key []byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return streamDec(AesCTR, key, cipherReader)
}

func aesCFBEnc(key []byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return streamEnc(AesCFB, key, plainReader)
}

func aesCFBDec(key []byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return streamDec(AesCFB, key, cipherReader)
}

func aesOFBEnc(key []byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	return streamEnc(AesOFB, key, plainReader)
}

func aesOFBDec(key []byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	return streamDec(AesOFB, key, cipherReader)
}

func streamEnc(aesMode AesMode, key []byte, plainReader io.Reader) (cipherReader io.Reader, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 随机初始化向量
	var iv [aes.BlockSize]byte
	cryptorand.Read(iv[:])

	var stream cipher.Stream
	switch aesMode {
	case AesCTR:
		stream = cipher.NewCTR(block, iv[:])
	case AesCFB:
		stream = cipher.NewCFBEncrypter(block, iv[:])
	case AesOFB:
		stream = cipher.NewOFB(block, iv[:])
	default:
		panic("unknown aes mode")
	}

	reader := &cipher.StreamReader{
		S: stream,
		R: plainReader,
	}

	return io.MultiReader(bytes.NewReader(iv[:]), reader), nil
}

func streamDec(aesMode AesMode, key []byte, cipherReader io.Reader) (plainReader io.Reader, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte

	// 读取头部向量
	_, err = cipherReader.Read(iv[:])
	if err != nil {
		return
	}

	var stream cipher.Stream
	switch aesMode {
	case AesCTR:
		stream = cipher.NewCTR(block, iv[:])
	case AesCFB:
		stream = cipher.NewCFBDecrypter(block, iv[:])
	case AesOFB:
		stream = cipher.NewOFB(block, iv[:])
	default:
		panic("unknown aes mode")
	}

	plainReader = &cipher.StreamReader{
		S: stream,
		R: cipherReader,
	}
	return
}

// PKCS5Padding PKCS5 Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	if padding < 0 {
		padding = 0
	}

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding PKCS5 UnPadding
func PKCS5UnPadding(plaintext []byte) []byte {
	length := len(plaintext)
	if length <= 0 {
		return nil
	}

	unpadding := int(plaintext[length-1])
	if length-unpadding < 0 {
		return nil
	}

	return plaintext[:(length - unpadding)]
}
