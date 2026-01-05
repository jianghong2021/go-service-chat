package tools

import (
	"encoding/json"
	"errors"
	"log"
)

type C8AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data C8AuthUserInfo
}

type C8AuthUserInfo struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
}

func init() {
	LoadAppConf()
}

func C8AuthAndGetInfo(token string) (C8AuthUserInfo, error) {
	var result = &C8AuthResponse{}

	data := map[string]string{
		"token": token,
	}

	heaher := map[string]string{
		"platform": "10",
		"language": "zh-cn",
		"timezone": "shanghai",
	}
	res, err := PostFormWithHeaders(AppConf.C8api+"/client/member/member-info/checkToken", data, heaher)
	if err != nil {
		return result.Data, err
	}

	err = json.Unmarshal([]byte(res), result)

	if err != nil {
		return result.Data, err
	}

	if result.Code != 200 {
		log.Println("c8认证失败:", result.Msg)
		return result.Data, errors.New(result.Msg)
	}
	return result.Data, nil
}
