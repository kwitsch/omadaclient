package model

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
	ConnectType    uint   `json:"connectType"`
	ConnectDevType string `json:"connectDevType"`
	Wireless       bool   `json:"wireless"`
	Ssid           string `json:"ssid,omitempty"`
	SignalLevel    uint   `json:"signalLevel,omitempty"`
	SignalRank     uint   `json:"signalRank,omitempty"`
	WifiMode       uint   `json:"wifiMode,omitempty"`
	ApName         string `json:"apName,omitempty"`
	ApMac          string `json:"apMac,omitempty"`
	RadioID        uint   `json:"radioId,omitempty"`
	Channel        uint   `json:"channel,omitempty"`
	RxRate         uint   `json:"rxRate,omitempty"`
	TxRate         uint   `json:"txRate,omitempty"`
	PowerSave      bool   `json:"powerSave,omitempty"`
	Rssi           int    `json:"rssi,omitempty"`
	Activity       uint   `json:"activity"`
	TrafficDown    uint64 `json:"trafficDown"`
	TrafficUp      uint64 `json:"trafficUp"`
	Uptime         uint64 `json:"uptime"`
	LastSeen       uint64 `json:"lastSeen"`
	AuthStatus     uint   `json:"authStatus"`
	Guest          bool   `json:"guest"`
	Active         bool   `json:"active"`
	Manager        bool   `json:"manager"`
	DownPacket     uint64 `json:"downPacket"`
	UpPacket       uint64 `json:"upPacket"`
	SwitchMac      string `json:"switchMac,omitempty"`
	SwitchName     string `json:"switchName,omitempty"`
	Vid            uint   `json:"vid,omitempty"`
	NetworkName    string `json:"networkName,omitempty"`
	Dot1XVlan      uint   `json:"dot1xVlan,omitempty"`
	Port           uint   `json:"port,omitempty"`
	IPSetting      struct {
		UseFixedAddr bool   `json:"useFixedAddr"`
		NetID        string `json:"netId"`
		IP           string `json:"ip"`
	} `json:"ipSetting"`
	RateLimit struct {
		RateLimitID string `json:"rateLimitId"`
		Enable      bool   `json:"enable"`
		UpEnable    bool   `json:"upEnable"`
		UpUnit      uint   `json:"upUnit"`
		UpLimit     uint   `json:"upLimit"`
		DownEnable  bool   `json:"downEnable"`
		DownUnit    uint   `json:"downUnit"`
		DownLimit   uint   `json:"downLimit"`
	} `json:"rateLimit"`
}
