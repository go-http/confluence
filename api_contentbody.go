package confluence

import (
	"encoding/json"
	"fmt"
)

// 获取指定ID的内容
func (cli *Client) ContentBodyConvertTo(value, from, to string) (string, error) {
	data := ContentBodyStorage{
		Value:          value,
		Representation: from,
	}

	resp, err := cli.ApiPOST("/contentbody/convert/"+to, data)
	if err != nil {
		return "", fmt.Errorf("执行请求失败: %s", err)
	}

	defer resp.Body.Close()

	var info struct {
		ErrorResp
		ContentBodyStorage
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %s", err)
	}

	if info.StatusCode != 0 {
		return "", fmt.Errorf("[%d]%s", info.StatusCode, info.Message)
	}

	return info.ContentBodyStorage.Value, nil
}
