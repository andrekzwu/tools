package encode

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// HmacSha256
func HmacSha256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
