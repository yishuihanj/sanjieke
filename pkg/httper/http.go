package http

import (
	"io"
	"net/http"
	"sanjieke/config"
	"time"
)

// HttpGet
func HttpGet(url string) (error, []byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil
	}
	setNormalHeader(req.Header)
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	return nil, all
}

func Get(url string) (error, io.ReadCloser) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil
	}
	setNormalHeader(req.Header)
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	return nil, resp.Body
}

func setNormalHeader(headers http.Header) {
	// 设置请求头
	headers.Set("accept", "application/json, text/plain, */*")
	headers.Set("accept-language", "zh-CN,zh;q=0.9")
	headers.Set("Authorization", config.Authorization)
	headers.Set("Cookie", config.Cookie)
	headers.Set("origin", "https://study.sanjieke.cn")
	headers.Set("priority", "u=1, i")
	headers.Set("referer", "https://study.sanjieke.cn/")
	headers.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	headers.Set("sec-ch-ua-mobile", "?0")
	headers.Set("sec-ch-ua-platform", `"Windows"`)
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("sec-fetch-mode", "cors")
	headers.Set("sec-fetch-site", "same-site")
	headers.Set("sjk-apikey", config.ApiKey)
	headers.Set("sjk-platform", "pc")
	headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	headers.Set("x-domain-prefix", "cos")
	headers.Set("x-requested-with", "XMLHttpRequest")
}
