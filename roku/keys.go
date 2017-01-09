package roku

import (
	"bytes"
	"fmt"
	"net/http"
)

type Command struct {
	count      int
	key        string
	url        string
	httpClient *http.Client
}

func Keypress(key string, count int) error {
	fmt.Printf("Hitting key %v, %v times\n", key, count)
	var command = Command{key: key, count: count, httpClient: http.DefaultClient}
	command.buildURL()

	req, err := command.buildReq()
	if err != nil {
		return err
	}
	command.do(req)
	if err != nil {
		return nil
	}
	fmt.Printf(command.url)
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

func (c *Command) do(req *http.Request) (*http.Response, error) {
	var lastResp *http.Response
	var sendErr error
	sendErr = nil
	for i := 0; i < c.count; i++ {
		go func(i int) {
			resp, err := c.httpClient.Do(req)
			if err != nil {
				sendErr = err
			}
			if resp.StatusCode != 200 {
				sendErr = fmt.Errorf("Got something other than 200")
				lastResp = resp
			}
		}(i)
	}
	return lastResp, sendErr
}

func (c *Command) buildReq() (*http.Request, error) {
	req, err := http.NewRequest("POST", c.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Golang-Roku")
	return req, nil
}
