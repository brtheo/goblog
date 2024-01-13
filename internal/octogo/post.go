package octogo

type Post struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	PublishedAt    int      `json:"publishedAt"`
	LastModifiedAt int      `json:"lastModifiedAt"`
	Tags           []string `json:"tags"`
	Content        string   `json:"content"`
	Author         string   `json:"author"`
	AuthorPic      string   `json:"authorPic"`
}
