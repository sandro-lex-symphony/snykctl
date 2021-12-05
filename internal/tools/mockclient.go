package tools

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"snykctl/internal/config"
)

type MockClient struct {
	config.ConfigProperties
	ResponseBody string
	Status       string
	StatusCode   int
	Debug        bool
}

func NewMockClient() MockClient {
	var client MockClient
	client.SetUrl("")
	return client
}

func (c MockClient) Request(path string, verb string) *http.Response {
	t := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(c.ResponseBody)),
	}

	t.StatusCode = c.StatusCode
	t.Status = c.Status

	return &t

}

func (c MockClient) RequestGet(path string) *http.Response {
	return c.Request(path, "GET")

}

func (c MockClient) RequestDelete(path string) *http.Response {
	return c.Request(path, "DELETE")

}

func (c MockClient) RequestPost(path string, data []byte) *http.Response {
	return c.Request(path, "POST")

}

func (c MockClient) Url() string {
	return ""
}

func (c MockClient) Token() string {
	return ""
}

func (c MockClient) Id() string {
	return ""
}

func (c MockClient) Timeout() int {
	return 0
}

func (c MockClient) WorkerSize() int {
	return 0
}
