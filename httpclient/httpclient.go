package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/kwitsch/omadaclient/log"
	"github.com/kwitsch/omadaclient/utils"
	"golang.org/x/net/publicsuffix"
)

type HttpClient struct {
	http *http.Client
	url  string
	l    *log.Log
}

func NewClient(url string, skipVerify, verbose bool) (*HttpClient, error) {
	l := log.New("HttpClient", verbose)
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

	result := HttpClient{
		http: &http,
		url:  url,
		l:    l,
	}
	return &result, nil
}

func (c *HttpClient) Get(path, body string, headers, params map[string]string) (*[]byte, error) {
	return c.request("GET", path, body, headers, params)
}

func (c *HttpClient) GetD(path, body string, headers, params map[string]string, result interface{}) error {
	res, err := c.Get(path, body, headers, params)
	if err != nil {
		return err
	}

	return c.decode(*res, &result)
}

func (c *HttpClient) Post(path, body string, headers, params map[string]string) (*[]byte, error) {
	return c.request("POST", path, body, headers, params)
}

func (c *HttpClient) PostD(path, body string, headers, params map[string]string, result interface{}) error {
	res, err := c.Post(path, body, headers, params)
	if err != nil {
		return err
	}

	return c.decode(*res, &result)
}

func (c *HttpClient) decode(data []byte, result interface{}) error {
	var aRes ApiResult
	if err := json.Unmarshal(data, &aRes); err != nil {
		return err
	}

	if !aRes.IsSuccess() {
		errCode, errMsg := aRes.GetHead()
		return utils.NewError("Errorcode:", errCode, "-", errMsg)
	}

	return aRes.GetResult(&result)
}

func (c *HttpClient) request(methode, path, body string, headers, params map[string]string) (*[]byte, error) {
	bodyData := []byte(body)
	url := c.url + path
	request, err := http.NewRequest(methode, url, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	query := request.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}

	query.Add("_", fmt.Sprint(time.Now().Unix()))

	request.URL.RawQuery = query.Encode()

	c.l.V("Request:", request.URL.String())

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.l.V("Response:", respBody)

	return &respBody, nil
}
