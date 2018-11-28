package confluence

import (
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"io/ioutil"
	"path"
)

type BlackFridayRenderer struct {
	blackfriday.HTMLRenderer
	ImageSrcPrefix string
}

func (r *BlackFridayRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//替换掉以./开头的本地图片的SRC地址
	if r.ImageSrcPrefix != "" && node.Type == blackfriday.Image {
		imageSrc := string(node.LinkData.Destination)
		if path.Dir(imageSrc) == "." {
			imageSrc = r.ImageSrcPrefix + path.Base(imageSrc)
			node.LinkData.Destination = []byte(imageSrc)
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

func parseMarkdownFile(file, imageSrcPrefix string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := &BlackFridayRenderer{
		ImageSrcPrefix: imageSrcPrefix,
	}
	mdData := blackfriday.Run(rawData, blackfriday.WithRenderer(r))

	return append([]byte(ConfluenceToc), mdData...), nil
}
