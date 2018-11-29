package confluence

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
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
	Link      LinkResp    `json:"_links"`
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

func (cli *Client) ContentBySpaceAndTitle(space, title string) (Content, error) {
	q := url.Values{
		"title":    {title},
		"spaceKey": {space},
		"expand":   {"version"},
	}

	resp, err := cli.GET("/content", q)
	if err != nil {
		return Content{}, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		PageResp
		Results []Content
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Content{}, fmt.Errorf("解析响应失败: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Content{}, fmt.Errorf("[%d]%s", info.StatusCode, info.Message)
	}

	switch info.Size {
	case 0:
		return Content{}, nil
	case 1:
		return info.Results[0], nil
	default:
		return Content{}, fmt.Errorf("找到%d条记录", info.Size)
	}
}

func (cli *Client) PageCreateInSpace(space, parentId, title, data string) (Content, error) {
	return cli.ContentCreateInSpace("page", space, parentId, title, data)
}

func (cli *Client) ContentCreateInSpace(contentType, space, parentId, title, data string) (Content, error) {
	content := Content{Type: contentType, Title: title}
	content.Space.Key = space
	content.Body.Storage.Value = data
	content.Body.Storage.Representation = "storage"

	//FIXME: 这里指定了创建信息，但是好像没什么用
	content.Version.Message = time.Now().Local().Format("机器人创建于2006-01-02 15:04:05")

	//设置父页面
	if parentId != "" {
		content.Ancestors = []Content{Content{Id: parentId}}
	}

	resp, err := cli.POST("/content", content)
	if err != nil {
		return Content{}, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Content{}, fmt.Errorf("解析响应失败: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Content{}, fmt.Errorf("[%d]%s", resp.StatusCode, info.Message)
	}

	return info.Content, nil
}

func (cli *Client) ContentUpdate(content Content) (Content, error) {
	resp, err := cli.PUT("/content/"+content.Id, content)
	if err != nil {
		return Content{}, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		Content
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Content{}, fmt.Errorf("解析响应失败: %s", err)
	}

	if info.StatusCode != 0 {
		fmt.Printf("%#v\n", info)
		return Content{}, fmt.Errorf("[%d]%s", info.StatusCode, info.Message)
	}

	return info.Content, nil
}

//从指定空间查找或创建指定标题的Content
func (cli *Client) PageFindOrCreateBySpaceAndTitle(space, parentId, title, data string) (Content, error) {
	content, err := cli.ContentBySpaceAndTitle(space, title)
	if err != nil {
		return Content{}, fmt.Errorf("查找%s出错: %s", title, err)
	}

	// 不存在则创建
	if content.Id == "" {
		return cli.PageCreateInSpace(space, parentId, title, data)
	}

	// 存在则否则更新
	content.Space.Key = space
	content.Version.Number += 1
	content.Version.Message = time.Now().Local().Format("机器人更新于2006-01-02 15:04:05")
	content.Body.Storage.Value = data
	content.Body.Storage.Representation = "storage"

	//设置父页面
	if parentId != "" {
		content.Ancestors = []Content{Content{Id: parentId}}
	}

	return cli.ContentUpdate(content)
}
