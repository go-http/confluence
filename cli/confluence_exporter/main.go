package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/athurg/go-confluence"
)

func main() {
	var addr, user, pass, space, dir string

	flag.StringVar(&addr, "addr", "https://www.confluence.com", "Confluence访问地址")
	flag.StringVar(&user, "u", "", "用户名")
	flag.StringVar(&pass, "p", "", "密码")
	flag.StringVar(&space, "s", "", "Confluence空间标识")
	flag.StringVar(&dir, "d", "", "要导出的目录")

	flag.Parse()

	err := exportSpaceTo(addr, user, pass, space, dir)
	if err != nil {
		log.Fatal(err)
	}
}

func exportSpaceTo(addr, user, pass, space, outDir string) error {
	client := confluence.New(addr, user, pass)

	pages, err := client.GetAllSpacePages(space)
	if err != nil {
		return err
	}

	//清空原目录
	os.RemoveAll(outDir)
	os.MkdirAll(outDir, 0755)

	total := len(pages)
	for i, page := range pages {

		//获取父页面信息
		pageDirs := []string{outDir}
		for _, ancestor := range page.Ancestors {
			pageDirs = append(pageDirs, ancestor.Title)
		}

		//创建所需目录
		pageDir := path.Join(pageDirs...)
		os.MkdirAll(pageDir, 0755)

		//将目录同名的文件挪为index.xml
		os.Rename(pageDir+".xml", path.Join(pageDir, "index.xml"))

		//输出文件
		file := path.Join(pageDir, page.Title+".xml")
		ioutil.WriteFile(file, []byte(page.Body.Storage.Value), 0755)

		fileKBSize := float32(len(page.Body.Storage.Value)) / 100
		log.Printf("[%3d/%3d] (%7.2f KiB) %s", i+1, total, fileKBSize, page.Title)

		//下载附件
		attachmentDir := path.Join(pageDir, page.Title)
		os.MkdirAll(attachmentDir, 0755)
		os.Rename(attachmentDir+".xml", path.Join(attachmentDir, "index.xml")) //将目录同名文件挪为index.xml

		attachments, err := client.AttachmentByContentId(page.Id)
		if err != nil {
			return fmt.Errorf("获取%s附件列表错误: %s", page.Title, err)
		}
		for j, att := range attachments {
			attachmentData, err := client.AttachmentDownload(att.Link.Download)
			if err != nil {
				return fmt.Errorf("下载%s附件%s错误: %s", page.Title, att.Title, err)
			}

			attachmentFile := path.Join(attachmentDir, att.Title)
			ioutil.WriteFile(attachmentFile, attachmentData, 0755)
			log.Printf("         附件%2d %s => %s", j+1, att.Title, attachmentFile)
		}
	}

	return nil
}
