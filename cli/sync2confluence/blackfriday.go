package main

import (
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"io/ioutil"
	"path"
)

type BlackFridayRenderer struct {
	blackfriday.HTMLRenderer
}

//预处理附件图片
func preRenderImage(w io.Writer, src string) bool {
	dir := path.Dir(src)
	if dir == "." || dir == AssetsDirName {
		basename := path.Base(src)
		result := `<ac:image><ri:attachment ri:filename="` + basename + `" /></ac:image>`
		w.Write([]byte(result))
		return true
	}
	return true
}

//预处理附件
func preRenderLink(w io.Writer, src, title string) bool {
	dir := path.Dir(src)
	if dir == "." || dir == AssetsDirName {
		basename := path.Base(src)
		result := `<ac:link><ri:attachment ri:filename="` + basename + `" /><ac:plain-text-link-body><![CDATA[` + title + `]]></ac:plain-text-link-body></ac:link>`
		w.Write([]byte(result))
		return true
	}
	return true
}

func (r *BlackFridayRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//预处理可能的附件图片
	if node.Type == blackfriday.Image {
		if entering {
			ok := preRenderImage(w, string(node.LinkData.Destination))
			if ok {
				return blackfriday.GoToNext
			}
		}
	}

	//预处理可能的附件
	if node.Type == blackfriday.Link {
		ok := preRenderLink(w, string(node.LinkData.Destination), string(node.LinkData.Title))
		if ok {
			return blackfriday.GoToNext
		}
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
