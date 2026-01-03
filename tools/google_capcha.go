package tools

import (
	"encoding/json"
	"log"
	"net/url"
	"time"
)

var (
	GoogleCaptchaUrl = "https://www.google.com/recaptcha/api/siteverify"
	GoogleCaptchaKey = ""
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func init() {
	LoadAppConf()
}

func SiteverifyWithLogin(token string) (bool, error) {
	data := url.Values{}
	data.Set("secret", AppConf.GoogleCaptcha.SecretKey)
	data.Set("response", token)
	res, err := PostForm(GoogleCaptchaUrl, data)
	if err != nil {
		return false, err
	}
	var result = &RecaptchaResponse{}
	err = json.Unmarshal([]byte(res), result)

	if err != nil {
		return false, err
	}
	ok := result.Success && result.Score >= 0.5 && result.Action == "login"
	if !ok {
		log.Println("google验证失败:", result.ErrorCodes)
	}
	return ok, nil
}
