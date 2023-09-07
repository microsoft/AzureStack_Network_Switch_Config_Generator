package main

type InputType struct {
	Version     string    `yaml:"Version"`
	Description string    `yaml:"Description"`
	InputData   InputData `yaml:"InputData"`
}

type InputData struct {
	Cloud             []CloudType  `yaml:"Cloud"`
	Switches          []SwitchType `yaml:"Switches"`
	SwitchUplink      string       `yaml:"SwitchUplink"`
	DeploymentPattern string       `yaml:"DeploymentPattern"`
	HostConnectivity  string       `yaml:"HostConnectivity"`
	Supernets         []Supernet   `yaml:"Supernets"`
	Setting           struct {
		TimeServer   []string `yaml:"TimeServer"`
		SyslogServer []string `yaml:"SyslogServer"`
		DNSForwarder []string `yaml:"DNSForwarder"`
	} `yaml:"Setting"`
	WANSIM WANSIMType `yaml:"WANSIM,omitempty"`
}

type OutputType struct {
	ToolBuildVersion  string                     `yaml:"ToolBuildVersion"`
	Switch            SwitchType                 `yaml:"Switch,omitempty"`
	DeploymentPattern string                     `yaml:"DeploymentPattern"`
	SwitchPeer        []SwitchType               `yaml:"SwitchPeer,omitempty"`
	SwitchBMC         []SwitchType               `yaml:"SwitchBMC,omitempty"`
	SwitchUplink      []SwitchType               `yaml:"SwitchUplink,omitempty"`
	SwitchDownlink    []SwitchType               `yaml:"SwitchDownlink,omitempty"`
	GlobalSetting     GlobalSettingType          `yaml:"GlobalSetting,omitempty"`
	Vlans             []VlanType                 `yaml:"Vlans,omitempty"`
	L3Interfaces      map[string]L3IntfType      `yaml:"L3Interfaces,omitempty"`
	PortChannel       map[string]PortChannelType `yaml:"PortChannel,omitempty"`
	Ports             []PortType                 `yaml:"Ports,omitempty"`
	PortGroup         []PortGroupType            `yaml:"PortGroup,omitempty"`
	Routing           RoutingType                `yaml:"Routing,omitempty"`
	WANSIM            WANSIMType                 `yaml:"WANSIM,omitempty"`
}

type GlobalSettingType struct {
	Username     string   `yaml:"Username"`
	Password     string   `yaml:"Password"`
	TimeServer   []string `yaml:"TimeServer"`
	SyslogServer []string `yaml:"SyslogServer"`
	DNSForwarder []string `yaml:"DNSForwarder"`
	DHCPInfra    []string `yaml:"DHCPInfra"`
	DHCPTenant   []string `yaml:"DHCPTenant"`
	OOB          string   `yaml:"OOB"`
}

type VlanType struct {
	GroupName      string            `yaml:"GroupName"`
	VlanName       string            `yaml:"VlanName"`
	VlanID         int               `yaml:"VlanID"`
	VirtualGroupID int               `yaml:"VirtualGroupID,omitempty"`
	IPAddress      string            `yaml:"IPAddress"`
	Cidr           int               `yaml:"Cidr"`
	Subnet         string            `yaml:"Subnet"`
	Mtu            int               `yaml:"MTU"`
	VIPAddress     string            `yaml:"VIPAddress,omitempty"`
	VIPPriorityId  int               `yaml:"VIPPriorityId,omitempty"`
	Shutdown       bool              `yaml:"Shutdown"`
	Others         map[string]string `yaml:"Others,omitempty"`
}

type PortChannelType struct {
	Description   string `yaml:"Description"`
	Function      string `yaml:"Function"`
	UntagVlan     int    `yaml:"UntagVlan"`
	TagVlans      int    `yaml:"TagVlans"`
	IPAddress     string `yaml:"IPAddress"`
	PortChannelID string `yaml:"PortChannelID"`
	VPC           string `yaml:"VPC"`
	Shutdown      bool   `yaml:"Shutdown"`
}

type L3IntfType struct {
	Name         string `yaml:"Name"`
	Function     string `yaml:"Function"`
	Description  string `yaml:"Description"`
	IPAddress    string `yaml:"IPAddress"`
	Cidr         int    `yaml:"Cidr"`
	NbrIPAddress string `yaml:"NbrIPAddress,omitempty"`
	Subnet       string `yaml:"Subnet"`
	Mtu          int    `yaml:"MTU"`
	Shutdown     bool   `yaml:"Shutdown"`
}

type CloudType struct {
	ID                      string   `yaml:"Id"`
	TimeServer              []string `yaml:"TimeServer"`
	SyslogServerIPv4Address string   `yaml:"SyslogServerIPv4Address"`
	DNSForwarder            []string `yaml:"DNSForwarder"`
}

type SwitchType struct {
	Make     string `yaml:"Make"`
	Model    string `yaml:"Model"`
	Type     string `yaml:"Type"`
	Hostname string `yaml:"Hostname"`
	Asn      int    `yaml:"ASN"`
	Firmware string `yaml:"Firmware"`
}

type Supernet struct {
	GroupName   string `yaml:"GroupName"`
	Description string `yaml:"Description"`
	Shutdown    bool   `yaml:"Shutdown"`
	IPv4        struct {
		Name        string     `yaml:"Name"`
		VlanID      int        `yaml:"VLANID"`
		NetworkType string     `yaml:"NetworkType"`
		SwitchSVI   bool       `yaml:"SwitchSVI"`
		Cidr        int        `yaml:"Cidr"`
		Subnet      string     `yaml:"Subnet"`
		Gateway     string     `yaml:"Gateway"`
		Assignment  []IPv4Unit `yaml:"Assignment"`
	} `yaml:"IPv4"`
	IPv6 struct {
	} `yaml:"IPv6"`
}

type IPv4Unit struct {
	Name      string `yaml:"Name,omitempty"`
	IP        string `yaml:"IP,omitempty"`
	IPNetwork string `yaml:"IPNetwork,omitempty"`
	LocalIP   string `yaml:"LocalIP,omitempty"`
	RemoteIP  string `yaml:"RemoteIP,omitempty"`
}

type PortJson struct {
	Model           string `yaml:"Model"`
	Type            string `yaml:"Type"`
	SupportNoBmc    bool   `yaml:"SupportNoBmc"`
	FirmwareVersion string `yaml:"FirmwareVersion"`
	Maxmtu          string `yaml:"MaxMTU"`
	Function        []struct {
		Function string   `yaml:"Function"`
		Port     []string `yaml:"Port"`
	} `yaml:"Function"`
	VlanGroup map[string][]string `yaml:"VlanGroup"`
	Port      []struct {
		Port      string `yaml:"Port"`
		Type      string `yaml:"Type"`
		Idx       int    `yaml:"Idx"`
		PortGroup string `yaml:"PortGroup,omitempty"`
		Mode      string `yaml:"Mode,omitempty"`
	} `yaml:"Port"`
}

type PortType struct {
	Port        string            `yaml:"Port"`
	Idx         int               `yaml:"Idx"`
	Type        string            `yaml:"Type"`
	Description string            `yaml:"Description"`
	Function    string            `yaml:"Function"`
	UntagVlan   int               `yaml:"UntagVlan,omitempty"`
	TagVlans    []int             `yaml:"TagVlans,omitempty"`
	IPAddress   string            `yaml:"IPAddress,omitempty"`
	Mtu         int               `yaml:"MTU"`
	Shutdown    bool              `yaml:"Shutdown"`
	Others      map[string]string `yaml:"Others,omitempty"`
	Mode        string            `yaml:"Mode,omitempty"`
	PortGroup   string            `yaml:"PortGroup,omitempty"`
}

type RoutingType struct {
	BGP        BGPType          `yaml:"BGP,omitempty"`
	Static     []StaticType     `yaml:"Static,omitempty"`
	PrefixList []PrefixListType `yaml:"PrefixList,omitempty"`
}

type BGPType struct {
	BGPAsn          int                `yaml:"BGPAsn"`
	RouterID        string             `yaml:"RouterID"`
	IPv4Network     []string           `yaml:"IPv4Network"`
	RouteMap        []RouteMapType     `yaml:"RouteMap,omitempty"`
	IPv4Neighbor    []IPv4NeighborType `yaml:"IPv4Neighbor"`
	TemplateNeigbor []IPv4NeighborType `yaml:"TemplateNeigbor,omitempty"`
}

type IPv4NeighborType struct {
	SwitchRelation    string `yaml:"SwitchRelation"`
	Description       string `yaml:"Description"`
	NeighborAsn       int    `yaml:"NeighborAsn"`
	NeighborIPAddress string `yaml:"NeighborIPAddress"`
	PrefixList        []struct {
		Name      string `yaml:"Name"`
		Direction string `yaml:"Direction"`
	} `yaml:"PrefixList,omitempty"`
	RemovePrivateAS bool   `yaml:"RemovePrivateAS,omitempty"`
	Shutdown        bool   `yaml:"Shutdown"`
	NbrPassword     string `yaml:"NbrPassword,omitempty"`
	UpdateSource    string `yaml:"UpdateSource,omitempty"`
	LocalAS         string `yaml:"LocalAS,omitempty"`
	EBGPMultiHop    int    `yaml:"EBGPMultiHop,omitempty"`
}

type StaticType struct {
	Network string `yaml:"Network"`
	NextHop string `yaml:"NextHop"`
	Name    string `yaml:"Name"`
}

type PrefixListType struct {
	Name   string `yaml:"Name"`
	Config []struct {
		Idx         int    `yaml:"Idx"`
		Action      string `yaml:"Action"`
		Description string `yaml:"Description"`
		Network     string `yaml:"Network"`
		Operation   string `yaml:"Operation"`
		Prefix      int    `yaml:"Prefix"`
	} `yaml:"Config"`
}

type RouteMapType struct {
	Name   string `yaml:"Name"`
	Action string `yaml:"Action"`
	Seq    int    `yaml:"Seq"`
	Rules  []struct {
		PrefixList string `yaml:"PrefixList"`
	} `yaml:"Rules,omitempty"`
}

// For Dell Port-Group
type PortGroupType struct {
	PortGroup string `yaml:"PortGroup,omitempty"`
	Mode      string `yaml:"Mode,omitempty"`
	Type      string `yaml:"Type,omitempty"`
	Idx       int    `yaml:"Idx,omitempty"`
}

type WANSIMType struct {
	Hostname string   `yaml:"Hostname"`
	Loopback IPv4Unit `yaml:"Loopback"`
	GRE1     IPv4Unit `yaml:"GRE1"`
	GRE2     IPv4Unit `yaml:"GRE2"`
	BGP      struct {
		ASN    int    `yaml:"ASN"`
		NbrIP  string `yaml:"NbrIP"`
		NbrASN int    `yaml:"NbrASN"`
	} `yaml:"BGP"`
	RerouteNetworks []string `yaml:"RerouteNetworks"`
}
