package model

type Networks struct {
	Paged
	Data []Network `json:"data"`
}

type Network struct {
	ID                string   `json:"id"`
	Site              string   `json:"site"`
	Name              string   `json:"name"`
	Purpose           string   `json:"purpose"`
	InterfaceIds      []string `json:"interfaceIds"`
	Vlan              uint64   `json:"vlan"`
	GatewaySubnet     string   `json:"gatewaySubnet"`
	Domain            string   `json:"domain,omitempty"`
	IgmpSnoopEnable   bool     `json:"igmpSnoopEnable"`
	Portal            bool     `json:"portal"`
	AccessControlRule bool     `json:"accessControlRule"`
	RateLimit         bool     `json:"rateLimit"`
	AllLan            bool     `json:"allLan"`
	Primary           bool     `json:"primary"`
	DhcpGuard         struct {
		Enable   bool   `json:"enable"`
		DhcpSvr1 string `json:"dhcpSvr1"`
	} `json:"dhcpGuard,omitempty"`
	DhcpSettings struct {
		Enable       bool   `json:"enable"`
		IpaddrStart  string `json:"ipaddrStart"`
		IpaddrEnd    string `json:"ipaddrEnd"`
		IPRangeStart int64  `json:"ipRangeStart"`
		IPRangeEnd   int64  `json:"ipRangeEnd"`
		Dhcpns       string `json:"dhcpns"`
		PriDNS       string `json:"priDns"`
		SndDNS       string `json:"sndDns"`
		Leasetime    uint64 `json:"leasetime"`
		Option138    string `json:"option138"`
	} `json:"dhcpSettings,omitempty"`
	LanNetworkIpv6Config struct {
		Proto  string `json:"proto"`
		Enable int    `json:"enable"`
	} `json:"lanNetworkIpv6Config,omitempty"`
}
