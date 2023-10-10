package dhttp

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

// 获取api的unix时间戳
func GetApiUnixTime() int64 {
	client := &http.Client{
		Timeout: 10 * time.Second, //请求超时时间
	}

	url := "https://time.is/t/?zh.0.347.2464.0p.480.43d.1574683214663.1594044507830."

	// 创建一个Get请求
	req, err := http.NewRequest("GET", url, nil)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return time.Now().Unix()
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Now().Unix()
	}
	str := string(body)[:10]
	expTime, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Now().Unix()
	}

	return expTime
}
