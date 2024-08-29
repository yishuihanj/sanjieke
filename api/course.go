package http

import (
	"encoding/json"
	"fmt"
)

type CourseNode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Nodes []struct {
			ContentType  string      `json:"contentType"`
			ContentId    interface{} `json:"contentId"`
			HtmlContent  interface{} `json:"htmlContent"`
			VideoContent struct {
				ResolutionRatioObjList []struct {
					Url             string `json:"url"`
					Name            string `json:"name"`
					VipFlag         bool   `json:"vipFlag"`
					ResolutionRatio string `json:"resolutionRatio"`
					MediaType       string `json:"mediaType"`
				} `json:"resolutionRatioObjList"`
				Subtitles []interface{} `json:"subtitles"`
				ContentId int           `json:"contentId"`
				Title     string        `json:"title"`
				Duration  int           `json:"duration"`
				Width     int           `json:"width"`
				Height    int           `json:"height"`
				Url       struct {
					P  string `json:"360p"`
					P1 string `json:"608p"`
				} `json:"url"`
			} `json:"videoContent"`
			ProgramQuestionContent interface{} `json:"programQuestionContent"`
			QuestionContent        interface{} `json:"questionContent"`
		} `json:"nodes"`
		AiFlag bool `json:"aiFlag"`
	} `json:"data"`
}

func GetCourseNode(nodeId int) (*CourseNode, error) {
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/34003473/content/%v", nodeId)
	err, body := HttpGet(url)
	if err != nil {
		return nil, err
	}
	resp := new(CourseNode)

	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
