package octogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"os"
	"sync"

	"github.com/brtheo/goblog/internal/octogo/util"

	"github.com/BooleanCat/go-functional/v2/it"
)

type OctoReqConf struct {
	id       string
	endpoint string
	verb     string
	payload  []byte
}

type OctoRequest struct {
	conf OctoReqConf
}

func Endpoint(endpoint string) ConfFunc[OctoReqConf] {
	return func(conf *OctoReqConf) {
		conf.endpoint = endpoint
	}
}
func Payload(payload []byte) ConfFunc[OctoReqConf] {
	return func(conf *OctoReqConf) {
		conf.payload = payload
	}
}
func Id(id string) ConfFunc[OctoReqConf] {
	return func(conf *OctoReqConf) {
		conf.id = id
	}
}
func Verb(v string) ConfFunc[OctoReqConf] {
	return func(conf *OctoReqConf) {
		conf.verb = v
	}
}

func NewOctoReq(funs ...ConfFunc[OctoReqConf]) *OctoRequest {
	conf := _OctoReqConf()
	construct[OctoReqConf](&conf, funs)
	return &OctoRequest{conf}
}

func _OctoReqConf() OctoReqConf {
	return OctoReqConf{
		verb: "GET",
	}
}

func parseResponseInto[T any](resp *http.Response) (parsedStruct T) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	util.Throw(err)
	json.Unmarshal(body, &parsedStruct)
	// fmt.Printf("%+v\n", &parsedStruct)
	return
}

func (o *Octogo) octoFetch(octoReq *OctoRequest) (resp *http.Response) {
	client := &http.Client{}
	url := BASE_URL + fmt.Sprintf("%s/", o.conf.user) + fmt.Sprintf("%s/", o.conf.repo) + octoReq.conf.endpoint
	// fmt.Println(url)
	req, err := func() (*http.Request, error) {
		if octoReq.conf.verb == "PUT" {
			return http.NewRequest(octoReq.conf.verb, url, bytes.NewBuffer(octoReq.conf.payload))
		} else {
			return http.NewRequest(octoReq.conf.verb, url, nil)
		}
	}()
	util.Throw(err)
	k := os.Getenv("GITKEY")
	req.Header.Set("Authorization", "token "+k)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err = client.Do(req)
	util.Throw(err)
	return resp
}

func (o *Octogo) octoFetches(octoReq map[string]OctoReqs) iter.Seq[Responses] {
	ch := make(chan Responses)
	var wg sync.WaitGroup
	wg.Add(len(octoReq))
	for id, req := range octoReq {
		go func(req OctoReqs, id string) {
			defer wg.Done()
			ch <- Responses{
				One: o.octoFetch(req.One),
				Two: o.octoFetch(req.Two),
				id:  id,
			}
		}(req, id)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	return it.FromChannel(ch)
}
