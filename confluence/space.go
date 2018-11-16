package confluence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type RepresentationValue struct {
	Representation string
	Value          string
}

type SpaceDescription struct {
	Plain RepresentationValue
	View  RepresentationValue
}

type SpaceLabel struct {
	Prefix string
	Name   string
	Id     string
}

type SpaceMetadata struct {
	Labels struct {
		PageResp
		Results []SpaceLabel
	}
}

type Space struct {
	Id          int                 `json:"id,omitempty"`
	Key         string              `json:"key,omitempty"`
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	Icon        *Icon               `json:"icon,omitempty"`
	Description *SpaceDescription   `json:"description,omitempty"`
	HomePage    *Content            `json:"homePage,omitempty"`
	Metadata    *SpaceMetadata      `json:"metadata,omitempty"`
	Links       *LinkResp           `json:"_links,omitempty"`
	Expandable  *ExpandableResponse `json:"_expandable,omitempty"`
}

func (cli *Client) SpaceByKey(key string) (Space, error) {
	resp, err := cli.GET("/space/"+key, nil)
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
func (cli *Client) SpaceContentExportToPath(key, outDir string) error {
	resp, err := cli.GET("/space/"+key+"/content", url.Values{"expand": {"body.storage,ancestors"}})
	if err != nil {
		return fmt.Errorf("执行请求失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[%d]%s", resp.StatusCode, resp.Status)
	}

	var info struct {
		Page struct {
			PageResp
			LinkResp
			Results []Content
		}
		BlogPost struct {
			PageResp
			LinkResp
			Results []Content
		}
		LinkResp
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return fmt.Errorf("解析响应失败: %s", err)
	}

	//清空原目录
	os.RemoveAll(outDir)

	pageOutDir := path.Join(outDir, "page")
	os.MkdirAll(pageOutDir, 0755)
	for i, page := range info.Page.Results {
		fmt.Printf("[%3d/%3d] %s(%d Bytes)", i+1, len(info.Page.Results), page.Title, len(page.Body.Storage.Value))

		//获取父页面信息
		ancestorDirs := make([]string, len(page.Ancestors))
		for i, ancestor := range page.Ancestors {
			ancestorDirs[i] = ancestor.Title
		}

		pageDirs := append([]string{pageOutDir}, ancestorDirs...)

		pageDir := path.Join(pageDirs...)
		os.MkdirAll(pageDir, 0755)

		file := path.Join(pageDir, page.Title+".xml")

		fmt.Println("       =>", file)
		ioutil.WriteFile(file, []byte(page.Body.Storage.Value), 0755)
	}

	postOutDir := path.Join(outDir, "post")
	os.MkdirAll(postOutDir, 0755)
	for i, post := range info.BlogPost.Results {
		fmt.Printf("[%3d/%3d] %s(%d Bytes)", i+1, len(info.BlogPost.Results), post.Title, len(post.Body.Storage.Value))

		//获取父页面信息
		ancestorDirs := make([]string, len(post.Ancestors))
		for i, ancestor := range post.Ancestors {
			ancestorDirs[i] = ancestor.Title
		}

		postDirs := append([]string{postOutDir}, ancestorDirs...)

		postDir := path.Join(postDirs...)
		os.MkdirAll(postDir, 0755)

		file := path.Join(postDir, post.Title+".xml")

		fmt.Println("       =>", file)
		ioutil.WriteFile(file, []byte(post.Body.Storage.Value), 0755)
	}

	return nil
}
