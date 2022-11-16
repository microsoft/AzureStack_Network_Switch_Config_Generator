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
	Switch         SwitchType        `json:"Switch"`
	SwitchPeers    []SwitchType      `json:"SwitchPeers"`
	SwitchBMC      []SwitchType      `json:"SwitchBMC"`
	SwitchUplink   []SwitchType      `json:"Uplinks"`
	SwitchDownlink []SwitchType      `json:"Downlinks"`
	GlobalSetting  GlobalSettingType `json:"GlobalSetting"`
	Vlans          []VlanType        `json:"Vlans"`
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
		SwitchSVI  bool       `json:"SwitchSVI"`
		Cidr       int        `json:"Cidr"`
		Assignment []IPv4Unit `json:"Assignment"`
	} `json:"IPv4"`
	IPv6 struct {
	} `json:"IPv6"`
}

type IPv4Unit struct {
	Name string `json:"Name"`
	IP   string `json:"IP"`
}
