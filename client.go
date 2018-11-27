package confluence

type Client struct {
	Hostname string
	Username string
	Password string
}

func New(addr, user, pass string) *Client {
	return &Client{
		Hostname: addr,
		Username: user,
		Password: pass,
	}
}

func (cli *Client) AttachmentUrlPrefix(contentId string) string {
	return cli.Hostname + "/download/attachments/" + contentId + "/"
}

func (cli *Client) APIPrefix() string {
	return cli.Hostname + "/rest/api"
}
