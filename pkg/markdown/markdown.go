package markdown

//引用于 github.com/nicoxiang/geektime-downloader

import (
	"context"
	"errors"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sanjieke/pkg/downloader"
	"sanjieke/pkg/filenamify"
	"sanjieke/pkg/tool"
	"strings"
	"sync"
)

var (
	converter *md.Converter
	imgRegexp = regexp.MustCompile(`!\[(.*?)]\((.*?)\)`)
)

// MDExtension ...
const MDExtension = ".md"

type markdownString struct {
	sync.Mutex
	s string
}

func (ms *markdownString) ReplaceAll(o, n string) {
	ms.Lock()
	defer ms.Unlock()
	ms.s = strings.ReplaceAll(ms.s, o, n)
}

// Download article as markdown
func Download(ctx context.Context, content, title, dir string, overwrite bool) (bool, error, string) {
	select {
	case <-ctx.Done():
		return false, context.Canceled, ""
	default:
	}
	fullName := path.Join(dir, tool.MakeValidFilename(filenamify.Filenamify(title)+MDExtension))
	if tool.CheckFileExists(fullName) && !overwrite {
		return true, nil, ""
	}

	// step1: convert to md string
	// 添加自定义转换规则来处理 video 标签
	getDefaultConverter().AddRules(md.Rule{
		Filter: []string{"video"},
		Replacement: func(content string, sel *goquery.Selection, options *md.Options) *string {
			// 获取 video 标签的原始 HTML 内容
			var buf strings.Builder
			sel.Each(func(i int, s *goquery.Selection) {
				// 开始构建 <video> 标签
				buf.WriteString("<video")
				// 添加 video 标签的属性
				for _, attr := range s.Nodes[0].Attr {
					buf.WriteString(" " + attr.Key + "=\"" + attr.Val + "\"")
				}
				buf.WriteString(">")
				// 添加标签内的 HTML 内容
				htmlContent, _ := s.Html()
				buf.WriteString(htmlContent)
				buf.WriteString("</video>")
			})

			result := buf.String()
			return &result
		},
	})
	markdown, err := getDefaultConverter().ConvertString(content)
	if err != nil {
		return false, err, ""
	}

	// step2: download images
	var ss = &markdownString{s: markdown}
	imageURLs := findAllImages(markdown)

	// images/aid/imageName.png
	imagesFolder := filepath.Join(dir, "图片")

	if _, err := os.Stat(imagesFolder); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(imagesFolder, os.ModePerm)
	}

	err = writeImageFile(ctx, imageURLs, dir, imagesFolder, ss)

	if err != nil {
		return false, err, ""
	}

	f, err := os.Create(fullName)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return false, err, ""
	}
	// step3: write md file
	_, err = f.WriteString("# " + title + "\n" + ss.s)
	if err != nil {
		return false, err, ""
	}
	return false, nil, ss.s
}

func findAllImages(md string) (images []string) {
	for _, matches := range imgRegexp.FindAllStringSubmatch(md, -1) {
		if len(matches) == 3 {
			s := matches[2]
			_, err := url.ParseRequestURI(s)
			if err == nil {
				images = append(images, s)
			}
			// sometime exists broken image url, just ignore
		}
	}
	return
}

func getDefaultConverter() *md.Converter {
	if converter == nil {
		converter = md.NewConverter("", true, nil)
	}
	return converter
}

func writeImageFile(ctx context.Context,
	imageURLs []string,
	dir,
	imagesFolder string,
	ms *markdownString,
) (err error) {
	for _, imageURL := range imageURLs {
		segments := strings.Split(imageURL, "/")
		f := segments[len(segments)-1]
		if i := strings.Index(f, "?"); i > 0 {
			f = f[:i]
		}
		imageLocalFullPath := filepath.Join(imagesFolder, f)

		headers := make(map[string]string, 0)
		_, err = downloader.DownloadFileConcurrently(ctx, imageLocalFullPath, imageURL, headers, 1)

		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(dir, imageLocalFullPath)
		ms.ReplaceAll(imageURL, filepath.ToSlash(rel))
	}
	return nil
}
