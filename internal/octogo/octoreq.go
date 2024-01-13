package octogo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/brtheo/goblog/internal/octogo/util"

	it "github.com/BooleanCat/go-functional/iter"
)

type OctoReqConf struct {
	id       string
	endpoint string
	verb     string
}

type OctoRequest struct {
	conf OctoReqConf
}

func (c *OctoReqConf) Endpoint(endpoint string) *OctoReqConf {
	c.endpoint = endpoint
	return c
}
func (c *OctoReqConf) Id(id string) *OctoReqConf {
	c.id = id
	return c
}

func NewOctoReq(conf *OctoReqConf) *OctoRequest {
	return &OctoRequest{*conf}
}

func NewOctoReqConf() *OctoReqConf {
	return &OctoReqConf{
		verb: "GET",
	}
}
func CommitOctoReq(url string) *OctoRequest {
	return &OctoRequest{conf: OctoReqConf{
		verb:     "GET",
		endpoint: COMMIT_URL + "?sha=" + os.Getenv("GITBRANCH") + "&path=" + url,
	}}
}

func parseResponseInto[T any](resp *http.Response) (parsedStruct T) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	util.Throw(err)
	json.Unmarshal(body, &parsedStruct)
	return
}

func (o *Octogo) octoFetch(octoReq *OctoRequest) (resp *http.Response) {
	client := &http.Client{}
	url := BASE_URL + fmt.Sprintf("%s/", o.conf.user) + fmt.Sprintf("%s/", o.conf.repo) + octoReq.conf.endpoint
	fmt.Println(url)
	req, err := http.NewRequest(octoReq.conf.verb, url, nil)
	util.Throw(err)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITKEY")))
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err = client.Do(req)
	util.Throw(err)
	return resp
}

func (o *Octogo) octoFetches(octoReq OctoReqs) Responses {
	ch := make(chan Responses)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(req OctoReqs) {
		defer wg.Done()
		ch <- Responses{
			One: o.octoFetch(req.One),
			Two: o.octoFetch(req.Two),
		}
	}(octoReq)

	go func() {
		wg.Wait()
		close(ch)
	}()
	return it.FromChannel(ch).Next().Unwrap()
}
