package config

import "path"

var Config *config

func init() {
	if Config == nil {
		Config = &config{
			VideoQuality: "1080p",
		}
	}
}

type config struct {
	Authorization string `yaml:"authorization"`
	Cookie        string `yaml:"cookie"`
	ApiKey        string `yaml:"sjk-apikey"`
	CourseId      string `yaml:"course_id"`
	Title         string
	OutDirectory  string `yaml:"out_directory"`
	VideoQuality  string
	FfmpegPath    string `yaml:"ffmpeg_path"`
}

// 例如 D://sanjieke//Config.Title
func (c *config) GetCourseRootFolder() string {
	return path.Join(c.OutDirectory, c.Title)
}
