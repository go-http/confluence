package confluence

type Client struct {
	Address  string
	Username string
	Password string
}

func New(addr, user, pass string) *Client {
	return &Client{
		Address:  addr + "/rest/api",
		Username: user,
		Password: pass,
	}
}
