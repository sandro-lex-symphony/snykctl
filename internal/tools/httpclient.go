package tools

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"snykctl/internal/config"
	"time"
)

type HttpClient interface {
	Request(path string, verb string) *http.Response
	RequestGet(path string) *http.Response
	RequestDelete(path string) *http.Response
	RequestPost(path string, data []byte) *http.Response
	Url() string
	Token() string
	Id() string
	Timeout() int
	WorkerSize() int
}

type DefaultClient struct {
	config.ConfigProperties
	Debug bool
}

func NewHttpclient(conf config.ConfigProperties, d bool) DefaultClient {
	var c DefaultClient
	c.SetUrl(conf.Url())
	c.SetToken(conf.Token())
	c.SetId(conf.Id())
	c.SetTimeout(conf.Timeout())
	c.SetWorkerSize(conf.WorkerSize())
	c.SetSync(true)
	c.Debug = d
	return c
}

func (sc DefaultClient) Request(path string, verb string) *http.Response {
	timeout := time.Duration(time.Duration(sc.Timeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req := sc.Url() + path

	if sc.Debug {
		fmt.Println(verb, req)
	}

	request, err := http.NewRequest(verb, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", "token "+sc.Token())

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (sc DefaultClient) RequestGet(path string) *http.Response {
	return sc.Request(path, "GET")
}

func (sc DefaultClient) RequestDelete(path string) *http.Response {
	return sc.Request(path, "DELETE")
}

func (sc DefaultClient) RequestPost(path string, data []byte) *http.Response {
	timeout := time.Duration(time.Duration(sc.Timeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req := sc.Url() + path
	if sc.Debug {
		fmt.Println("POST", req)
		fmt.Println(string(data))
	}

	request, err := http.NewRequest("POST", req, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+sc.Token())

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
