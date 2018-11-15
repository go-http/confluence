package confluence

type Space struct {
	Id   int
	Key  string
	Name string
	Type string

	Expandable ExpandableResponse
	Links      LinkResponse `json:"_links"`
}
