package confluence

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

type SpaceDescription struct {
	Plain RepresentationValue
	View  RepresentationValue
}

type RepresentationValue struct {
	Representation string
	Value          string
}

type SpaceMetadata struct {
	Labels struct {
		PageResp
		Results []SpaceLabel
	}
}

type SpaceLabel struct {
	Prefix string
	Name   string
	Id     string
}
