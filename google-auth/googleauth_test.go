package googleauth

import (
	"fmt"
	"testing"
)

func TestGoogleAuthSecret(t *testing.T) {
	secret := GetSecret()
	fmt.Println(secret)
	fmt.Println(GetQrcodeUrl("todo", secret))
}

func TestGoogleAuthVerifyCode(t *testing.T) {
	fmt.Println(VerifyCode("todo", "925604"))
}
