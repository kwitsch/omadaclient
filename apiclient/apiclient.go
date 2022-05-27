package apiclient

import (
	"fmt"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/log"
	"github.com/kwitsch/omadaclient/model"
)

type Apiclient struct {
	url      string
	siteName string
	username string
	password string
	http     *httpclient.HttpClient
	omadaId  string
	siteId   string
	headers  map[string]string
	l        *log.Log
}

var empty = map[string]string{}

func New(url, siteName, username, password string, skipVerify, verbose bool) (*Apiclient, error) {
	l := log.New("OmadaApi", verbose)
	http, err := httpclient.NewClient(url, skipVerify, verbose)
	if err != nil {
		return nil, l.E(err)
	}
	result := Apiclient{
		url:      url,
		siteName: siteName,
		username: username,
		password: password,
		http:     http,
		l:        l,
		siteId:   "Default",
		headers:  map[string]string{},
	}

	return &result, nil
}

func (ac *Apiclient) Start() error {
	ac.l.V("Start")
	ai, err := ac.ApiInfo()
	if err != nil {
		return ac.l.E(err)
	}

	if ai.OmadacId == "" {
		return ac.l.E("Couldn't optain Omada ID.")
	}

	ac.omadaId = ai.OmadacId

	if err := ac.Login(); err != nil {
		return ac.l.E(err)
	}

	cu, err := ac.UsersCurrent()
	if err != nil {
		return ac.l.E(err)
	}

	ac.l.V("SiteName:", ac.siteName)
	siteAvailable := false
	for _, v := range cu.Privilege.Sites {
		if v.Name == ac.siteName {
			ac.siteId = v.Key
			siteAvailable = true
			ac.l.V("SiteId:", ac.siteId)
			break
		}
	}
	if !siteAvailable {
		return ac.l.E("Site " + ac.siteName + " is not available for user " + ac.username)
	}

	ac.l.ReturnSuccess()

	return nil
}

func (ac *Apiclient) ApiInfo() (*model.ApiInfo, error) {
	ac.l.V("ApiInfo")
	var result model.ApiInfo
	if err := ac.http.GetD("/api/info", "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.ReturnSuccess()

	return &result, nil
}

func (ac *Apiclient) Login() error {
	ac.l.V("Login")
	bodyData := `{
		"username": "` + ac.username + `",
		"password": "` + ac.password + `"
	}`

	var result model.Login
	if err := ac.http.PostD(ac.getPath("login"), bodyData, ac.headers, empty, &result); err != nil {
		return ac.l.E(err)
	}

	if result.Token == "" {
		return ac.l.E("Couldn't optain Logintoken.")
	}

	ac.headers = map[string]string{"Csrf-Token": result.Token}
	ac.l.ReturnSuccess()

	return nil
}

func (ac *Apiclient) LoginStatus() (bool, error) {
	ac.l.V("LoginStatus")
	var result model.LoginStatus
	if err := ac.http.GetD(ac.getPath("loginStatus"), "", ac.headers, empty, &result); err != nil {
		return false, ac.l.E(err)
	}

	ac.l.Return(result.Login)
	return result.Login, nil
}

func (ac *Apiclient) UsersCurrent() (*model.UsersCurrent, error) {
	ac.l.V("UsersCurrent")
	var result model.UsersCurrent
	if err := ac.http.GetD(ac.getPath("users/current"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

func (ac *Apiclient) Logout() error {
	ac.l.V("Logout")
	if _, err := ac.http.Post(ac.getPath("logout"), "", ac.headers, empty); err != nil {
		return ac.l.E(err)
	}
	ac.l.ReturnSuccess()
	return nil
}

func (ac *Apiclient) Devices() (*[]model.Device, error) {
	ac.l.V("Devices")
	var result []model.Device
	if err := ac.http.GetD(ac.getSitesPath("devices"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

func (ac *Apiclient) DevicesDetailed() (*[]model.Device, error) {
	ac.l.V("DevicesDetailed")
	devices, err := ac.Devices()
	if err != nil {
		return nil, ac.l.E(err)
	}

	result := []model.Device{}
	for _, d := range *devices {
		if err := ac.GetDeviceDetail(&d); err != nil {
			return nil, ac.l.E(err)
		}
		result = append(result, d)
	}

	ac.l.Return(result)
	return &result, nil
}

func (ac *Apiclient) GetDeviceDetail(device *model.Device) error {
	var dtype string
	switch device.Type {
	case "switch":
		dtype = "switches"
	case "gateway":
		dtype = "gateways"
	case "ap":
		dtype = "eaps"
	default:
		return ac.l.E("Unknown device type: " + device.Type)
	}

	if err := ac.http.GetD(ac.getSitesPath(dtype+"/"+device.Mac), "", ac.headers, empty, device); err != nil {
		return ac.l.E(err)
	}

	ac.l.Return(*device)
	return nil
}

func (ac *Apiclient) Clients() (*[]model.Client, error) {
	result := []model.Client{}
	page := 1
	params := map[string]string{
		"currentPageSize": "10",
	}

	for {
		params["currentPage"] = fmt.Sprint(page)

		var hres model.Clients
		if err := ac.http.GetD(ac.getSitesPath("clients"), "", ac.headers, params, &hres); err != nil {
			return nil, ac.l.E(err)
		}
		result = append(result, hres.Data...)
		page = hres.CurrentPage + 1

		if len(result) >= hres.TotalRows {
			break
		}
	}

	return &result, nil
}

func (ac *Apiclient) getPath(endPoint string) string {
	return "/" + ac.omadaId + "/api/v2/" + endPoint
}

func (ac *Apiclient) getSitesPath(endPoint string) string {
	return ac.getPath("sites/" + ac.siteId + "/" + endPoint)
}
