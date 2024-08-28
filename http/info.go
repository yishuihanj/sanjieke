package http

import (
	"encoding/json"
	"fmt"
	"sanjieke/config"
)

type InfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id               int         `json:"id"`
		UpToCourseId     interface{} `json:"upToCourseId"`
		Type             int         `json:"type"`
		ClassId          int         `json:"classId"`
		ProductId        int         `json:"productId"`
		CosVipFlag       bool        `json:"cosVipFlag"`
		UserIdentity     string      `json:"userIdentity"`
		UpToProductId    interface{} `json:"upToProductId"`
		Title            string      `json:"title"`
		WithService      int         `json:"withService"`
		CompanyId        int         `json:"companyId"`
		CompanyType      string      `json:"companyType"`
		CaasVersion      string      `json:"caasVersion"`
		PickedCourseFlag bool        `json:"pickedCourseFlag"`
		Teachers         []struct {
			UserId           int         `json:"userId"`
			ShowName         string      `json:"showName"`
			ShowAvatar       string      `json:"showAvatar"`
			Avatar           string      `json:"avatar"`
			Title            string      `json:"title"`
			Advantage        string      `json:"advantage"`
			SelfIntro        string      `json:"selfIntro"`
			AuthStatus       string      `json:"authStatus"`
			PublishedBook    interface{} `json:"publishedBook"`
			Recommendation   interface{} `json:"recommendation"`
			GoodAtCategories []string    `json:"goodAtCategories"`
			ContactFlag      interface{} `json:"contactFlag"`
			SocialMediaFlag  interface{} `json:"socialMediaFlag"`
			IpFlag           interface{} `json:"ipFlag"`
			Wxoa             string      `json:"wxoa"`
			Dy               string      `json:"dy"`
			Weibo            string      `json:"weibo"`
			Xiaohongshu      string      `json:"xiaohongshu"`
			Zhihu            string      `json:"zhihu"`
			Github           string      `json:"github"`
			PersonalWeb      string      `json:"personalWeb"`
			Awards           interface{} `json:"awards"`
			CourseSource     interface{} `json:"courseSource"`
		} `json:"teachers"`
		Editors          []interface{} `json:"editors"`
		CompanyInnerFlag bool          `json:"companyInnerFlag"`
		CategoryInfo     struct {
			Direction       string `json:"direction"`
			DirectionId     string `json:"directionId"`
			CategoryId      int    `json:"categoryId"`
			CategoryName    string `json:"categoryName"`
			SubcategoryId   int    `json:"subcategoryId"`
			SubcategoryName string `json:"subcategoryName"`
		} `json:"categoryInfo"`
		Certifications        interface{} `json:"certifications"`
		JoinCertificationFlag bool        `json:"joinCertificationFlag"`
	} `json:"data"`
}

func GetInfo() (*InfoResp, error) {
	// 创建一个新的 HTTP 请求
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/%v/info", config.CourseId)
	err, all := HttpGet(url)
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
