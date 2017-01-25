package roku

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type Command struct {
	count      int
	key        string
	url        string
	httpClient *http.Client
}

type HttpResponse struct {
	count    int
	response *http.Response
	err      error
}

func Keypress(key string, count int) error {
	fmt.Printf("Hitting key %v, %v times\n", key, count)
	var command = Command{key: key, count: count, httpClient: http.DefaultClient}
	command.buildURL()

	req, err := command.buildReq()
	if err != nil {
		return err
	}
	reqs := command.do(req)
	for _, req := range reqs {
		time.Sleep(1 * time.Second)
		if req.err != nil {
			return fmt.Errorf("Not all parts called")
		}
	}
	return nil
}

func (c *Command) buildURL() {
	var buffer bytes.Buffer
	buffer.WriteString(rokuIP)
	buffer.WriteString("keypress")
	buffer.WriteString(c.key)
	c.url = buffer.String()
	return
}

func (c *Command) do(req *http.Request) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	for i := 0; i < c.count; i++ {
		go func(i int) {
			resp, err := c.httpClient.Do(req)
			ch <- &HttpResponse{i, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(i)
	}

	for {
		select {
		case r := <-ch:
			if r.err != nil {
				fmt.Printf("Had trouble the %v request, for key %v\n", r.count, c.key)
			}
			responses = append(responses, r)
			if len(responses) == c.count {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

func (c *Command) buildReq() (*http.Request, error) {
	req, err := http.NewRequest("POST", c.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Golang-Roku")
	return req, nil
}
