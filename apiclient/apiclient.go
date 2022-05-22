package apiclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/log"
	"github.com/kwitsch/omadaclient/model"
	"github.com/kwitsch/omadaclient/utils"
)

type Apiclient struct {
	url       string
	http      *http.Client
	id        string
	loginData model.Login
	l         *log.Log
}

func New(url string, skipVerify, verbose bool) (*Apiclient, error) {
	l := log.New("OmadaApi", verbose)
	http, err := httpclient.New(skipVerify)
	if err != nil {
		return nil, l.E(err)
	}
	result := Apiclient{
		url:  url,
		http: http,
		l:    l,
	}

	ai, err := result.ApiInfo()
	if err != nil {
		return nil, result.l.E(err)
	}

	if ai.OmadacId == "" {
		return nil, result.l.E("Couldn't optain Omada ID.")
	}

	result.id = ai.OmadacId

	return &result, nil
}

func (ac *Apiclient) ApiInfo() (*model.ApiInfo, error) {
	ac.l.V("ApiInfo")
	resp, err := ac.http.Get(ac.url + "/api/info")
	if err != nil {
		return nil, ac.l.E(err)
	}
	defer resp.Body.Close()

	var result model.ApiInfoResponse
	if err := ac.unmarshalResponse(resp, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result.Result)
	return &result.Result, nil
}

func (ac *Apiclient) Login(username, password string) error {
	ac.l.V("Login")
	bodyData := `{
		"username": "` + username + `",
		"password": "` + password + `"
	}`

	var result model.LoginResponse
	if err := ac.request("POST", "login", bodyData, &result); err != nil {
		return ac.l.E(err)
	}

	if result.Result.Token == "" {
		return ac.l.E("Couldn't optain Logintoken.")
	}

	ac.loginData = result.Result
	ac.l.Return(result.Msg)

	return nil
}

func (ac *Apiclient) LoginStatus() (bool, error) {
	ac.l.V("LoginStatus")
	var result model.LoginStatusResponse
	if err := ac.request("GET", "loginStatus", "", &result); err != nil {
		return false, ac.l.E(err)
	}

	ac.l.Return(result.Result.Login)
	return result.Result.Login, nil
}

func (ac *Apiclient) request(methode, endpoint, body string, result interface{}) error {
	if endpoint != "login" && ac.loginData.Token == "" {
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

	if ac.loginData.Token != "" {
		request.Header.Set("Csrf-Token", ac.loginData.Token)
	}

	resp, err := ac.http.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := ac.unmarshalResponse(resp, result); err != nil {
		return err
	}

	return nil
}

func (ac *Apiclient) getPath(endPoint string) string {
	return ac.url + "/" + ac.id + "/api/v2/" + endPoint
}

func (ac *Apiclient) unmarshalResponse(resp *http.Response, result interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	ac.l.V("Http Request:", resp.Request.URL.String())
	ac.l.V("Http Response:", body)

	if err := ac.testResponse(result); err != nil {
		return err
	}

	return nil
}

func (ac *Apiclient) testResponse(respObj interface{}) error {
	rv := reflect.ValueOf(respObj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return utils.NewError("Response is not a valid struct")
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
		return utils.NewError("ErrorCode is missing")
	}
}
