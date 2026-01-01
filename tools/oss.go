package tools

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type OssConfig struct {
	SecretId  string
	SecretKey string
}

var (
	client  *cos.Client
	OssConf string = "config/oss.json"
)

func init() {
	conf := GetOssConf()
	u, _ := url.Parse("https://c8chat-1386757550.cos.ap-hongkong.myqcloud.com")
	su, _ := url.Parse("https://cos.ap-hongkong.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}

	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})
}

func GetOssConf() *OssConfig {
	var oss = &OssConfig{}
	isExist, _ := IsFileExist(OssConf)
	if !isExist {
		return oss
	}
	info, err := ioutil.ReadFile(OssConf)
	if err != nil {
		return oss
	}
	err = json.Unmarshal(info, oss)
	return oss
}

func OssUpload(file *multipart.FileHeader, fname string) error {
	fd, err := file.Open()
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = client.Object.Put(context.Background(), fname, fd, nil)
	if err != nil {
		return err
	}

	return nil
}

func GeOsstUrl(fname string) (string, error) {
	src, err := client.Object.GetPresignedURL2(context.Background(), "GET", fname, time.Minute*5, nil)
	if err != nil {
		return "", err
	}
	return src.String(), nil
}
