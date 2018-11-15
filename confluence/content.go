package confluence

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ContentBody struct {
	Storage             Storage     `json:"storage,omitempty"`
	Editor              interface{} `json:"editor,omitempty"`
	View                interface{} `json:"view,omitempty"`
	ExportView          interface{} `json:"export_view,omitempty"`
	StyledView          interface{} `json:"styled_view,omitempty"`
	AnonymousExportView interface{} `json:"anonymous_export_view,omitempty"`
}

type Content struct {
	Id        string      `json:"id,omitempty"`
	Type      string      `json:"type,omitempty"`
	Title     string      `json:"title,omitempty"`
	Space     Space       `json:"space,omitempty"`
	Body      ContentBody `json:"body,omitempty"`
	Version   Version     `json:"version,omitempty"`
	Ancestors []Content   `json:"ancestors,omitempty"`
}

// 通过ID获取内容
func (cli *Client) ContentById(id string) (Content, error) {
	return cli.ContentByIdWithOpt(id, nil)
}

// 通过ID获取内容（可以设置获取选项）
func (cli *Client) ContentByIdWithOpt(id string, opt url.Values) (Content, error) {
	if opt == nil {
		opt = url.Values{}
	}

	// 缺省情况下，需要展开version，以便于后期编辑
	if opt.Get("expand") == "" {
		opt.Set("expand", "version")
	}

	resp, err := cli.GET("/content/"+id, opt)
	if err != nil {
		return Content{}, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
		Results []Content
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Content{}, fmt.Errorf("解析响应失败: %s", err)
	}

	if info.StatusCode != 0 {
		return Content{}, fmt.Errorf("[%d]%s", info.StatusCode, info.Message)
	}

	return info.Content, nil
}
