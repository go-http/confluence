package confluence

//Confluence的空间
type Space struct {
	Id          int                 `json:"id,omitempty"`
	Key         string              `json:"key,omitempty"`
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	Icon        *SpaceIcon          `json:"icon,omitempty"`
	Description *SpaceDescription   `json:"description,omitempty"`
	HomePage    *Content            `json:"homePage,omitempty"`
	Metadata    *SpaceMetadata      `json:"metadata,omitempty"`
	Links       *LinkResp           `json:"_links,omitempty"`
	Expandable  *ExpandableResponse `json:"_expandable,omitempty"`
}

//Confluence的空间描述
type SpaceDescription struct {
	Plain RepresentationValue
	View  RepresentationValue
}

//Confluence的空间描述值
type RepresentationValue struct {
	Representation string
	Value          string
}

//Confluence的空间附加信息
type SpaceMetadata struct {
	Labels struct {
		PageResp
		Results []SpaceLabel
	}
}

//Confluence的空间标签
type SpaceLabel struct {
	Prefix string
	Name   string
	Id     string
}

//Confluence的空间图标
type SpaceIcon struct {
	Path      string
	Width     int
	Height    int
	IsDefault bool
}
