package octogo

type CommentResponse struct {
	Author      string `json:"author"`
	PublishedAt int    `json:"publishedAt"`
	Content     string `json:"content"`
}
