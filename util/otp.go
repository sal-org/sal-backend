package util

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

// GenerateOTP - generate otp for a phone number
func GenerateOTP(phone string) (string, bool) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "sal.foundation",
		AccountName: "sal@foundation",
		Secret:      []byte(phone),
	})
	if err != nil {
		fmt.Println("GenerateOTP", err)
		return "", false
	}
	otp, err := totp.GenerateCodeCustom(key.Secret(), time.Now().UTC(), totp.ValidateOpts{Period: 60, Digits: 4, Skew: 1})
	if err != nil {
		fmt.Println("GenerateOTP", err)
		return "", false
	}

	return otp, true
}

// VerifyOTP - verfity given otp for a phone number
func VerifyOTP(phone, otp string) bool {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "sal.foundation",
		AccountName: "sal@foundation",
		Secret:      []byte(phone),
	})
	if err != nil {
		fmt.Println("VerifyOTP", err)
		return false
	}
	valid, err := totp.ValidateCustom(otp, key.Secret(), time.Now().UTC(), totp.ValidateOpts{Period: 60, Digits: 4, Skew: 1})
	if err != nil {
		fmt.Println("VerifyOTP", err)
		return false
	}

	return valid
}
