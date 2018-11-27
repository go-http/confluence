package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 错误数据的结构
type ErrorData struct {
	Authorized            bool
	Valid                 bool
	AllowedInReadOnlyMode bool
	Successful            bool
	Errors                []interface{}
}

// 错误信息的响应结构
type ErrorResp struct {
	StatusCode int
	Data       ErrorData
	Message    string
	Reason     string
}

type ExpandableResponse map[string]string

// 链接的响应结构
type LinkResp struct {
	Base    string
	Context string
	Next    string
	Self    string
	WebUI   string
}

// 分页的响应结构
type PageResp struct {
	Size  int
	Start int
	Limit int
	Links LinkResp `json:"_links"`
}

func (cli *Client) GET(path string, query url.Values) (*http.Response, error) {
	return cli.Request("GET", path, query, nil)
}

func (cli *Client) POST(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}
	return cli.Request("POST", path, nil, r)
}

func (cli *Client) PUT(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}

	return cli.Request("PUT", path, nil, r)
}

// 执行指定的HTTP请求，执行前会自动添加上认证信息和Content-Type信息
func (cli *Client) Request(method, path string, query url.Values, body io.Reader) (*http.Response, error) {
	// 检查添加Query参数
	if query != nil {
		path += "?" + query.Encode()
	}

	// 构造请求
	req, err := http.NewRequest(method, cli.Address+path, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(cli.Username, cli.Password)

	return http.DefaultClient.Do(req)
}

// 数据转换为JSON流reader
func dataToJsonReader(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("无法编码Data: %s", err)
	}

	return bytes.NewReader(jsonData), nil
}
