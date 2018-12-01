package confluence

type ContentBody struct {
	Storage             Storage     `json:"storage,omitempty"`
	Editor              interface{} `json:"editor,omitempty"`
	View                interface{} `json:"view,omitempty"`
	ExportView          interface{} `json:"export_view,omitempty"`
	StyledView          interface{} `json:"styled_view,omitempty"`
	AnonymousExportView interface{} `json:"anonymous_export_view,omitempty"`
}

type Content struct {
	Id        string      `json:"id,omitempty"`
	Type      string      `json:"type,omitempty"`
	Title     string      `json:"title,omitempty"`
	Space     Space       `json:"space,omitempty"`
	Body      ContentBody `json:"body,omitempty"`
	Version   Version     `json:"version,omitempty"`
	Ancestors []Content   `json:"ancestors,omitempty"`
	Link      LinkResp    `json:"_links"`
}
