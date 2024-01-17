package rsaprovider

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"os"
)

func VerifySignature(publicKey *rsa.PublicKey, signed, plain string) error {
	message := []byte(plain)
	hashed := sha1.Sum(message)
	signatureByte, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], signatureByte)
	if err != nil {
		return err
	}
	return nil
}

func ParseRsaKey(path string) (*rsa.PrivateKey, error) {
	rawPrivateKey, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	privatePem, _ := pem.Decode(rawPrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func ParsePublicRsaKey(path string) (*rsa.PublicKey, error) {
	publicKeyPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyPem)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return publicKey, nil
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func GenerateSignatureByPrivateKey(priv *rsa.PrivateKey, mes string) (*string, error) {
	rawMes, err := base64.StdEncoding.DecodeString(mes)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(rawMes)
	signature, errG := rsa.SignPKCS1v15(nil, priv, crypto.SHA256, hashed[:])
	if errG != nil {
		return nil, errG
	}
	signatureStr := hex.EncodeToString(signature)
	return &signatureStr, nil
}
