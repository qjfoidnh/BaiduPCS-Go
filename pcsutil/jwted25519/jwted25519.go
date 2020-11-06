package jwted25519

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ed25519"
	"unsafe"
)

type (
	SigningMethodEd25519 struct{}
)

var (
	SigningMethodED25519 *SigningMethodEd25519
)

func init() {
	// ED25519
	SigningMethodED25519 = &SigningMethodEd25519{}
	jwt.RegisterSigningMethod(SigningMethodED25519.Alg(), func() jwt.SigningMethod {
		return SigningMethodED25519
	})
}

func (m *SigningMethodEd25519) Alg() string {
	return "ED25519"
}

func (m *SigningMethodEd25519) Sign(signingString string, key interface{}) (string, error) {
	if privkey, ok := key.(ed25519.PrivateKey); ok {
		if len(privkey) != ed25519.PrivateKeySize {
			return "", jwt.ErrInvalidKey
		}

		signature := ed25519.Sign(privkey, *(*[]byte)(unsafe.Pointer(&signingString)))
		return jwt.EncodeSegment(signature), nil
	}

	return "", jwt.ErrInvalidKeyType
}

func (m *SigningMethodEd25519) Verify(signingString, signature string, key interface{}) error {
	pubkey, ok := key.(ed25519.PublicKey)
	if !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(pubkey) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	message, err := jwt.DecodeSegment(signature)
	if err != nil {
		return err
	}

	ok = ed25519.Verify(pubkey, *(*[]byte)(unsafe.Pointer(&signingString)), message)
	if !ok {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
