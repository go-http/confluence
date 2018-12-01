package confluence

//Confluence内容
type Content struct {
	Id        string      `json:"id,omitempty"`
	Type      string      `json:"type,omitempty"`
	Title     string      `json:"title,omitempty"`
	Space     Space       `json:"space,omitempty"`
	Body      ContentBody `json:"body,omitempty"`
	Link      LinkResp    `json:"_links",omitempty`
	Version   Version     `json:"version,omitempty"`
	Ancestors []Content   `json:"ancestors,omitempty"`
}

const (
	ContentTypePage = "page" //页面类型的Content
	ContentTypeBlog = "blog" //博客类型的Content
)

//Confluence内容体
type ContentBody struct {
	Storage             ContentBodyStorage `json:"storage,omitempty"`
	Editor              interface{}        `json:"editor,omitempty"`
	View                interface{}        `json:"view,omitempty"`
	ExportView          interface{}        `json:"export_view,omitempty"`
	StyledView          interface{}        `json:"styled_view,omitempty"`
	AnonymousExportView interface{}        `json:"anonymous_export_view,omitempty"`
}

//Storage类型的内容体
type ContentBodyStorage struct {
	Value          string `json:"value,omitempty"`
	Representation string `json:"representation,omitempty"`
}

//设置Storage类型的内容体
func (content *Content) SetStorageBody(value string) {
	content.Body.Storage.Representation = "storage"
	content.Body.Storage.Value = value
}
