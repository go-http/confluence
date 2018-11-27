package confluence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

// 将指定目录的附件全部上传
func (cli *Client) AttachmentCreateFromDir(contentId string, dir string) error {
	//获取页面原有附件的清单
	attachments, err := cli.AttachmentByContentId(contentId)
	if err != nil {
		return fmt.Errorf("获取页面原有附件清单失败: %s", err)
	}

	attIdByName := make(map[string]string)
	for _, att := range attachments {
		attIdByName[att.Title] = att.Id
	}

	//读取目录，获取需要添加、更新的附件清单
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录错误: %s", err)
	}

	updateFiles := make(map[string]string, 0)
	newFileList := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !isAttachmentFilename(file.Name()) {
			continue
		}

		id, found := attIdByName[file.Name()]
		if found {
			updateFiles[id] = path.Join(dir, file.Name())
		} else {
			newFileList = append(newFileList, path.Join(dir, file.Name()))
		}
	}

	if len(newFileList) > 0 {
		_, err = cli.AttachmentCreate(contentId, newFileList)
		if err != nil {
			return fmt.Errorf("添加新附件错误: %s", err)
		}
	}

	for attchmentId, file := range updateFiles {
		_, err = cli.AttachmentUpdate(contentId, attchmentId, file)
		if err != nil {
			return fmt.Errorf("更新附件%s错误: %s", file, err)
		}
	}

	return nil
}

// 创建附件
func (cli *Client) AttachmentCreate(contentId string, fileList []string) ([]Content, error) {
	if len(fileList) <= 0 {
		return nil, fmt.Errorf("file list is empty")
	}

	resp, err := cli.POSTFiles("/content/"+contentId+"/child/attachment", fileList)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
		Results []Content
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%d]%s", resp.StatusCode, info.Message)
	}

	return info.Results, nil
}

// 创建附件
func (cli *Client) AttachmentUpdate(contentId, attachmentId, file string) ([]Content, error) {
	resp, err := cli.POSTFiles("/content/"+contentId+"/child/attachment/"+attachmentId+"/data", []string{file})
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
		Results []Content
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%d]%s", resp.StatusCode, info.Message)
	}

	return info.Results, nil
}

// 获取指定页面的附件
func (cli *Client) AttachmentByContentId(contentId string) ([]Content, error) {
	resp, err := cli.GET("/content/"+contentId+"/child/attachment", nil)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
		Results []Content
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%d]%s", resp.StatusCode, info.Message)
	}

	return info.Results, nil
}

func isAttachmentFilename(filename string) bool {
	//去掉路径名，获取干净的文件名
	fname := path.Base(filename)

	//点开头的文件不算附件
	if fname[0] == '.' {
		return false
	}

	//被解析器支持的内容文件不算附件
	ext := path.Ext(filename)
	for _, supportedExt := range supportedFileExts {
		if supportedExt == ext {
			return false
		}
	}

	return true
}
