package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/athurg/go-confluence"
)

//作为Content解析的文件后缀名，当前支持Markdown文件和直接存储的XML文件
var AssetsDirName = "assets"
var ContentFileExts = []string{".md", ".xml"}
var EnableHardLineBreak bool

func main() {
	var addr, user, pass, space, dir string

	flag.StringVar(&addr, "addr", "https://www.confluence.com", "Confluence访问地址")
	flag.StringVar(&user, "u", "", "用户名")
	flag.StringVar(&pass, "p", "", "密码")
	flag.StringVar(&space, "s", "", "Confluence空间标识")
	flag.StringVar(&dir, "d", "", "要导入的目录")
	flag.StringVar(&AssetsDirName, "assets", "assets", "图片和附件专用的目录名")

	flag.BoolVar(&EnableHardLineBreak, "hardLineBreak", false, "是否将换行符视为页面换行")

	flag.Parse()

	err := ImportToSpace(addr, user, pass, space, dir)
	if err != nil {
		log.Fatal(err)
	}
}

func ImportToSpace(addr, user, pass, space, from string) error {
	client := confluence.New(addr, user, pass)
	_, err := os.Stat(from)
	if err != nil {
		return fmt.Errorf("读取目录失败: %s", err)
	}

	dirs, files, err := getContentInfoLists(from)
	if err != nil {
		return fmt.Errorf("获取列表错误: %s", err)
	}

	//缓存已经创建的Content ID，以便其子Content查找父Content的ID
	contentIds := make(map[string]string)

	//处理目录，由于目录存在依赖关系，因此必需串行处理
	total := len(dirs)
	for i, item := range dirs {
		log.Printf("[%3d/%d]目录: %s", i+1, total, item.Path)
		parentId := contentIds[item.ParentTitle]

		data, err := getDirContentData(item.Path, item.ParentTitle)
		if err != nil {
			return fmt.Errorf("处理目录%s失败: %s", item.Path, err)
		}

		content, err := client.PageFindOrCreateBySpaceAndTitle(space, parentId, item.Title, string(data))
		if err != nil {
			return fmt.Errorf("%s\n创建/更新%s错误: %s", string(data), item.Path, err)
		}

		contentIds[item.Title] = content.Id

		//处理目录下的附件，优先使用其assets子目录，缺省则使用自身
		attachmentFiles, err := getAttachmentFiles(item.Path)
		if err != nil {
			return fmt.Errorf("更新目录%s附件错误: %s", item.Path, err)
		}

		err = client.UpdateContentAttachments(content.Id, attachmentFiles)
		if err != nil {
			return fmt.Errorf("更新目录%s附件错误: %s", item.Path, err)
		}
	}

	//并发处理文件
	var wg sync.WaitGroup
	errs := make(chan error, len(files))
	for _, item := range files {
		wg.Add(1)

		go func(info FileContentInfo) {
			defer wg.Done()

			parentId := contentIds[info.ParentTitle]

			buff, err := getFileContentData(info.Path, info.Ext, info.ParentTitle)
			if err != nil {
				errs <- fmt.Errorf("\033[31m错误文件%s: %s\033[0m", info.Path, err)
				return
			}

			_, err = client.PageFindOrCreateBySpaceAndTitle(space, parentId, info.Title, string(buff))
			if err != nil {
				errs <- fmt.Errorf("%s\n\033[31m错误文件%s错误: %s\033[0m", string(buff), info.Path, err)
				return
			}

			log.Printf("完成文件: %s", info.Path)
		}(item)
	}

	wg.Wait()

	//输出所有的错误信息
	for len(errs) > 0 {
		log.Println("", <-errs, "\033[0m")
	}

	return nil
}

// 从指定目录获取有效Conten列表
func getContentInfoLists(rootPath string) ([]FileContentInfo, []FileContentInfo, error) {
	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, nil, fmt.Errorf("获取%s的绝对路径失败: %s", rootPath, err)
	}

	dirs := make([]FileContentInfo, 0)
	files := make([]FileContentInfo, 0)
	titles := make(map[string][]string)

	err = filepath.Walk(absRootPath, func(path string, info os.FileInfo, err error) error {
		//遍历的filepath和rootPath取相对路径肯定是始终成功的
		relPath, _ := filepath.Rel(absRootPath, path)

		//顶层目录自身不需处理
		if relPath == "." {
			return nil
		}

		contentInfo := GetFileContentInfo(relPath, info.IsDir())
		contentInfo.Path = path

		title := contentInfo.Title

		//以.开头的文件跳过、以.开头的目录及其子目录跳过
		if title == "" || strings.HasPrefix(title, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		//目录直接处理
		if info.IsDir() {
			//assets目录跳过
			if info.Name() == AssetsDirName {
				return filepath.SkipDir
			}

			dirs = append(dirs, contentInfo)

			if _, found := titles[title]; !found {
				titles[title] = make([]string, 0, 1)
			}
			titles[title] = append(titles[title], path)

			return nil
		}

		//只支持普通文件，不支持符号链接、设备等其他类型的文件
		if !info.Mode().IsRegular() {
			return fmt.Errorf("文件%s不是普通文件", path)
		}

		//目前只支持md、xml格式的文件
		var isExtSupport bool
		for _, ext := range ContentFileExts {
			if ext == contentInfo.Ext {
				isExtSupport = true
			}
		}
		if !isExtSupport {
			return nil
		}

		//索引文件会在目录列表处理时读取，文件列表直接忽略
		if title == "index" {
			return nil
		}

		files = append(files, contentInfo)

		if _, found := titles[title]; !found {
			titles[title] = make([]string, 0, 1)
		}
		titles[title] = append(titles[title], path)

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("遍历目录%s错误: %s", rootPath, err)
	}

	var duplicatedCount int
	for title, matches := range titles {
		if len(matches) == 1 {
			continue
		}
		duplicatedCount += 1
		log.Println(title, "x", len(matches))
		for _, match := range matches {
			log.Println("\t", match)
		}
	}

	if duplicatedCount > 0 {
		return nil, nil, fmt.Errorf("有%d个重复的标题", duplicatedCount)
	}

	return dirs, files, nil
}

type FileContentInfo struct {
	Path        string
	Title       string
	ParentTitle string
	Ext         string
}

// 获取指定文件的信息
func GetFileContentInfo(path string, isDir bool) FileContentInfo {
	var ext string
	title := filepath.Base(path)

	// 目录不需要分离后缀名，非目录则需要
	if !isDir {
		ext = filepath.Ext(title)
		title = strings.TrimSuffix(title, ext)
	}

	parentTitle := filepath.Base(filepath.Dir(path))
	if parentTitle == "." {
		parentTitle = ""
	}

	return FileContentInfo{
		Ext:         ext,
		Title:       title,
		ParentTitle: parentTitle,
	}
}

var DefaultDirContentData = []byte(`<ac:structured-macro ac:name="children"><ac:parameter ac:name="all">true</ac:parameter></ac:structured-macro>`)

func getDirContentData(dir, parentTitle string) ([]byte, error) {
	//检查是否有索引文件，如果有则用索引替换掉缺省的标准模板
	for _, ext := range ContentFileExts {
		indexFile := filepath.Join(dir, "index"+ext)
		buff, err := getFileContentData(indexFile, ext, parentTitle)
		if err == nil {
			return buff, nil
		}

		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("读取失败: %s", err)
		}
	}

	//如果没有找到任何合适的索引文件，就用缺省模板
	return DefaultDirContentData, nil
}

func getFileContentData(file, ext, parentTitle string) ([]byte, error) {
	if ext == ".xml" {
		return ioutil.ReadFile(file)
	}

	if ext == ".md" {
		return parseMarkdownFile(file, parentTitle)
	}

	return nil, fmt.Errorf("不支持的文件格式: %s", ext)
}

// 获取指定目录下的附件清单
func getAttachmentFiles(dir string) ([]string, error) {
	assetsDir := filepath.Join(dir, AssetsDirName)
	if fileInfo, err := os.Stat(assetsDir); !os.IsNotExist(err) && fileInfo.IsDir() {
		dir = assetsDir
	}

	//读取目录，获取需要添加、更新的附件清单
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录错误: %s", err)
	}

	attachmentFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !isAttachmentFilename(file.Name()) {
			continue
		}

		attachmentFiles = append(attachmentFiles, filepath.Join(dir, file.Name()))
	}

	return attachmentFiles, nil
}

func isAttachmentFilename(filename string) bool {
	//去掉路径名，获取干净的文件名
	fname := filepath.Base(filename)

	//点开头的文件不算附件
	if fname[0] == '.' {
		return false
	}

	//被解析器支持的内容文件不算附件
	ext := filepath.Ext(filename)
	for _, e := range ContentFileExts {
		if e == ext {
			return false
		}
	}

	return true
}
