package confluence

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//根据SpaceKey获取空间的信息
func (cli *Client) SpaceByKey(key string) (Space, error) {
	resp, err := cli.ApiGET("/space/"+key, nil)
	if err != nil {
		return Space{}, fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Space{}, fmt.Errorf("[%d]%s", resp.StatusCode, resp.Status)
	}

	var info Space
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Space{}, fmt.Errorf("解析响应失败: %s", err)
	}

	return info, nil
}

//获取空间特定类型的内容
func (cli *Client) SpaceContentByType(key, contentType string, start int) ([]Content, int, error) {
	query := url.Values{
		"start":  {fmt.Sprintf("%d", start)},
		"expand": {"body.storage,ancestors"},
	}
	resp, err := cli.ApiGET("/space/"+key+"/content/"+contentType, query)
	if err != nil {
		return nil, 0, fmt.Errorf("执行请求失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("[%d]%s", resp.StatusCode, resp.Status)
	}

	var info struct {
		PageResp
		Results []Content
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, 0, fmt.Errorf("解析响应失败: %s", err)
	}

	//是否存在Next链接表示是否包含下一页
	nextStart := 0
	if info.Links.Next != "" {
		nextStart = info.Start + info.Size
	}

	return info.Results, nextStart, nil
}

//获取空间所有的页面
func (cli *Client) AllSpacePages(key string) ([]Content, error) {
	return cli.AllSpaceContents(key, ContentTypePage)
}

//获取空间所有的博客
func (cli *Client) AllSpaceBlogs(key string) ([]Content, error) {
	return cli.AllSpaceContents(key, ContentTypeBlog)
}

//获取空间所有的内容
func (cli *Client) AllSpaceContents(key, contentType string) ([]Content, error) {
	var pages []Content

	start := 0
	for {
		contents, nextStart, err := cli.SpaceContentByType(key, contentType, start)
		if err != nil {
			return nil, err
		}

		pages = append(pages, contents...)

		if nextStart <= 0 {
			break
		}

		start = nextStart
	}

	return pages, nil
}
