package octogo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/BooleanCat/go-functional/v2/it/itx"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/brtheo/goblog/internal/octogo/util"
	"github.com/google/uuid"
	headingid "github.com/jkboxomine/goldmark-headingid"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	gutil "github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	BASE_URL        = "https://api.github.com/repos/"
	TREE_URL        = "git/trees/"
	COMMIT_URL      = "commits"
	CONTENT_URL     = "contents/"
	BLOB            = "git/blobs/"
	COMMENTS_TREE   = "git/trees/comments?recursive=true"
	COMMENTS_COMMIT = "commits?sha=comments&path="
	CONTENTS        = "contents/"
)

type Pair[T any, K any] struct {
	One T
	Two K
}

type Slug_Id Pair[string, string]
type OctoReqs Pair[*OctoRequest, *OctoRequest]
type Responses struct {
	One *http.Response
	Two *http.Response
	id  string
}
type Commits_B64 Pair[[]CommitResponse, string]

type ConfFunc[T any] func(conf *T)

type OctoConf struct {
	user string
	repo string
}
type Octogo struct {
	conf OctoConf
}

func User(u string) ConfFunc[OctoConf] {
	return func(conf *OctoConf) {
		conf.user = u
	}
}

func Repo(r string) ConfFunc[OctoConf] {
	return func(conf *OctoConf) {
		conf.repo = r
	}
}

func _OctoConf() OctoConf {
	return OctoConf{
		user: "admin",
		repo: "default",
	}
}

func construct[T any](obj *T, funs []ConfFunc[T]) {
	for _, fn := range funs {
		fn(obj)
	}
}
func NewOctogo(funs ...ConfFunc[OctoConf]) *Octogo {
	conf := _OctoConf()
	construct[OctoConf](&conf, funs)
	return &Octogo{conf}
}

func (o *Octogo) newPost(commits_b64 Commits_B64, _slug string) *Post {
	var matter struct {
		Tags []string `yaml:"tags"`
	}
	last := len(commits_b64.One) - 1
	slug := strings.ReplaceAll(_slug, ".md", "")
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
func (o *Octogo) newComment(commits_b64 Commits_B64, _slug string) *CommentResponse {
	var matter struct {
		Author string `yaml:"author"`
	}
	last := len(commits_b64.One) - 1
	html := extractHTMLWithMatter(commits_b64.Two, &matter)

	return &CommentResponse{
		PublishedAt: int(commits_b64.One[last].Commit.Author.Date.Unix()),
		Content:     html,
		Author:      matter.Author,
	}
}

func (o *Octogo) GetAllPosts(nPerPage uint) []*Post {
	start := time.Now()

	treeElements := itx.FromSlice[TreeElement](parseResponseInto[TreeResponse](
		o.octoFetch(
			NewOctoReq(Endpoint(TREE_URL + os.Getenv("GITBRANCH"))),
		),
	).Tree).
		Filter(func(treeElement TreeElement) bool {
			return strings.Contains(treeElement.Path, ".md")
		}).Take(nPerPage)

	slug_id := it.Map[TreeElement, Slug_Id](treeElements, func(treeElement TreeElement) Slug_Id {
		return Slug_Id{
			One: treeElement.Path,
			Two: treeElement.Sha,
		}
	})

	p := slices.Collect(it.Map[Responses, *Post](
		o.octoFetches(it.Fold(
			slug_id,
			func(m map[string]OctoReqs, s Slug_Id) map[string]OctoReqs {
				m[s.One] = OctoReqs{
					One: NewOctoReq(Endpoint(COMMIT_URL + "?sha=" + os.Getenv("GITBRANCH") + "&path=" + s.One)),
					Two: NewOctoReq(Endpoint(BLOB + s.Two)),
				}
				return m
			},
			map[string]OctoReqs{},
		)),
		func(r Responses) *Post {
			return o.newPost(Commits_B64{
				One: parseResponseInto[[]CommitResponse](r.One),
				Two: parseResponseInto[BlobResponse](r.Two).Content,
			}, r.id)
		},
	))
	By[Post](func(a, b *Post) bool {
		return a.PublishedAt > b.PublishedAt
	}).Sort(p)
	end := time.Since(start)
	fmt.Println(end)
	return p
}

type CommentRequest struct {
	Message string `json:"message"`
	Branch  string `json:"branch"`
	Content string `json:"content"`
}

func (o *Octogo) CommitComment(props map[string]string, slug string) *CommentResponse {
	author := props["author"]
	body := props["body"]
	comment := "---\n" + "author: " + author + "\n" + "---\n" + body
	content := base64.StdEncoding.EncodeToString([]byte(comment))
	req := &CommentRequest{
		Message: "Message from " + author,
		Branch:  "comments",
		Content: content,
	}
	jjson, err := json.Marshal(req)
	util.Throw(err)
	uid := uuid.New()
	fmt.Println(slug)
	o.octoFetch(NewOctoReq(
		Verb("PUT"),
		Endpoint(CONTENTS+slug+"/"+uid.String()+".md"),
		Payload(jjson),
	))
	return &CommentResponse{
		Author:      author,
		PublishedAt: int(time.Now().Unix()),
		Content:     "<p>" + body + "</p>",
	}
}

func (o *Octogo) GetCommentsBySlug(slug string) Pair[[]*CommentResponse, string] {
	treeElements := itx.FromSlice[TreeElement](parseResponseInto[TreeResponse](
		o.octoFetch(
			NewOctoReq(Endpoint(COMMENTS_TREE)),
		),
	).Tree).
		Filter(func(treeElement TreeElement) bool {
			return strings.Contains(treeElement.Path, slug)
		})
	comments := []*CommentResponse{}
	elements := treeElements.Filter(func(treeElement TreeElement) bool {
		return strings.ContainsAny(treeElement.Path, "/")
	})
	slug_id := it.Map[TreeElement, Slug_Id](elements, func(treeElement TreeElement) Slug_Id {
		return Slug_Id{
			One: treeElement.Path,
			Two: treeElement.Sha,
		}
	})

	comments = slices.Collect(it.Map[Responses, *CommentResponse](
		o.octoFetches(it.Fold[Slug_Id, map[string]OctoReqs](
			slug_id,
			func(m map[string]OctoReqs, s Slug_Id) map[string]OctoReqs {
				m[s.One] = OctoReqs{
					One: NewOctoReq(Endpoint(COMMENTS_COMMIT + s.One)),
					Two: NewOctoReq(Endpoint(BLOB + s.Two)),
				}
				return m
			},
			map[string]OctoReqs{},
		)),
		func(r Responses) *CommentResponse {
			return o.newComment(Commits_B64{
				One: parseResponseInto[[]CommitResponse](r.One),
				Two: parseResponseInto[BlobResponse](r.Two).Content,
			}, r.id)
		},
	))
	if len(comments) > 0 {
		By[CommentResponse](func(a, b *CommentResponse) bool {
			return a.PublishedAt > b.PublishedAt
		}).Sort(comments)
	}
	return Pair[[]*CommentResponse, string]{One: comments, Two: slug}
}

func (o *Octogo) GetPostBySlug(slug string) *Post {
	_slug := slug + ".md"
	octoReqs := slices.Collect(o.octoFetches(
		map[string]OctoReqs{
			"_": {
				One: NewOctoReq(Endpoint(COMMIT_URL + "?sha=" + os.Getenv("GITBRANCH") + "&path=" + _slug)),
				Two: NewOctoReq(Endpoint(CONTENT_URL + _slug + "?ref=" + os.Getenv("GITBRANCH"))),
			},
		},
	))
	return o.newPost(Commits_B64{
		One: parseResponseInto[[]CommitResponse](octoReqs[0].One),
		Two: parseResponseInto[ContentResponse](octoReqs[0].Two).Content,
	}, slug)
}

func wrapperRenderer(w gutil.BufWriter, ctx hl.CodeBlockContext, entering bool) {
	language, ok := ctx.Language()
	lang := string(language)
	if ok && lang != "" {
		if entering {
			w.WriteString("<section data-lang=" + lang + ">")
		} else {
			w.WriteString(`</section>`)
		}
		return
	}
	if language == nil {
		if entering {
			w.WriteString("<pre><code>")
		} else {
			w.WriteString(`</code></pre>`)
		}
	}
}

func extractHTMLWithMatter[T any](content string, matter *T) string {
	unparsedFrontmatter, err := base64.StdEncoding.DecodeString(content)
	util.Throw(err)
	var buf bytes.Buffer
	ctx := parser.NewContext(parser.WithIDs(headingid.NewIDs()))
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			&frontmatter.Extender{},
			hl.NewHighlighting(
				hl.WithStyle("xcode-dark"),
				// hl.WithCustomStyle(chroma.MustNewStyle("mystyle", chroma.StyleEntries{
				// 	chroma.Background: "var(--crBackground)",
				// 	chroma.Keyword:    "var(-crKeyword)",
				// })),
				// hl.WithStyle("mystyle"),
				hl.WithWrapperRenderer(wrapperRenderer),
				hl.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(2),
				),
			),
		),
	)
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
