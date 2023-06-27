package googleauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// GetSecret 生成秘钥信息
func GetSecret() string {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, time.Now().UnixNano()/1000/30)
	return strings.ToUpper(base32.StdEncoding.EncodeToString(hmacSha1(buf.Bytes(), nil)))
}

// GetQrcodeUrl 获取二维码
func GetQrcodeUrl(user, secret string) string {
	qrcode := fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
}

// VerifyCode 验证动态码
func VerifyCode(secret, code string) bool {
	return getCode(secret) == code
}

// GetCode 获取动态验证码
func GetCode(secret string) string {
	return getCode(secret)
}

func hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if len(data) > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

// getCode
func getCode(secret string) string {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}
	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	return fmt.Sprintf("%06d", oneTimePassword(key, toBytes(epochSeconds/30)))
}
func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}
func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)
	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F
	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]
	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F
	number := toUint32(hashParts)
	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000
	return pwd
}
