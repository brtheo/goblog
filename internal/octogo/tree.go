package octogo

type TreeResponse struct {
	Sha       string        `json:"sha"`
	Url       string        `json:"url"`
	Tree      []TreeElement `json:"tree"`
	Truncated bool          `json:"truncated"`
}

type TreeElement struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
	Size int    `json:"size"`
	Url  string `json:"url"`
}
