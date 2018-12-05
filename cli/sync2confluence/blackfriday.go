package main

import (
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"io/ioutil"
	"net/url"
	"path"
)

type BlackFridayRenderer struct {
	blackfriday.HTMLRenderer
}

func getAttachmentDir(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return ""
	}

	if u.Scheme != "" {
		return ""
	}

	dir := path.Dir(src)
	if path.IsAbs(src) {
		return ""
	}

	//如果附件位于assets目录，则提取其上级目录
	if path.Base(dir) == AssetsDirName {
		dir = path.Dir(dir)
	}

	return dir
}

func (r *BlackFridayRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//如果是图片，需要检查是否符合图片附件的规则，如果符合，就用Confluence宏而不是HTML输出
	if node.Type == blackfriday.Image {
		dest := string(node.LinkData.Destination)
		dir := getAttachmentDir(dest)
		//预处理图片附件链接
		if dir != "" {
			if entering {
				filename := path.Base(dest)
				result := fmt.Sprintf(`<ac:image><ri:attachment ri:filename="%s">`, filename)

				if dir != "." {
					result += fmt.Sprintf(`<ri:page ri:content-title="%s"/>`, path.Base(dir))
				}

				w.Write([]byte(result))

			} else {
				w.Write([]byte("</ri:attachment></ac:image>"))
			}
			return blackfriday.GoToNext
		}
	}

	//如果是链接，需要检查是否符合附件的规则，如果符合，就用Confluence宏而不是HTML输出
	if node.Type == blackfriday.Link {
		dest := string(node.LinkData.Destination)
		dir := getAttachmentDir(dest)
		//预处理附件链接
		if dir != "" {
			if entering {
				filename := path.Base(dest)
				result := fmt.Sprintf(`<ac:link><ri:attachment ri:filename="%s">`, filename)
				if dir != "." {
					result += fmt.Sprintf(`<ri:page ri:content-title="%s"/>`, path.Base(dir))
				}

				result += "</ri:attachment><ac:plain-text-link-body><![CDATA["
				w.Write([]byte(result))
			} else {
				result := "]]></ac:plain-text-link-body></ac:link>"
				w.Write([]byte(result))
			}
			return blackfriday.GoToNext
		}
	}

	//代码块使用Confluence官方的宏
	//但mermaid代码块除外，因为他需要配合JS渲染成流程图，而不是语法高亮
	if node.Type == blackfriday.CodeBlock && string(node.Info) != "mermaid" {
		result := `<ac:structured-macro ac:name="code">`
		result += `<ac:parameter ac:name="linenumbers">true</ac:parameter>`
		result += `<ac:parameter ac:name="theme">RDark</ac:parameter>`
		result += `<ac:parameter ac:name="language">` + string(node.Info) + `</ac:parameter>`
		result += `<ac:plain-text-body><![CDATA[` + string(node.Literal) + `]]></ac:plain-text-body>`
		result += `</ac:structured-macro>`
		w.Write([]byte(result))
		return blackfriday.GoToNext
	}

	return r.HTMLRenderer.RenderNode(w, node, entering)
}

//Confluence的目录宏，用于自动添加到编译后的页面
const ConfluenceToc = `
<ac:structured-macro ac:name="toc">
	<ac:parameter ac:name="outline">true</ac:parameter>
</ac:structured-macro>
`

func parseMarkdownFile(file string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := &BlackFridayRenderer{}
	r.Flags = blackfriday.UseXHTML

	extensions := blackfriday.CommonExtensions
	if EnableHardLineBreak {
		extensions |= blackfriday.HardLineBreak
	}

	mdData := blackfriday.Run(rawData, blackfriday.WithRenderer(r), blackfriday.WithExtensions(extensions))

	return append([]byte(ConfluenceToc), mdData...), nil
}
