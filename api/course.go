package api

import (
	"encoding/json"
	"fmt"
	"sanjieke/config"
	"sanjieke/pkg/httper"
)

type CourseResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *CourseData `json:"data"`
}

type CourseData struct {
	Nodes []*CourseNodes `json:"nodes"`
}

type CourseNodes struct {
	ContentType  string        `json:"contentType"`  // b-video 为视频 html为文本
	HTMLContent  any           `json:"htmlContent"`  // 如果是html，则有此字段，
	VideoContent *VideoContent `json:"videoContent"` //如果是video，则有此字段，
}

// 视频集合
type ResolutionRatioObj struct {
	URL             string `json:"url"`             // 下载地址
	ResolutionRatio string `json:"resolutionRatio"` // 1080p 720p
	MediaType       string `json:"mediaType"`       // M3U8
}

type VideoContent struct {
	ResolutionRatioObjList []*ResolutionRatioObj `json:"resolutionRatioObjList"`
}

func GetCourseNode(nodeId int) (*CourseResp, error) {
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/%v/content/%v", config.Config.CourseId, nodeId)
	err, body := httper.HttpGet(url)
	if err != nil {
		return nil, err
	}
	resp := new(CourseResp)

	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
