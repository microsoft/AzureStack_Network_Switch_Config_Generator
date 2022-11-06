package main

type InputType struct {
	Version     string `json:"Version"`
	Description string `json:"Description"`
	InputData   struct {
		Cloud            []CloudType    `json:"Cloud"`
		Switches         []SwitchesType `json:"Switches"`
		SwitchUplink     string         `json:"SwitchUplink"`
		HostConnectivity string         `json:"HostConnectivity"`
		Supernets        []Supernet     `json:"Supernets"`
	} `json:"InputData"`
}

type OutputType struct {
	Switch    SwitchesType   `json:"Switch"`
	Uplinks   []SwitchesType `json:"Uplinks"`
	Downlinks []SwitchesType `json:"Downlinks"`
	Vlans     []VlanType     `json:"Vlans"`
}

type VlanType struct {
	VlanName  string `json:"VlanName"`
	VlanID    int    `json:"VlanID"`
	Group     string `json:"Group"`
	IPAddress string `json:"IPAddress"`
	Mtu       int    `json:"MTU"`
	Vip       struct {
		PriorityId int    `json:"PriorityId"`
		VIPAddress string `json:"VIPAddress"`
	} `json:"VIP"`
	Shutdown bool `json:"Shutdown"`
}

type CloudType struct {
	ID                      string   `json:"Id"`
	TimeServer              []string `json:"TimeServer"`
	SyslogServerIPv4Address string   `json:"SyslogServerIPv4Address"`
	DNSForwarder            []string `json:"DNSForwarder"`
}

type SwitchesType struct {
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
	Description string `json:"Description"`
	IPv4        struct {
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
