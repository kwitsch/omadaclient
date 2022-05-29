package omadaclient

import (
	"fmt"

	"github.com/kwitsch/omadaclient/httpclient"
	"github.com/kwitsch/omadaclient/log"
	"github.com/kwitsch/omadaclient/model"
)

type SiteClient struct {
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

// Creates a new SiteClient
//
// Parameters:
//   url : Omada controller address(example: https://192.168.0.2)
//   siteName : Visible site name(empty string for default site)
//   username : Username for login(it is advised to create a seperate api user)
//   password : Password for login
//   skipVerify : Ignore SSL errors(necessary for ip addresses as url or selfsigned certificates)
//   erbose : Debug logging to console(should only be enabled for debugging scenarios)
//
// Return:
//   SiteClient instance
//   error
func NewSiteClient(url, siteName, username, password string, skipVerify, verbose bool) (*SiteClient, error) {
	l := log.New("SiteClient", verbose)
	l.V("New")
	http, err := httpclient.NewClient(url, skipVerify, verbose)
	if err != nil {
		return nil, l.E(err)
	}
	result := SiteClient{
		url:      url,
		siteName: siteName,
		username: username,
		password: password,
		http:     http,
		l:        l,
		siteId:   "Default",
		headers:  map[string]string{},
	}

	ai, err := result.GetApiInfo()
	if err != nil {
		return nil, result.l.E(err)
	}

	if ai.OmadacId == "" {
		return nil, result.l.E("Couldn't optain Omada ID.")
	}

	result.omadaId = ai.OmadacId

	if len(result.siteName) != 0 {
		cu, err := result.GetUserInfo()
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

// SiteClient finalizer
func (c *SiteClient) Close() {
	c.l.V("Close")
	if err := c.EndSession(); err != nil {
		c.l.E(err)
	}
}

// Start a new session
//
// Will be called by api methodes if required
//
// Return:
//   error
func (c *SiteClient) StartSession() error {
	c.l.V("StartSession")
	if c.HasActiveSession() {
		return nil
	}

	bodyData := `{
		"username": "` + c.username + `",
		"password": "` + c.password + `"
	}`

	var result model.Login
	if err := c.http.PostD(c.getPath("login"), bodyData, c.headers, empty, &result); err != nil {
		return c.l.E(err)
	}

	if result.Token == "" {
		return c.l.E("Couldn't optain Logintoken.")
	}

	c.setToken(result.Token)

	return nil
}

// End current session
//
// Will be called by Close
//
// Return:
//   error
func (c *SiteClient) EndSession() error {
	c.l.V("EndSession")
	if !c.HasActiveSession() {
		return nil
	}

	if _, err := c.http.Post(c.getPath("logout"), "", c.headers, empty); err != nil {
		return c.l.E(err)
	}

	c.removeToken()

	return nil
}

// Determins if client has an active session
//
// Return:
//   sesion state
func (c *SiteClient) HasActiveSession() bool {
	if !c.hasToken() {
		return false
	}

	var result model.LoginStatus
	if err := c.http.GetD(c.getPath("loginStatus"), "", c.headers, empty, &result); err != nil {
		c.removeToken()
		c.l.E(err)
		return false
	}

	if !result.Login {
		c.removeToken()
	}

	return result.Login
}

// Get API information
//
// Return:
//   API information
//   error
func (c *SiteClient) GetApiInfo() (*model.ApiInfo, error) {
	c.l.V("ApiInfo")
	var result model.ApiInfo
	if err := c.http.GetD("/api/info", "", c.headers, empty, &result); err != nil {
		return nil, c.l.E(err)
	}

	c.l.ReturnSuccess()

	return &result, nil
}

// Get user information for current session
//
// Return:
//   User information
//   error
func (c *SiteClient) GetUserInfo() (*model.UsersCurrent, error) {
	c.l.V("GetUserInfo")
	if err := c.ensureLoggedIn(); err != nil {
		return nil, c.l.E(err)
	}

	var result model.UsersCurrent
	if err := c.http.GetD(c.getPath("users/current"), "", c.headers, empty, &result); err != nil {
		return nil, c.l.E(err)
	}

	c.l.Return(result)
	return &result, nil
}

// Get list of devices
//
// Parameters:
//   detailed : get detailed device information
//
// Return:
//   Devices list
//   error
func (c *SiteClient) GetDevices(detailed bool) (*[]model.Device, error) {
	c.l.V("Devices")
	if err := c.ensureLoggedIn(); err != nil {
		return nil, c.l.E(err)
	}

	var devices []model.Device
	if err := c.http.GetD(c.getSitesPath("devices"), "", c.headers, empty, &devices); err != nil {
		return nil, c.l.E(err)
	}

	result := []model.Device{}
	for _, d := range devices {
		if d.Site == c.siteId {
			if detailed {
				if err := c.GetDeviceDetails(&d); err != nil {
					return nil, c.l.E(err)
				}
			}
			result = append(result, d)
		}
	}

	c.l.Return(result)

	return &result, nil
}

// Get enhanced information for a provided device and enhance the struct by it
//
// Parameters:
//   device : Device to enhance(Type and Mac have to be provided as minimal information)
//
// Return:
//   error
func (c *SiteClient) GetDeviceDetails(device *model.Device) error {
	c.l.V("GetDeviceDetails")
	if err := c.ensureLoggedIn(); err != nil {
		return c.l.E(err)
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
		return c.l.E("Unknown device type: " + device.Type)
	}

	if err := c.http.GetD(c.getSitesPath(dtype+"/"+device.Mac), "", c.headers, empty, device); err != nil {
		return c.l.E(err)
	}

	c.l.Return(*device)

	return nil
}

// Get active clients
//
// Parameters:
//   detailed : get detailed device information
//
// Return:
//   Client list
//   error
func (c *SiteClient) GetClients(detailed bool) (*[]model.Client, error) {
	c.l.V("GetClients")
	if err := c.ensureLoggedIn(); err != nil {
		return nil, c.l.E(err)
	}

	clients := []model.Client{}
	page := 1
	total := 1
	params := map[string]string{
		"currentPageSize": "10",
	}

	for len(clients) < total {
		params["currentPage"] = fmt.Sprint(page)

		var hres model.Clients
		if err := c.http.GetD(c.getSitesPath("clients"), "", c.headers, params, &hres); err != nil {
			return nil, c.l.E(err)
		}
		clients = append(clients, hres.Data...)
		page = hres.CurrentPage + 1
		total = hres.TotalRows
	}

	result := []model.Client{}
	for _, d := range clients {
		if detailed {
			if err := c.GetClientDetails(&d); err != nil {
				return nil, c.l.E(err)
			}
		}
		result = append(result, d)

	}

	c.l.Return(result)

	return &result, nil
}

// Get enhanced information for a provided cliend and enhance the struct by it
//
// Parameters:
//   client : Client to enhance(Type and Mac have to be provided as minimal information)
//
// Return:
//   error
func (c *SiteClient) GetClientDetails(client *model.Client) error {
	c.l.V("GetClientDetails")
	if err := c.ensureLoggedIn(); err != nil {
		return c.l.E(err)
	}

	if err := c.http.GetD(c.getSitesPath("clients/"+client.Mac), "", c.headers, empty, client); err != nil {
		return c.l.E(err)
	}

	c.l.Return(*client)

	return nil
}

// Get networks
//
// Return:
//   Network list
//   error
func (c *SiteClient) GetNetworks() (*[]model.Network, error) {
	c.l.V("GetNetworks")
	if err := c.ensureLoggedIn(); err != nil {
		return nil, c.l.E(err)
	}

	networks := []model.Network{}
	page := 1
	total := 1
	params := map[string]string{
		"currentPageSize": "10",
	}

	for len(networks) < total {
		params["currentPage"] = fmt.Sprint(page)

		var hres model.Networks
		if err := c.http.GetD(c.getSitesPath("setting/lan/networks"), "", c.headers, params, &hres); err != nil {
			return nil, c.l.E(err)
		}
		networks = append(networks, hres.Data...)
		page = hres.CurrentPage + 1
		total = hres.TotalRows
	}

	result := []model.Network{}
	for _, d := range networks {
		if d.Site == c.siteId {
			result = append(result, d)
		}
	}

	c.l.Return(result)

	return &result, nil
}

func (c *SiteClient) getPath(endPoint string) string {
	return "/" + c.omadaId + "/api/v2/" + endPoint
}

func (c *SiteClient) getSitesPath(endPoint string) string {
	return c.getPath("sites/" + c.siteId + "/" + endPoint)
}

func (c *SiteClient) hasToken() bool {
	_, ok := c.headers[tokenKey]
	return ok
}

func (c *SiteClient) setToken(token string) {
	c.headers[tokenKey] = token
}

func (c *SiteClient) removeToken() {
	delete(c.headers, tokenKey)
}

func (c *SiteClient) ensureLoggedIn() error {
	if c.HasActiveSession() {
		return nil
	}

	return c.StartSession()
}
