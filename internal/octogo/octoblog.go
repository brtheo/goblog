package octogo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	it "github.com/BooleanCat/go-functional/iter"
	"github.com/brtheo/goblog/internal/octogo/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

const (
	BASE_URL    = "https://api.github.com/repos/"
	TREE_URL    = "git/trees/"
	COMMIT_URL  = "commits"
	CONTENT_URL = "contents/"
	BLOB        = "git/blobs/"
)

type Slug_Id it.Pair[string, string]
type OctoReqs it.Pair[*OctoRequest, *OctoRequest]
type Responses it.Pair[*http.Response, *http.Response]
type Commits_B64 it.Pair[[]CommitResponse, string]

type OctoConf struct {
	user string
	repo string
}
type Octogo struct {
	conf OctoConf
}

func (conf OctoConf) User(u string) OctoConf {
	conf.user = u
	return conf
}

func (conf OctoConf) Repo(r string) OctoConf {
	conf.repo = r
	return conf
}

func NewOctoConf() OctoConf {
	return OctoConf{
		user: "admin",
		repo: "default",
	}
}

func NewOctogo(conf OctoConf) Octogo {
	return Octogo{conf}
}

func (o *Octogo) newPost(commits_b64 Commits_B64, _slug string) *Post {
	var matter struct {
		Tags []string `yaml:"tags"`
	}
	last := len(commits_b64.One) - 1
	slug := strings.Replace(_slug, ".md", "", 0)
	title := strings.ReplaceAll(slug, "-", " ")
	// fmt.Println(pair.One)
	html := extractHTMLWithMatter(commits_b64.Two, &matter)
	return &Post{
		Slug:        slug,
		Title:       title,
		PublishedAt: int(commits_b64.One[last].Commit.Author.Date.Unix()),
		Content:     html,
		Author:      commits_b64.One[last].Committer.Login,
		AuthorPic:   commits_b64.One[last].Committer.Avatar_url,
		Tags:        matter.Tags,
	}
}

func (o *Octogo) GetAllPosts(nPerPage uint) []*Post {
	start := time.Now()

	treeElements := it.Lift[TreeElement](
		parseResponseInto[TreeResponse](o.octoFetch(NewOctoReq(NewOctoReqConf().Endpoint(TREE_URL + os.Getenv("GITBRANCH"))))).Tree,
	).Filter(func(treeElement TreeElement) bool {
		return strings.Contains(treeElement.Path, ".md")
	}).Take(nPerPage)

	slug_id := it.Map[TreeElement, Slug_Id](treeElements, func(treeElement TreeElement) Slug_Id {
		return Slug_Id{
			One: treeElement.Path,
			Two: treeElement.Sha,
		}
	})

	it.Map[Slug_Id, *Post](slug_id, o.getPost)

	p := it.Map[Slug_Id, *Post](slug_id, o.getPost).Collect()
	By[Post](func(a, b *Post) bool {
		return a.PublishedAt > b.PublishedAt
	}).Sort(p)
	end := time.Since(start)
	fmt.Println(end)
	return p
}

func (o *Octogo) getPost(slug_id Slug_Id) *Post {
	octoReqs := o.octoFetches(
		OctoReqs{
			One: CommitOctoReq(slug_id.One),
			Two: NewOctoReq(NewOctoReqConf().Endpoint(BLOB + slug_id.Two)),
		},
	)
	fmt.Print(slug_id.Two)
	return o.newPost(Commits_B64{
		One: parseResponseInto[[]CommitResponse](octoReqs.One),
		Two: parseResponseInto[BlobResponse](octoReqs.Two).Content,
	}, slug_id.One)
}

func (o *Octogo) GetPostBySlug(_slug string) *Post {
	slug := _slug + ".md"
	octoReqs := o.octoFetches(
		OctoReqs{
			One: CommitOctoReq(_slug),
			Two: NewOctoReq(NewOctoReqConf().Endpoint(CONTENT_URL + _slug + "?ref=" + os.Getenv("GITBRANCH"))),
		},
	)
	return o.newPost(Commits_B64{
		One: parseResponseInto[[]CommitResponse](octoReqs.One),
		Two: parseResponseInto[ContentResponse](octoReqs.Two).Content,
	}, slug)
}

func extractHTMLWithMatter[T any](content string, matter *T) string {
	unparsedFrontmatter, err := base64.StdEncoding.DecodeString(content)
	util.Throw(err)
	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(&frontmatter.Extender{}))
	ctx := parser.NewContext()
	if err := md.Convert([]byte(unparsedFrontmatter), &buf, parser.WithContext(ctx)); err != nil {
		log.Fatal(err)
	}
	fm := frontmatter.Get(ctx)
	if fm == nil {
		log.Fatal("no frontmatter found")
	}
	if err := fm.Decode(&matter); err != nil {
		log.Fatal(err)
	}
	return string(buf.Bytes())
}
