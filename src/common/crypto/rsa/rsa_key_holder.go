package rsaprovider

import (
	"crypto/rsa"
	"go-kit/src/common/configs"
	"strings"
)

type RsaKeyHolder struct {
	rsaMapKey        map[string]*rsa.PublicKey
	rsaMapPrivateKey map[string]*rsa.PrivateKey
}

func (h *RsaKeyHolder) GetKey(code string) (*rsa.PublicKey, bool) {
	key, exists := h.rsaMapKey[strings.ToLower(code)]
	return key, exists
}

func (h *RsaKeyHolder) GetPrivateKey(code string) (*rsa.PrivateKey, bool) {
	key, exists := h.rsaMapPrivateKey[code]
	return key, exists
}

func NewRsaKeyHolder(cf *configs.Config) (*RsaKeyHolder, error) {
	mapPublicRsaKeys := make(map[string]*rsa.PublicKey)
	mapPrivateRsaKeys := make(map[string]*rsa.PrivateKey)

	for key, path := range cf.Rsa.PublicKeys {
		publicKey, err := ParsePublicRsaKey(path)
		if err == nil {
			mapPublicRsaKeys[strings.ToLower(key)] = publicKey
		}
	}

	for key, path := range cf.Rsa.PrivateKeys {
		private, err := ParseRsaKey(path)
		if err == nil {
			mapPrivateRsaKeys[strings.ToLower(key)] = private
		}
	}

	return &RsaKeyHolder{mapPublicRsaKeys, mapPrivateRsaKeys}, nil
}
