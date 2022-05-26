package model

type Device struct {
	Type            string `json:"type"`
	Mac             string `json:"mac"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	ModelVersion    string `json:"modelVersion"`
	CompoundModel   string `json:"compoundModel"`
	ShowModel       string `json:"showModel"`
	FirmwareVersion string `json:"firmwareVersion"`
	Version         string `json:"version"`
	HwVersion       string `json:"hwVersion"`
	Status          int    `json:"status"`
	StatusCategory  int    `json:"statusCategory"`
	Site            string `json:"site"`
	Compatible      int    `json:"compatible"`
	Sn              string `json:"sn"`
	PortNum         int    `json:"portNum"`
	LedSetting      int    `json:"ledSetting"`
	SnmpSeting      struct {
		Location string `json:"location"`
		Contact  string `json:"contact"`
	} `json:"snmpSeting"`
	IptvSetting struct {
		IgmpEnable  bool   `json:"igmpEnable"`
		IgmpVersion string `json:"igmpVersion"`
	} `json:"iptvSetting"`
	HwOffloadEnable bool   `json:"hwOffloadEnable"`
	LldpEnable      bool   `json:"lldpEnable"`
	EchoServer      string `json:"echoServer"`
	IP              string `json:"ip"`
	Uptime          string `json:"uptime"`
	UptimeLong      int    `json:"uptimeLong"`
	CPUUtil         int    `json:"cpuUtil"`
	MemUtil         int    `json:"memUtil"`
	LastSeen        int64  `json:"lastSeen"`
	PortStats       []struct {
		Port              int    `json:"port"`
		Name              string `json:"name"`
		Type              int    `json:"type"`
		Mode              int    `json:"mode"`
		Status            int    `json:"status"`
		Rx                int    `json:"rx"`
		RxPkt             int    `json:"rxPkt"`
		RxPktRate         int    `json:"rxPktRate"`
		RxRate            int    `json:"rxRate"`
		Tx                int    `json:"tx"`
		TxPkt             int    `json:"txPkt"`
		TxPktRate         int    `json:"txPktRate"`
		TxRate            int    `json:"txRate"`
		InternetState     int    `json:"internetState,omitempty"`
		IP                string `json:"ip,omitempty"`
		Speed             int    `json:"speed,omitempty"`
		Duplex            int    `json:"duplex,omitempty"`
		Proto             string `json:"proto,omitempty"`
		WanIpv6Comptent   int    `json:"wanIpv6Comptent,omitempty"`
		WanPortIpv6Config struct {
			Enable        int    `json:"enable"`
			Addr          string `json:"addr"`
			Gateway       string `json:"gateway"`
			PriDNS        string `json:"priDns"`
			SndDNS        string `json:"sndDns"`
			InternetState int    `json:"internetState"`
		} `json:"wanPortIpv6Config,omitempty"`
		WanPortIpv4Config struct {
			IP       string `json:"ip"`
			Gateway  string `json:"gateway"`
			Gateway2 string `json:"gateway2"`
			PriDNS   string `json:"priDns"`
			SndDNS   string `json:"sndDns"`
			PriDNS2  string `json:"priDns2"`
			SndDNS2  string `json:"sndDns2"`
		} `json:"wanPortIpv4Config,omitempty"`
	} `json:"portStats"`
	LanClientStats []struct {
		LanName           string `json:"lanName"`
		Vlan              int    `json:"vlan"`
		IP                string `json:"ip"`
		Rx                int    `json:"rx"`
		Tx                int    `json:"tx"`
		ClientNum         int    `json:"clientNum"`
		LanPortIpv6Config struct {
			Addr string `json:"addr"`
		} `json:"lanPortIpv6Config"`
	} `json:"lanClientStats"`
	NeedUpgrade     bool  `json:"needUpgrade"`
	Download        int64 `json:"download"`
	Upload          int64 `json:"upload"`
	NetworkComptent int   `json:"networkComptent"`
}
