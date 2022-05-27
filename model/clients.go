package model

import "strings"

type Clients struct {
	Paged
	Data       []Client `json:"data"`
	ClientStat struct {
		Total       int `json:"total"`
		Wireless    int `json:"wireless"`
		Wired       int `json:"wired"`
		Num2G       int `json:"num2g"`
		Num5G       int `json:"num5g"`
		Num5G2      int `json:"num5g2"`
		NumUser     int `json:"numUser"`
		NumGuest    int `json:"numGuest"`
		Num2GUser   int `json:"num2gUser"`
		Num5GUser   int `json:"num5gUser"`
		Num5G2User  int `json:"num5g2User"`
		Num2GGuest  int `json:"num2gGuest"`
		Num5GGuest  int `json:"num5gGuest"`
		Num5G2Guest int `json:"num5g2Guest"`
	} `json:"clientStat"`
}

type Client struct {
	Mac            string `json:"mac"`
	Name           string `json:"name"`
	HostName       string `json:"hostName,omitempty"`
	DeviceType     string `json:"deviceType"`
	IP             string `json:"ip"`
	ConnectType    int    `json:"connectType"`
	ConnectDevType string `json:"connectDevType"`
	Wireless       bool   `json:"wireless"`
	Ssid           string `json:"ssid,omitempty"`
	SignalLevel    int    `json:"signalLevel,omitempty"`
	SignalRank     int    `json:"signalRank,omitempty"`
	WifiMode       int    `json:"wifiMode,omitempty"`
	ApName         string `json:"apName,omitempty"`
	ApMac          string `json:"apMac,omitempty"`
	RadioID        int    `json:"radioId,omitempty"`
	Channel        int    `json:"channel,omitempty"`
	RxRate         int    `json:"rxRate,omitempty"`
	TxRate         int    `json:"txRate,omitempty"`
	PowerSave      bool   `json:"powerSave,omitempty"`
	Rssi           int    `json:"rssi,omitempty"`
	Activity       int    `json:"activity"`
	TrafficDown    int    `json:"trafficDown"`
	TrafficUp      int    `json:"trafficUp"`
	Uptime         int    `json:"uptime"`
	LastSeen       int64  `json:"lastSeen"`
	AuthStatus     int    `json:"authStatus"`
	Guest          bool   `json:"guest"`
	Active         bool   `json:"active"`
	Manager        bool   `json:"manager"`
	DownPacket     int    `json:"downPacket"`
	UpPacket       int    `json:"upPacket"`
	SwitchMac      string `json:"switchMac,omitempty"`
	SwitchName     string `json:"switchName,omitempty"`
	Vid            int    `json:"vid,omitempty"`
	NetworkName    string `json:"networkName,omitempty"`
	Dot1XVlan      int    `json:"dot1xVlan,omitempty"`
	Port           int    `json:"port,omitempty"`
}

func (c *Client) GetCleanName() string {
	return strings.ReplaceAll(c.Name, " ", "")
}
