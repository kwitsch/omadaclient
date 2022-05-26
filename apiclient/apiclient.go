package apiclient

import (
	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/log"
	"github.com/kwitsch/omadaclient/model"
)

type Apiclient struct {
	url       string
	http      *httpclient.HttpClient
	id        string
	headers   map[string]string
	loginData model.Login
	l         *log.Log
}

func New(url string, skipVerify, verbose bool) (*Apiclient, error) {
	l := log.New("OmadaApi", verbose)
	http, err := httpclient.NewClient(url, skipVerify, verbose)
	if err != nil {
		return nil, l.E(err)
	}
	result := Apiclient{
		url:     url,
		http:    http,
		l:       l,
		headers: map[string]string{},
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
	var result model.ApiInfo
	if err := ac.http.GetD("/api/info", "", ac.headers, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.ReturnSuccess()

	return &result, nil
}

func (ac *Apiclient) Login(username, password string) error {
	ac.l.V("Login")
	bodyData := `{
		"username": "` + username + `",
		"password": "` + password + `"
	}`

	var result model.Login
	if err := ac.http.PostD(ac.getPath("login"), bodyData, ac.headers, &result); err != nil {
		return ac.l.E(err)
	}

	if result.Token == "" {
		return ac.l.E("Couldn't optain Logintoken.")
	}

	ac.loginData = result
	ac.headers = map[string]string{"Csrf-Token": result.Token}
	ac.l.ReturnSuccess()

	return nil
}

func (ac *Apiclient) LoginStatus() (bool, error) {
	ac.l.V("LoginStatus")
	var result model.LoginStatus
	if err := ac.http.GetD(ac.getPath("loginStatus"), "", ac.headers, &result); err != nil {
		return false, ac.l.E(err)
	}

	ac.l.Return(result.Login)
	return result.Login, nil
}

func (ac *Apiclient) Logout() error {
	ac.l.V("Logout")
	if _, err := ac.http.Post(ac.getPath("logout"), "", ac.headers); err != nil {
		return ac.l.E(err)
	}
	ac.l.ReturnSuccess()
	return nil
}

func (ac *Apiclient) getPath(endPoint string) string {
	return "/" + ac.id + "/api/v2/" + endPoint
}
