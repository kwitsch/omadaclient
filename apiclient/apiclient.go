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

// Creates a new Apiclient
// url: Omada controller address(example: https://192.168.0.2)
// siteName: Visible site name(empty string for default site)
// username: Username for login(it is advised to create a seperate api user)
// password: Password for login
// skipVerify: Ignore SSL errors(necessary for ip addresses as url or selfsigned certificates)
// verbose: Debug logging to console(should only be enabled for debugging scenarios)
// return Apiclient instance and possible error
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

// Initialize the Apiclient by logging in and fetching necesarry informations
// return possible error
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

	if len(ac.siteName) == 0 {
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

// Start current session
// return possible error
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

// Determins if current session is active
// return sesion state and possible error
func (ac *Apiclient) LoginStatus() (bool, error) {
	ac.l.V("LoginStatus")
	var result model.LoginStatus
	if err := ac.http.GetD(ac.getPath("loginStatus"), "", ac.headers, empty, &result); err != nil {
		return false, ac.l.E(err)
	}

	ac.l.Return(result.Login)
	return result.Login, nil
}

// Get user information for current session
// return user information and possible error
func (ac *Apiclient) UsersCurrent() (*model.UsersCurrent, error) {
	ac.l.V("UsersCurrent")
	var result model.UsersCurrent
	if err := ac.http.GetD(ac.getPath("users/current"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

// End current session
// return possible error
func (ac *Apiclient) Logout() error {
	ac.l.V("Logout")
	if _, err := ac.http.Post(ac.getPath("logout"), "", ac.headers, empty); err != nil {
		return ac.l.E(err)
	}
	ac.l.ReturnSuccess()
	return nil
}

// Fetches list of devices for all sites with basic information
// return list of devices and possible error
func (ac *Apiclient) Devices() (*[]model.Device, error) {
	ac.l.V("Devices")
	var result []model.Device
	if err := ac.http.GetD(ac.getSitesPath("devices"), "", ac.headers, empty, &result); err != nil {
		return nil, ac.l.E(err)
	}

	ac.l.Return(result)
	return &result, nil
}

// Fetches list of devices for initialized site with enhanced information
// return list of devices and possible error
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
// device: Device to enhance(Type and Mac have to be provided as minimal information)
// return possible error
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

// Get all active clients for initialised site
// return clients list and possible error
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
