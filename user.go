package confluence

type User struct {
	Type           string      `json:"type,omitempty"`
	Username       string      `json:"username,omitempty"`
	UserKey        string      `json:"userKey,omitempty"`
	ProfilePicture interface{} `json:"profilePicture,omitempty"`
	DisplayName    string      `json:"displayName,omitempty"`
	Links          *LinkResp   `json:"_links"`
}
