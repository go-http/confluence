package confluence

type Version struct {
	By        *User    `json:"by,omitempty"`
	When      string   `json:"when,omitempty"`
	Message   string   `json:"message,omitempty"`
	Number    int      `json:"number,omitempty"`
	MinorEdit bool     `json:"minorEdit,omitempty"`
	Hidden    bool     `json:"hidden,omitempty"`
	Links     LinkResp `json:"_links"`
}
