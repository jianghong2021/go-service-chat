package tools

import (
	"image"

	"github.com/pquerna/otp/totp"
)

func init() {
	LoadAppConf()
}

func GenerateOpts(userName string) (string, image.Image, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "C8 IM ",
		AccountName: userName,
	})

	if err != nil {
		return "", nil, err
	}
	qrcode, err := key.Image(300, 300)
	if err != nil {
		return "", nil, err
	}
	return key.Secret(), qrcode, nil
}

func ValidateOtps(passcode string, secret string) bool {
	ok := totp.Validate(passcode, secret)
	return ok
}

func DecodeOtpsKey(str string) (string, error) {
	return AesDecrypt(str, AppConf.OtpsKey)
}

func EncodeOtpsKey(str string) (string, error) {
	return AesEncrypt(str, AppConf.OtpsKey)
}
