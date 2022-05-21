package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/model"
	"github.com/kwitsch/omadaclient/utils"
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

func (ac *Apiclient) Login(username, password string) (bool, error) {
	return true, nil
}

func (ac *Apiclient) getPath(endPoint string) string {
	return ac.url
}

func getId(url string, http *http.Client) (string, error) {
	resp, err := http.Get(url + "/api/info")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result model.ApiInfoResponse
	if err := unmarshalResponse(resp.Body, &result); err != nil {
		return "", err
	}

	if result.Result.OmadacId == "" {
		return "", errors.New("Couldn't optain Omada ID.")
	}

	return result.Result.OmadacId, nil
}

func unmarshalResponse(r io.Reader, result interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if err := testResponse(result); err != nil {
		return err
	}

	return nil
}

func testResponse(respObj interface{}) error {
	rv := reflect.ValueOf(respObj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("Response is not a valid struct")
	}
	ecF := rv.FieldByName("ErrorCode")
	if ecF.IsValid() {
		ecV := ecF.Int()
		if ecV == 0 {
			return nil
		} else {
			mF := rv.FieldByName("Msg")
			if mF.IsValid() {
				return utils.NewError("ErrorCode:", ecV, "Message:", mF.String())
			} else {
				return utils.NewError("ErrorCode:", ecV)
			}
		}

	} else {
		return errors.New("ErrorCode is missing")
	}
}
