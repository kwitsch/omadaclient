package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/model"
)

type Apiclient struct {
	url   string
	id    string
	token string
	http  *http.Client
}

func New(url string, skipVerify bool) (*Apiclient, error) {
	http, err := httpclient.New(skipVerify)
	if err != nil {
		return nil, err
	}

	id, err := getId(url, http)
	if err != nil {
		return nil, err
	}

	result := Apiclient{
		url:  url,
		id:   id,
		http: http,
	}

	fmt.Println("Optained id", id, "for", url)
	return &result, nil
}

func getId(url string, http *http.Client) (string, error) {
	resp, err := http.Get(url + "/api/info")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("response:", body)
	var result model.ApiInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.ErrorCode != 0 || result.Result.OmadacId == "" {
		return "", errors.New("Couldn't optain Omada ID.")
	}

	return result.Result.OmadacId, nil
}
