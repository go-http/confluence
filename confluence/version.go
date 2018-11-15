package confluence

type Version struct {
	By        User
	When      string
	Message   string
	Number    int
	MinorEdit bool
	Hidden    bool
	Links     LinkResp `josn:"_links"`
}
