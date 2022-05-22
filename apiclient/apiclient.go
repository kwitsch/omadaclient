package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/model"
	"github.com/kwitsch/omadaclient/utils"
)

type Apiclient struct {
	url   string
	http  *http.Client
	id    string
	token string
	role  int
}

func New(url string, skipVerify bool) (*Apiclient, error) {
	http, err := httpclient.New(skipVerify)
	if err != nil {
		return nil, err
	}
	result := Apiclient{
		url:  url,
		http: http,
	}

	ai, err := result.ApiInfo()
	if err != nil {
		return nil, err
	}

	if ai.OmadacId == "" {
		return nil, errors.New("Couldn't optain Omada ID.")
	}

	result.id = ai.OmadacId

	return &result, nil
}

func (ac *Apiclient) ApiInfo() (*model.ApiInfo, error) {
	resp, err := ac.http.Get(ac.url + "/api/info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result model.ApiInfoResponse
	if err := unmarshalResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Result, nil
}

func (ac *Apiclient) Login(username, password string) error {
	bodyData := `{
		"username": "` + username + `",
		"password": "` + password + `"
	}`

	var result model.LoginResponse
	if err := ac.request("POST", "login", bodyData, &result); err != nil {
		return err
	}

	if result.Result.Token == "" {
		return errors.New("Couldn't optain Logintoken.")
	}

	ac.token = result.Result.Token
	ac.role = result.Result.RoleType

	return nil
}

func (ac *Apiclient) LoginStatus() (bool, error) {
	var result model.LoginStatusResponse
	if err := ac.request("GET", "loginStatus", "", &result); err != nil {
		return false, err
	}

	return result.Result.Login, nil
}

func (ac *Apiclient) request(methode, endpoint, body string, result interface{}) error {
	if endpoint != "login" && ac.token == "" {
		return utils.NewError("Not logged in yet.")
	}

	var bodyData = []byte(body)
	request, err := http.NewRequest(methode, ac.getPath(endpoint), bytes.NewBuffer(bodyData))
	if err != nil {
		return err
	}

	if methode == "POST" {
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	if ac.token != "" {
		request.Header.Set("Csrf-Token", ac.token)
	}

	resp, err := ac.http.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := unmarshalResponse(resp, result); err != nil {
		return err
	}

	return nil
}

func (ac *Apiclient) getPath(endPoint string) string {
	return ac.url + "/" + ac.id + "/api/v2/" + endPoint
}

func unmarshalResponse(resp *http.Response, result interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	fmt.Println("- Debug Request:", resp.Request.URL.String())
	fmt.Println("-- Debug Response:", string(body))

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
