package httpclient

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type HttpClient struct {
	http  *http.Client
	url   string
	token string
}

func NewHttp(skipVerify bool) (*http.Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	cookies, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	http := http.Client{
		Jar: cookies,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipVerify,
			},
		},
	}

	return &http, nil
}

func NewClient(url string, skipVerify bool) (*HttpClient, error) {
	http, err := NewHttp(skipVerify)
	if err != nil {
		return nil, err
	}
	result := HttpClient{
		http: http,
		url:  url,
	}
	return &result, nil
}

func (c *HttpClient) SetToken(token string) {
	c.token = token
}

func (c *HttpClient) Get(path string) (*http.Response, error) {
	var bodyData = []byte{}
	request, err := http.NewRequest("GET", c.url+path, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if len(c.token) > 0 {
		request.Header.Set("Csrf-Token", c.token)
	}

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HttpClient) Post(path, body string) (*http.Response, error) {
	var bodyData = []byte(body)
	request, err := http.NewRequest("POST", c.url+path, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if len(c.token) > 0 {
		request.Header.Set("Csrf-Token", c.token)
	}

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
