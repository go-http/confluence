package confluence

import (
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"io/ioutil"
	"path"
)

type BlackFridayRenderer struct {
	blackfriday.HTMLRenderer
}

func (r *BlackFridayRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//替换掉以./开头的相对路径为/开头，以适配目前blackfriday的用法
	if len(node.LinkData.Destination) > 0 {
		dest := string(node.LinkData.Destination)
		dir := path.Dir(dest)
		if dir == "." || dir == "assets" {
			node.LinkData.Destination = []byte("/" + path.Base(dest))
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

func parseMarkdownFile(file, absolutePrefix string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := &BlackFridayRenderer{}
	r.Flags = blackfriday.UseXHTML
	r.AbsolutePrefix = absolutePrefix
	mdData := blackfriday.Run(rawData, blackfriday.WithRenderer(r))

	return append([]byte(ConfluenceToc), mdData...), nil
}
