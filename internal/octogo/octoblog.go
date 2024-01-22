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
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/brtheo/goblog/internal/octogo/util"
	headingid "github.com/jkboxomine/goldmark-headingid"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	gutil "github.com/yuin/goldmark/util"
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
type Responses struct {
	One *http.Response
	Two *http.Response
	id  string
}
type Commits_B64 it.Pair[[]CommitResponse, string]

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

func (o *Octogo) GetAllPosts(nPerPage uint) []*Post {
	start := time.Now()

	treeElements := it.Lift[TreeElement](parseResponseInto[TreeResponse](
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

	p := it.Map[Responses, *Post](
		o.octoFetches(it.Fold[Slug_Id, map[string]OctoReqs](
			slug_id,
			map[string]OctoReqs{},
			func(m map[string]OctoReqs, s Slug_Id) map[string]OctoReqs {
				m[s.One] = OctoReqs{
					One: NewOctoReq(Endpoint(COMMIT_URL + "?sha=" + os.Getenv("GITBRANCH") + "&path=" + s.One)),
					Two: NewOctoReq(Endpoint(BLOB + s.Two)),
				}
				return m
			},
		)),
		func(r Responses) *Post {
			return o.newPost(Commits_B64{
				One: parseResponseInto[[]CommitResponse](r.One),
				Two: parseResponseInto[BlobResponse](r.Two).Content,
			}, r.id)
		},
	).Collect()
	By[Post](func(a, b *Post) bool {
		return a.PublishedAt > b.PublishedAt
	}).Sort(p)
	end := time.Since(start)
	fmt.Println(end)
	return p
}

func (o *Octogo) GetPostBySlug(slug string) *Post {
	_slug := slug + ".md"
	octoReqs := o.octoFetches(
		map[string]OctoReqs{
			"_": {
				One: NewOctoReq(Endpoint(COMMIT_URL + "?sha=" + os.Getenv("GITBRANCH") + "&path=" + _slug)),
				Two: NewOctoReq(Endpoint(CONTENT_URL + _slug + "?ref=" + os.Getenv("GITBRANCH"))),
			},
		},
	).Collect()
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
