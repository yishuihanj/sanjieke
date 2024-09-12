package markdown

import (
	"context"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"path"
	"sanjieke/pkg/filenamify"
	"sanjieke/pkg/tool"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

const HtmlExtension = ".html"

func DownloadHtml(ctx context.Context, content, title, dir string, overwrite bool) (bool, error) {
	select {
	case <-ctx.Done():
		return false, context.Canceled
	default:
	}
	defer func() {
		_ = recover()
	}()
	fullName := path.Join(dir, tool.MakeValidFilename(filenamify.Filenamify(title)+HtmlExtension))
	if tool.CheckFileExists(fullName) && !overwrite {
		return true, nil
	}
	hByte := mdToHTML([]byte(content))

	f, err := os.Create(fullName)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return false, err
	}
	// step3: write md file
	_, err = f.Write(hByte)
	if err != nil {
		return false, err
	}
	return false, nil

}
