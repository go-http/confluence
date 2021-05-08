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
	"path/filepath"
)

// 错误信息的响应结构
type ErrorResp struct {
	StatusCode int
	Data       ErrorData
	Message    string
	Reason     string
}

// 错误响应中的错误数据
type ErrorData struct {
	Authorized            bool
	Valid                 bool
	AllowedInReadOnlyMode bool
	Successful            bool
	Errors                []interface{}
}

//可供展开的字段信息
type ExpandableResponse map[string]string

// 响应信息中的链接信息
type LinkResp struct {
	Base     string
	Context  string
	Next     string
	Self     string
	WebUI    string
	Download string
}

// 响应信息中的分页信息
type PageResp struct {
	Size  int
	Start int
	Limit int
	Links LinkResp `json:"_links,omitempty"`
}

//下载指定链接的内容
func (cli *Client) Download(downloadUrl string) ([]byte, error) {
	u, err := url.Parse(downloadUrl)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Request("GET", u.Path, u.Query(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

//发起GET类型的API请求
func (cli *Client) ApiGET(path string, query url.Values) (*http.Response, error) {
	return cli.ApiRequest("GET", path, query, nil, nil)
}

//发起POST类型的API请求
func (cli *Client) ApiPOST(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}
	return cli.ApiRequest("POST", path, nil, nil, r)
}

//发起PUT类型的API请求
func (cli *Client) ApiPUT(path string, data interface{}) (*http.Response, error) {
	r, err := dataToJsonReader(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求数据失败: %s", err)
	}

	return cli.ApiRequest("PUT", path, nil, nil, r)
}

//发起POST类型的文件上传请求
func (cli *Client) ApiPOSTFiles(path string, files []string) (*http.Response, error) {
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
		"X-Atlassian-Token": {"no-check"},
		"Content-Type":      {w.FormDataContentType()},
	}

	return cli.ApiRequest("POST", path, nil, header, &body)
}

//发起指定方法的API请求
func (cli *Client) ApiRequest(method, path string, query, header url.Values, body io.Reader) (*http.Response, error) {
	return cli.Request(method, filepath.Join("/rest/api", path), query, header, body)
}

//执行指定的HTTP请求，执行前会自动添加上认证信息和Content-Type信息
func (cli *Client) Request(method, path string, query, header url.Values, body io.Reader) (*http.Response, error) {
	// 检查添加Query参数
	if query != nil {
		path += "?" + query.Encode()
	}

	// 构造请求
	req, err := http.NewRequest(method, cli.Hostname+path, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}

	if cli.Username != "" {
		req.SetBasicAuth(cli.Username, cli.Password)
	}

	req.Header.Set("Content-Type", "application/json")

	for name, _ := range header {
		req.Header.Set(name, header.Get(name))
	}

	return http.DefaultClient.Do(req)
}

//数据转换为JSON流reader
func dataToJsonReader(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("无法编码Data: %s", err)
	}

	return bytes.NewReader(jsonData), nil
}
