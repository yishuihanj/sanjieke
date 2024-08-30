package markdown

import (
	"encoding/json"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"log"
	"os"
)

type Node struct {
	ContentType            string  `json:"contentType"`
	ContentId              *int    `json:"contentId"`
	HtmlContent            string  `json:"htmlContent"`
	VideoContent           *string `json:"videoContent"`
	ProgramQuestionContent *string `json:"programQuestionContent"`
	QuestionContent        *string `json:"questionContent"`
}

type Resp struct {
	Data *Data `json:"data"`
}

type Data struct {
	Nodes []Node `json:"nodes"`
}

func htm(html string) string {
	opt := &md.Options{
		EscapeMode: "basic", // default
	}
	converter := md.NewConverter("", true, opt)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	return markdown + "\n\n"
}

func main() {
	resp := new(Resp)
	err := json.Unmarshal([]byte(str), resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建Markdown文件
	file, err := os.Create("output.md")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, node := range resp.Data.Nodes {
		if node.ContentType == "html" {
			//markdown, err := htmlToMarkdown(node.HtmlContent)
			//if err != nil {
			//	fmt.Println("Error converting HTML to Markdown:", err)
			//	continue
			//}

			markdown := htm(node.HtmlContent)

			// 将Markdown内容写入文件
			_, err = file.WriteString(markdown)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}

	fmt.Println("Markdown content successfully written to output.md")
}
