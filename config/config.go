package config

var Authorization string
var Cookie string
var ApiKey string
var StudyId string
var Title string
var CourseIdMap map[int]*Course
var CourseNameMap map[string]*Course

// 课程数据
type Course struct {
	Name       string //名字
	Id         int    //课程id
	VideoList  []CourseVideo
	NextCourse *Course
}

type CourseVideo struct {
	Url             string //地址
	ResolutionRatio string // 1080p 720p
}

func ClearCache() {
	Authorization = ""
	Cookie = ""
	ApiKey = ""
	StudyId = ""
	CourseIdMap = make(map[int]*Course)
	CourseNameMap = make(map[string]*Course)
	Title = ""
}

func CheckConfigValid() bool {
	if Authorization == "" || Cookie == "" || ApiKey == "" || StudyId == "" {
		return false
	}
	return true
}
