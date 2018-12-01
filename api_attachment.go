package confluence

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// 上传附件到指定页面
func (cli *Client) UpdateContentAttachments(contentId string, files []string) error {
	//获取页面原有附件的清单
	attachments, err := cli.AttachmentsByContentId(contentId)
	if err != nil {
		return fmt.Errorf("获取页面原有附件清单失败: %s", err)
	}

	attIdByName := make(map[string]string)
	for _, att := range attachments {
		attIdByName[att.Title] = att.Id
	}

	createFiles := make([]string, 0)
	updateFiles := make(map[string]string, 0)
	for _, file := range files {
		basename := filepath.Base(file)
		id, found := attIdByName[basename]
		if found {
			updateFiles[id] = file
		} else {
			createFiles = append(createFiles, file)
		}
	}

	if len(createFiles) > 0 {
		_, err = cli.AttachmentCreate(contentId, createFiles)
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

// 在指定页面创建附件
func (cli *Client) AttachmentCreate(contentId string, fileList []string) ([]Content, error) {
	if len(fileList) <= 0 {
		return nil, fmt.Errorf("file list is empty")
	}

	resp, err := cli.ApiPOSTFiles("/content/"+contentId+"/child/attachment", fileList)
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

// 更新指定页面的附件
func (cli *Client) AttachmentUpdate(contentId, attachmentId, file string) ([]Content, error) {
	resp, err := cli.ApiPOSTFiles("/content/"+contentId+"/child/attachment/"+attachmentId+"/data", []string{file})
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

// 获取指定页面的所有附件
func (cli *Client) AttachmentsByContentId(contentId string) ([]Content, error) {
	resp, err := cli.ApiGET("/content/"+contentId+"/child/attachment", nil)
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
