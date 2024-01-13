package octogo

type BlobResponse struct {
	Sha      string `json:"sha"`
	Node_id  string `json:"node_id"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}
