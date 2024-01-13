package octogo

import "time"

type CommitResponse struct {
	Sha         string        `json:"sha"`
	NodeId      string        `json:"node_id"`
	Commit      Commit        `json:"commit"`
	Url         string        `json:"url"`
	HtmlUrl     string        `json:"html_url"`
	CommentsUrl string        `json:"comments_url"`
	Author      CommitsAuthor `json:"author"`
	Committer   CommitsAuthor `json:"committer"`
	Parents     Parent        `json:"parents"`
}

type CommitsAuthor struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	Avatar_url        string `json:"avatar_url"`
	Gravatar_id       string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Commit struct {
	Author       CommitAuthor `json:"author"`
	Committer    CommitAuthor `json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `json:"tree"`
	Url          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

type CommitAuthor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Tree struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Verification struct {
	Verified bool   `json:"verified"`
	Reason   string `json:"reason"`
	// Signature nil    `json:"signature"`
	// Payload   nil    `json:"payload"`
}

type Parent struct {
	Sha     string `json:"sha"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}
