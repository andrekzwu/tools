package encode

import (
	"encoding/base64"
)

/*
*
base64
*/
func Base64Encrypt(origData []byte) string {
	result := base64.StdEncoding.EncodeToString(origData)
	return result
}
func Base64Decrypt(ciphertext string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	return decodeBytes, err
}
