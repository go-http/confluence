package confluence

//Confluence客户端，封装了Confluence的常见资源操作
type Client struct {
	Hostname string
	Username string
	Password string
}

//创建新的Confluence客户端
func New(addr, user, pass string) *Client {
	return &Client{
		Hostname: addr,
		Username: user,
		Password: pass,
	}
}

//获取指定内容的附件访问前缀
func (cli *Client) ContentAttachmentUrlPrefix(contentId string) string {
	return cli.Hostname + "/download/attachments/" + contentId + "/"
}
