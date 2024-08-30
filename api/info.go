package api

import (
	"encoding/json"
	"fmt"
	"sanjieke/config"
	"sanjieke/pkg/httper"
)

type InfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id        int    `json:"id"`
		Type      int    `json:"type"`
		ClassId   int    `json:"classId"`
		ProductId int    `json:"productId"`
		Title     string `json:"title"`
	} `json:"data"`
}

func GetInfo() (*InfoResp, error) {
	// 创建一个新的 HTTP 请求
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/%v/info", config.Config.CourseId)
	err, all := httper.HttpGet(url)
	if err != nil {
		return nil, err
	}
	infoResp := new(InfoResp)
	err = json.Unmarshal(all, infoResp)
	if err != nil {
		return nil, err
	}
	return infoResp, nil
}
