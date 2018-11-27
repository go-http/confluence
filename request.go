package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	return cli.Request("GET", path, query, nil, nil)
}

func (cli *Client) POST(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}
	return cli.Request("POST", path, nil, nil, r)
}

func (cli *Client) PUT(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}

	return cli.Request("PUT", path, nil, nil, r)
}

func (cli *Client) POSTFiles(path string, files []string) (*http.Response, error) {
	var body bytes.Buffer

	w := multipart.NewWriter(&body)
	for _, file := range files {
		fw, err := w.CreateFormFile("file", file)
		if err != nil {
			return nil, fmt.Errorf("创建上传字段错误: %s", err)
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("读取文件%s错误: %s", file, err)
		}

		_, err = fw.Write(content)
		if err != nil {
			return nil, fmt.Errorf("添加上传文件%s错误: %s", file, err)
		}
	}
	w.Close()

	header := url.Values{
		"X-Atlassian-Token": {"nocheck"},
		"Content-Type":      {w.FormDataContentType()},
	}

	return cli.Request("POST", path, nil, header, &body)
}

// 执行指定的HTTP请求，执行前会自动添加上认证信息和Content-Type信息
func (cli *Client) Request(method, path string, query, header url.Values, body io.Reader) (*http.Response, error) {
	// 检查添加Query参数
	if query != nil {
		path += "?" + query.Encode()
	}

	// 构造请求
	req, err := http.NewRequest(method, cli.APIPrefix()+path, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}

	req.SetBasicAuth(cli.Username, cli.Password)
	req.Header.Set("Content-Type", "application/json")

	for name, _ := range header {
		req.Header.Set(name, header.Get(name))
	}

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
