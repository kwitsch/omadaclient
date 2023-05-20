package model

type DeviceSnmpSeting struct {
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

type DeviceIptvSetting struct {
	IgmpEnable  bool   `json:"igmpEnable"`
	IgmpVersion string `json:"igmpVersion"`
}

type DevicePortStatWanPortIpv6Config struct {
	Enable        uint64 `json:"enable"`
	Addr          string `json:"addr"`
	Gateway       string `json:"gateway"`
	PriDNS        string `json:"priDns"`
	SndDNS        string `json:"sndDns"`
	InternetState uint64 `json:"internetState"`
}

type DevicePortStatWanPortIpv4Config struct {
	IP       string `json:"ip"`
	Gateway  string `json:"gateway"`
	Gateway2 string `json:"gateway2"`
	PriDNS   string `json:"priDns"`
	SndDNS   string `json:"sndDns"`
	PriDNS2  string `json:"priDns2"`
	SndDNS2  string `json:"sndDns2"`
}

type DevicePortStat struct {
	Port              uint64                          `json:"port"`
	Name              string                          `json:"name"`
	Type              uint64                          `json:"type"`
	Mode              int64                           `json:"mode"`
	Status            uint64                          `json:"status"`
	Rx                uint64                          `json:"rx"`
	RxPkt             uint64                          `json:"rxPkt"`
	RxPktRate         uint64                          `json:"rxPktRate"`
	RxRate            uint64                          `json:"rxRate"`
	Tx                uint64                          `json:"tx"`
	TxPkt             uint64                          `json:"txPkt"`
	TxPktRate         uint64                          `json:"txPktRate"`
	TxRate            uint64                          `json:"txRate"`
	InternetState     uint64                          `json:"internetState,omitempty"`
	IP                string                          `json:"ip,omitempty"`
	Speed             uint64                          `json:"speed,omitempty"`
	Duplex            uint64                          `json:"duplex,omitempty"`
	Proto             string                          `json:"proto,omitempty"`
	WanIpv6Comptent   uint64                          `json:"wanIpv6Comptent,omitempty"`
	WanPortIpv6Config DevicePortStatWanPortIpv6Config `json:"wanPortIpv6Config,omitempty"`
	WanPortIpv4Config DevicePortStatWanPortIpv4Config `json:"wanPortIpv4Config,omitempty"`
}

type DeviceLanClientStatLanPortIpv6Config struct {
	Addr string `json:"addr"`
}

type DeviceLanClientStat struct {
	LanName           string                               `json:"lanName"`
	Vlan              uint64                               `json:"vlan"`
	IP                string                               `json:"ip"`
	Rx                uint64                               `json:"rx"`
	Tx                uint64                               `json:"tx"`
	ClientNum         uint64                               `json:"clientNum"`
	LanPortIpv6Config DeviceLanClientStatLanPortIpv6Config `json:"lanPortIpv6Config"`
}

type Device struct {
	Type            string                `json:"type"`
	Mac             string                `json:"mac"`
	Name            string                `json:"name"`
	Model           string                `json:"model"`
	ModelVersion    string                `json:"modelVersion"`
	CompoundModel   string                `json:"compoundModel"`
	ShowModel       string                `json:"showModel"`
	FirmwareVersion string                `json:"firmwareVersion"`
	Version         string                `json:"version"`
	HwVersion       string                `json:"hwVersion"`
	Status          uint64                `json:"status"`
	StatusCategory  uint64                `json:"statusCategory"`
	Site            string                `json:"site"`
	Compatible      uint64                `json:"compatible"`
	Sn              string                `json:"sn"`
	PortNum         uint64                `json:"portNum"`
	LedSetting      uint64                `json:"ledSetting"`
	SnmpSeting      DeviceSnmpSeting      `json:"snmpSeting"`
	IptvSetting     DeviceIptvSetting     `json:"iptvSetting"`
	HwOffloadEnable bool                  `json:"hwOffloadEnable"`
	LldpEnable      bool                  `json:"lldpEnable"`
	EchoServer      string                `json:"echoServer"`
	IP              string                `json:"ip"`
	Uptime          string                `json:"uptime"`
	UptimeLong      uint64                `json:"uptimeLong"`
	CPUUtil         uint64                `json:"cpuUtil"`
	MemUtil         uint64                `json:"memUtil"`
	LastSeen        uint64                `json:"lastSeen"`
	PortStats       []DevicePortStat      `json:"portStats"`
	LanClientStats  []DeviceLanClientStat `json:"lanClientStats"`
	NeedUpgrade     bool                  `json:"needUpgrade"`
	Download        uint64                `json:"download"`
	Upload          uint64                `json:"upload"`
	NetworkComptent uint64                `json:"networkComptent"`
}
