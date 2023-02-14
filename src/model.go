package main

type InputType struct {
	Version     string    `json:"Version"`
	Description string    `json:"Description"`
	InputData   InputData `json:"InputData"`
}

type InputData struct {
	Cloud            []CloudType  `json:"Cloud"`
	Switches         []SwitchType `json:"Switches"`
	SwitchUplink     string       `json:"SwitchUplink"`
	HostConnectivity string       `json:"HostConnectivity"`
	Supernets        []Supernet   `json:"Supernets"`
	Setting          struct {
		TimeServer   []string `json:"TimeServer"`
		SyslogServer []string `json:"SyslogServer"`
		DNSForwarder []string `json:"DNSForwarder"`
	}
}

type OutputType struct {
	Switch         SwitchType                 `json:"Switch"`
	SwitchPeers    []SwitchType               `json:"SwitchPeers"`
	SwitchBMC      []SwitchType               `json:"SwitchBMC"`
	SwitchUplink   []SwitchType               `json:"Uplinks"`
	SwitchDownlink []SwitchType               `json:"Downlinks"`
	GlobalSetting  GlobalSettingType          `json:"GlobalSetting"`
	Vlans          []VlanType                 `json:"Vlans"`
	L3Interfaces   map[string]L3IntfType      `json:"L3Interfaces"`
	PortChannel    map[string]PortChannelType `json:"PortChannel"`
	Ports          []PortType                 `json:"Ports"`
}

type GlobalSettingType struct {
	Username     string   `json:"Username"`
	Password     string   `json:"Password"`
	TimeServer   []string `json:"TimeServer"`
	SyslogServer []string `json:"SyslogServer"`
	DNSForwarder []string `json:"DNSForwarder"`
	OOB          string   `json:"OOB"`
}

type VlanType struct {
	VlanName      string `json:"VlanName"`
	VlanID        int    `json:"VlanID"`
	GroupID       string `json:"GroupID"`
	IPAddress     string `json:"IPAddress"`
	Mtu           int    `json:"MTU"`
	VIPAddress    string `json:"VIPAddress"`
	VIPPriorityId int    `json:"VIPPriorityId"`
	Shutdown      bool   `json:"Shutdown"`
}

type PortChannelType struct {
	Description   string `json:"Description"`
	Function      string `json:"Function"`
	UntagVlan     int    `json:"UntagVlan"`
	TagVlans      int    `json:"TagVlans"`
	IPAddress     string `json:"IPAddress"`
	PortChannelID string `json:"PortChannelID"`
	VPC           string `json:"VPC"`
	Shutdown      bool   `json:"Shutdown"`
}

type L3IntfType struct {
	Function  string `json:"Function"`
	IPAddress string `json:"IPAddress"`
	Cidr      int    `json:"Cidr"`
	Mtu       int    `json:"MTU"`
	Shutdown  bool   `json:"Shutdown"`
}

type CloudType struct {
	ID                      string   `json:"Id"`
	TimeServer              []string `json:"TimeServer"`
	SyslogServerIPv4Address string   `json:"SyslogServerIPv4Address"`
	DNSForwarder            []string `json:"DNSForwarder"`
}

type SwitchType struct {
	Make     string `json:"Make"`
	Model    string `json:"Model"`
	Type     string `json:"Type"`
	Hostname string `json:"Hostname"`
	Asn      int    `json:"ASN"`
	Firmware string `json:"Firmware"`
}

type Supernet struct {
	GroupID     string `json:"GroupID"`
	Name        string `json:"Name"`
	VlanID      int    `json:"VLANID"`
	Description string `json:"Description"`
	Shutdown    bool   `json:"Shutdown"`
	IPv4        struct {
		Name        string     `json:"Name"`
		NetworkType string     `json:"NetworkType"`
		SwitchSVI   bool       `json:"SwitchSVI"`
		Cidr        int        `json:"Cidr"`
		Assignment  []IPv4Unit `json:"Assignment"`
	} `json:"IPv4"`
	IPv6 struct {
	} `json:"IPv6"`
}

type IPv4Unit struct {
	Name string `json:"Name"`
	IP   string `json:"IP"`
}

type PortJson struct {
	Model           string `json:"Model"`
	Type            string `json:"Type"`
	SupportNoBmc    bool   `json:"SupportNoBmc"`
	FirmwareVersion string `json:"FirmwareVersion"`
	Maxmtu          string `json:"MaxMTU"`
	Function        []struct {
		Function string   `json:"Function"`
		Port     []string `json:"Port"`
	} `json:"Function"`
	Port []struct {
		Port string `json:"Port"`
		Type string `json:"Type"`
		Idx  int    `json:"Idx"`
	} `json:"Port"`
}

type PortType struct {
	Port        string            `json:"Port"`
	Idx         int               `json:"Idx"`
	Type        string            `json:"Type"`
	Description string            `json:"Description"`
	Function    string            `json:"Function"`
	UntagVlan   int               `json:"UntagVlan"`
	TagVlans    []int             `json:"TagVlans"`
	IPAddress   string            `json:"IPAddress"`
	Mtu         int               `json:"MTU"`
	Shutdown    bool              `json:"Shutdown"`
	Others      map[string]string `json:"Others"`
}
