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

const tokenKey string = "Csrf-Token"

// Creates a new Apiclient
//
// Parameters:
//   - `url` : Omada controller address(example: https://192.168.0.2)
//   - `siteName` : Visible site name(empty string for default site)
//   - `username` : Username for login(it is advised to create a seperate api user)
//   - `password` : Password for login
//   - `skipVerify` : Ignore SSL errors(necessary for ip addresses as url or selfsigned certificates)
//   - `erbose` : Debug logging to console(should only be enabled for debugging scenarios)
//
// return Apiclient instance or occured error
func New(url, siteName, username, password string, skipVerify, verbose bool) (*Apiclient, error) {
	l := log.New("OmadaApi", verbose)
	l.V("New")
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

	ai, err := result.ApiInfo()
	if err != nil {
		return nil, result.l.E(err)
	}

	if ai.OmadacId == "" {
		return nil, result.l.E("Couldn't optain Omada ID.")
	}

	result.omadaId = ai.OmadacId

	if len(result.siteName) != 0 {
		cu, err := result.UsersCurrent()
		if err != nil {
			return nil, result.l.E(err)
		}

		result.l.V("SiteName:", result.siteName)
		siteAvailable := false
		for _, v := range cu.Privilege.Sites {
			if v.Name == result.siteName {
				result.siteId = v.Key
				siteAvailable = true
				result.l.V("SiteId:", result.siteId)
				break
			}
		}
		if !siteAvailable {
			return nil, result.l.E("Site " + result.siteName + " is not available for user " + result.username)
		}
	}

	if err := result.EndSession(); err != nil {
		return nil, result.l.E(err)
	}

	return &result, nil
}

// Endsession wrapper for the use with defer
func (ac *Apiclient) Close() {
	ac.l.V("Close")
	if err := ac.EndSession(); err != nil {
		ac.l.E(err)
	}
}

// Fetch API information
//
// return API information or occured error
func (ac *Apiclient) ApiInfo() (*model.ApiInfo, error) {
	ac.l.V("ApiInfo")
	var result model.ApiInfo
	if err := ac.http.GetD("/api/info", "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.ReturnSuccess()

	return &result, nil
}

// Start session
//
// return occured error
func (ac *Apiclient) StartSession() error {
	ac.l.V("StartSession")
	if ac.HasActiveSession() {
		return nil
	}

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

	ac.setToken(result.Token)

	return nil
}

// Determins if is sn sctive session
//
// return sesion state
func (ac *Apiclient) HasActiveSession() bool {
	if !ac.hasToken() {
		return false
	}

	var result model.LoginStatus
	if err := ac.http.GetD(ac.getPath("loginStatus"), "", ac.headers, empty, &result); err != nil {
		ac.removeToken()
		ac.l.E(err)
		return false
	}

	if !result.Login {
		ac.removeToken()
	}

	return result.Login
}

// End session
//
// return occured error
func (ac *Apiclient) EndSession() error {
	ac.l.V("EndSession")
	if !ac.HasActiveSession() {
		return nil
	}

	if _, err := ac.http.Post(ac.getPath("logout"), "", ac.headers, empty); err != nil {
		return ac.l.E(err)
	}

	ac.removeToken()

	return nil
}

// Get user information for current session
//
// return user information or occured error
func (ac *Apiclient) UsersCurrent() (*model.UsersCurrent, error) {
	ac.l.V("UsersCurrent")
	if err := ac.ensureLoggedIn(); err != nil {
		return nil, ac.l.E(err)
	}

	var result model.UsersCurrent
	if err := ac.http.GetD(ac.getPath("users/current"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

// Fetches list of devices for all sites with basic information
//
// return list of devices or occured error
func (ac *Apiclient) Devices() (*[]model.Device, error) {
	ac.l.V("Devices")
	if err := ac.ensureLoggedIn(); err != nil {
		return nil, ac.l.E(err)
	}

	var result []model.Device
	if err := ac.http.GetD(ac.getSitesPath("devices"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

// Fetches list of devices for initialized site with enhanced information
//
// return list of devices or occured error
func (ac *Apiclient) DevicesDetailed() (*[]model.Device, error) {
	ac.l.V("DevicesDetailed")
	devices, err := ac.Devices()
	if err != nil {
		return nil, ac.l.E(err)
	}

	result := []model.Device{}
	for _, d := range *devices {
		if d.Site == ac.siteId {
			if err := ac.GetDeviceDetail(&d); err != nil {
				return nil, ac.l.E(err)
			}
			result = append(result, d)
		}
	}

	ac.l.Return(result)

	return &result, nil
}

// Fetch enhanced information for a provided device and enhance the struct by it
//
// Parameters:
//   - `device` : Device to enhance(Type and Mac have to be provided as minimal information)
//
// return occured error
func (ac *Apiclient) GetDeviceDetail(device *model.Device) error {
	ac.l.V("GetDeviceDetail")
	if err := ac.ensureLoggedIn(); err != nil {
		return ac.l.E(err)
	}

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

// Get all active clients for initialised site
//
// return clients list or occured error
func (ac *Apiclient) Clients() (*[]model.Client, error) {
	ac.l.V("Clients")
	if err := ac.ensureLoggedIn(); err != nil {
		return nil, ac.l.E(err)
	}

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

func (ac *Apiclient) hasToken() bool {
	_, ok := ac.headers[tokenKey]
	return ok
}

func (ac *Apiclient) setToken(token string) {
	ac.headers[tokenKey] = token
}

func (ac *Apiclient) removeToken() {
	delete(ac.headers, tokenKey)
}

func (ac *Apiclient) ensureLoggedIn() error {
	if ac.HasActiveSession() {
		return nil
	}

	return ac.StartSession()
}
