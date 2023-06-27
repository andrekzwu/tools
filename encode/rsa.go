package encode

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RsaSign
func RsaSign(content, privateKeyStr, signType string) (string, error) {
	encodedKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey([]byte(encodedKey))
	if err != nil {
		return "", err
	}
	var hashed []byte
	var hash crypto.Hash
	switch signType {
	case "RSA":
		rsaHashed := sha1.Sum([]byte(content))
		hashed = rsaHashed[:]
		hash = crypto.SHA1
	case "RSA2":
		rsa2Hashed := sha256.Sum256([]byte(content))
		hashed = rsa2Hashed[:]
		hash = crypto.SHA256
	}
	s, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), hash, hashed)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return data, nil

}

// RsaSignCheck
func RsaSignCheck(content, sign, publicKeyStr, signType string) (bool, error) {
	encodedKey, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return false, err
	}
	publicKey, err := x509.ParsePKIXPublicKey(encodedKey)
	if err != nil {
		return false, err
	}
	var hashed []byte
	var hash crypto.Hash
	switch signType {
	case "RSA":
		rsaHashed := sha1.Sum([]byte(content))
		hashed = rsaHashed[:]
		hash = crypto.SHA1
	case "RSA2":
		rsa2Hashed := sha256.Sum256([]byte(content))
		hashed = rsa2Hashed[:]
		hash = crypto.SHA256
	}
	data, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}
	if err := rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), hash, hashed, data); err != nil {
		return false, err
	}
	return true, nil
}

// RsaEncrypt
func RsaEncrypt(content, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, content)
}

// RsaDecrypt
func RsaDecrypt(content, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, content)
}
