package main

type InputType struct {
	Version     string    `json:"Version"`
	Description string    `json:"Description"`
	InputData   InputData `json:"InputData"`
}

type InputData struct {
	Cloud             []CloudType  `json:"Cloud"`
	Switches          []SwitchType `json:"Switches"`
	SwitchUplink      string       `json:"SwitchUplink"`
	DeploymentPattern string       `json:"DeploymentPattern"`
	HostConnectivity  string       `json:"HostConnectivity"`
	Supernets         []Supernet   `json:"Supernets"`
	Setting           struct {
		TimeServer   []string `json:"TimeServer"`
		SyslogServer []string `json:"SyslogServer"`
		DNSForwarder []string `json:"DNSForwarder"`
	}
}

type OutputType struct {
	Switch            SwitchType                 `json:"Switch,omitempty"`
	DeploymentPattern string                     `json:"DeploymentPattern"`
	SwitchPeer        []SwitchType               `json:"SwitchPeer,omitempty"`
	SwitchBMC         []SwitchType               `json:"SwitchBMC,omitempty"`
	SwitchUplink      []SwitchType               `json:"SwitchUplink,omitempty"`
	SwitchDownlink    []SwitchType               `json:"SwitchDownlink,omitempty"`
	GlobalSetting     GlobalSettingType          `json:"GlobalSetting,omitempty"`
	Vlans             []VlanType                 `json:"Vlans,omitempty"`
	L3Interfaces      map[string]L3IntfType      `json:"L3Interfaces,omitempty"`
	PortChannel       map[string]PortChannelType `json:"PortChannel,omitempty"`
	Ports             []PortType                 `json:"Ports,omitempty"`
	PortGroup         []PortGroupType            `json:"PortGroup,omitempty"`
	Routing           RoutingType                `json:"Routing,omitempty"`
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
	GroupName     string `json:"GroupName"`
	IPAddress     string `json:"IPAddress"`
	Cidr          int    `json:"Cidr"`
	Subnet        string `json:"Subnet"`
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
	Function     string `json:"Function"`
	Description  string `json:"Description"`
	IPAddress    string `json:"IPAddress"`
	Cidr         int    `json:"Cidr"`
	NbrIPAddress string `json:"NbrIPAddress,omitempty"`
	Subnet       string `json:"Subnet"`
	Mtu          int    `json:"MTU"`
	Shutdown     bool   `json:"Shutdown"`
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
	GroupName   string `json:"GroupName"`
	Description string `json:"Description"`
	Shutdown    bool   `json:"Shutdown"`
	IPv4        struct {
		Name        string     `json:"Name"`
		VlanID      int        `json:"VLANID"`
		NetworkType string     `json:"NetworkType"`
		SwitchSVI   bool       `json:"SwitchSVI"`
		Cidr        int        `json:"Cidr"`
		Subnet      string     `json:"Subnet"`
		Gateway     string     `json:"Gateway"`
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
	VlanGroup map[string][]string `json:"VlanGroup"`
	Port      []struct {
		Port      string `json:"Port"`
		Type      string `json:"Type"`
		Idx       int    `json:"Idx"`
		PortGroup string `json:"PortGroup,omitempty"`
		Mode      string `json:"Mode,omitempty"`
	} `json:"Port"`
}

type PortType struct {
	Port        string            `json:"Port"`
	Idx         int               `json:"Idx"`
	Type        string            `json:"Type"`
	Description string            `json:"Description"`
	Function    string            `json:"Function"`
	UntagVlan   int               `json:"UntagVlan,omitempty"`
	TagVlans    []int             `json:"TagVlans,omitempty"`
	IPAddress   string            `json:"IPAddress,omitempty"`
	Mtu         int               `json:"MTU"`
	Shutdown    bool              `json:"Shutdown"`
	Others      map[string]string `json:"Others,omitempty"`
	Mode        string            `json:"Mode,omitempty"`
	PortGroup   string            `json:"PortGroup,omitempty"`
}

type RoutingType struct {
	BGP        BGPType          `json:"BGP,omitempty"`
	Static     []StaticType     `json:"Static,omitempty"`
	PrefixList []PrefixListType `json:"PrefixList,omitempty"`
}

type BGPType struct {
	BGPAsn          int                `json:"BGPAsn"`
	RouterID        string             `json:"RouterID"`
	IPv4Network     []string           `json:"IPv4Network"`
	RouteMap        []RouteMapType     `json:"RouteMap,omitempty"`
	IPv4Neighbor    []IPv4NeighborType `json:"IPv4Neighbor"`
	TemplateNeigbor []IPv4NeighborType `json:"TemplateNeigbor,omitempty"`
}

type IPv4NeighborType struct {
	SwitchRelation    string `json:"SwitchRelation"`
	Description       string `json:"Description"`
	NeighborAsn       int    `json:"NeighborAsn"`
	NeighborIPAddress string `json:"NeighborIPAddress"`
	PrefixList        []struct {
		Name      string `json:"Name"`
		Direction string `json:"Direction"`
	} `json:"PrefixList,omitempty"`
	RemovePrivateAS bool   `json:"RemovePrivateAS,omitempty"`
	Shutdown        bool   `json:"Shutdown"`
	NbrPassword     string `json:"NbrPassword,omitempty"`
	UpdateSource    string `json:"UpdateSource,omitempty"`
	LocalAS         string `json:"LocalAS,omitempty"`
	EBGPMultiHop    int    `json:"EBGPMultiHop,omitempty"`
}

type StaticType struct {
	Network string `json:"Network"`
	NextHop string `json:"NextHop"`
	Name    string `json:"Name"`
}

type PrefixListType struct {
	Name   string `json:"Name"`
	Config []struct {
		Idx         int    `json:"Idx"`
		Action      string `json:"Action"`
		Description string `json:"Description"`
		Network     string `json:"Network"`
		Operation   string `json:"Operation"`
		Prefix      int    `json:"Prefix"`
	} `json:"Config"`
}

type RouteMapType struct {
	Name   string `json:"Name"`
	Action string `json:"Action"`
	Seq    int    `json:"Seq"`
	Rules  []struct {
		PrefixList string `json:"PrefixList"`
	} `json:"Rules,omitempty"`
}

// For Dell Port-Group
type PortGroupType struct {
	PortGroup string `json:"PortGroup,omitempty"`
	Mode      string `json:"Mode,omitempty"`
	Type      string `json:"Type,omitempty"`
	Idx       int    `json:"Idx,omitempty"`
}
