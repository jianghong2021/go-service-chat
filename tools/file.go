package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// 判断文件文件夹是否存在(字节0也算不存在)
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

// 判断文件文件夹不存在
func IsFileNotExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

func ImageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	var mimeType string

	// 尝试编码为PNG
	err := png.Encode(&buf, img)
	if err == nil {
		mimeType = "image/png"
	} else {
		// 如果PNG失败，尝试JPEG
		buf.Reset()
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		if err != nil {
			return "", fmt.Errorf("图片编码失败: %v", err)
		}
		mimeType = "image/jpeg"
	}

	// 生成Base64字符串
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str), nil
}
