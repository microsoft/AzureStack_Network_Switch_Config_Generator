package main

type InputType struct {
	Version   string                 `json:"Version"`
	Settings  map[string]interface{} `json:"Settings"`
	Devices   []DeviceType           `json:"Devices"`
	Supernets interface{}            `json:"Supernets"`
}

type DeviceType struct {
	Make                 string `json:"Make"`
	Type                 string `json:"Type"`
	Hostname             string `json:"Hostname"`
	Asn                  int    `json:"ASN"`
	Model                string `json:"Model"`
	Firmware             string `json:"Firmware"`
	GenerateDeviceConfig bool   `json:"GenerateDeviceConfig"`
	StaticRouting        bool   `json:"StaticRouting"`
	Username             string `json:"Username"`
	Password             string `json:"Password"`
}

type SupernetInputType struct {
	VlanID           int    `json:"VlanID"`
	Group            string `json:"Group"`
	Name             string `json:"Name"`
	Subnet           string `json:"Subnet"`
	Shutdown         bool   `json:"Shutdown"`
	SubnetAssignment []struct {
		Name         string                  `json:"Name"`
		Netmask      int                     `json:"Netmask"`
		IPSize       int                     `json:"IPSize"`
		IPAssignment []IPAssignmentInputItem `json:"IPAssignment"`
	} `json:"SubnetAssignment"`
}

type IPAssignmentInputItem struct {
	Name     string `json:"Name"`
	Position int    `json:"Position"`
}

type SupernetOutputType struct {
	VlanID       int                      `json:"VlanID"`
	Group        string                   `json:"Group"`
	Name         string                   `json:"Name"`
	Subnet       string                   `json:"Subnet"`
	Shutdown     bool                     `json:"Shutdown"`
	IPAssignment []IPAssignmentOutputItem `json:"IPAssignment"`
}

type IPAssignmentOutputItem struct {
	Name      string `json:"Name"`
	IPAddress string `json:"IPAddress"`
}

type InterfaceFrameworkType struct {
	Device struct {
		Type  string `json:"Type"`
		Make  string `json:"Make"`
		Model string `json:"Model"`
	} `json:"Device"`
	InterfaceName []struct {
		Speed int    `json:"Speed"`
		Name  string `json:"Name"`
	} `json:"InterfaceName"`
	Port []struct {
		ID          int               `json:"ID"`
		Port        string            `json:"Port"`
		Speed       int               `json:"Speed"`
		Description string            `json:"Description"`
		Mtu         int               `json:"MTU"`
		PortType    string            `json:"PortType"`
		Shutdown    bool              `json:"Shutdown,omitempty"`
		IPAddress   string            `json:"IPAddress"`
		UntagVlan   string            `json:"UntagVlan"`
		TagVlan     []string          `json:"TagVlan"`
		Others      map[string]string `json:"Others"`
	} `json:"Port"`
	Vlan []struct {
		Group             string `json:"Group"`
		IPAssignment      string `json:"IPAssignment"`
		Mtu               int    `json:"MTU"`
		Native            bool   `json:"Native"`
		VirtualAssignment struct {
			PriorityID   int    `json:"PriorityId"`
			IPAssignment string `json:"IPAssignment"`
		} `json:"VirtualAssignment"`
		IPHelper []string      `json:"IPHelper"`
		ACL      []interface{} `json:"ACL"`
		Shutdown bool          `json:"Shutdown,omitempty"`
	} `json:"VLAN"`
	Loopback    []LoopbackType `json:"Loopback"`
	PortChannel []struct {
		Name     string        `json:"Name"`
		ID       int           `json:"ID"`
		PortType string        `json:"PortType"`
		Settings []interface{} `json:"Settings"`
	} `json:"PortChannel"`
}

type OutputType struct {
	Device    DeviceType             `json:"Device"`
	Settings  map[string]interface{} `json:"Settings"`
	Port      []PortType             `json:"Port"`
	Vlan      []VlanType             `json:"Vlan"`
	Loopback  []LoopbackType         `json:"Loopback"`
	Routing   *RoutingType           `json:"Routing"`
	Supernets *[]SupernetOutputType  `json:"Supernets"`
}

type PortType struct {
	Port        string            `json:"Port"`
	PortName    string            `json:"PortName"`
	PortType    string            `json:"PortType"`
	Description string            `json:"Description"`
	Mtu         int               `json:"MTU"`
	Shutdown    bool              `json:"Shutdown"`
	IPAddress   string            `json:"IPAddress"`
	UntagVlan   int               `json:"UntagVlan"`
	TagVlan     string            `json:"TagVlan"`
	Others      map[string]string `json:"Others"`
}

type VlanType struct {
	VlanName  string `json:"VlanName"`
	VlanID    int    `json:"VlanID"`
	Group     string `json:"Group"`
	IPAddress string `json:"IPAddress"`
	Mtu       int    `json:"MTU"`
	Shutdown  bool   `json:"Shutdown"`
}

type LoopbackType struct {
	Description string `json:"Description"`
	IPAddress   string `json:"IPAddress"`
}

type RoutingType struct {
	Bgp           BGPType    `json:"BGP"`
	Static        StaticType `json:"Static"`
	RoutingPolicy struct {
		PrefixList []PrefixListType `json:"PrefixList"`
		RouteMap   []RouteMapType   `json:"RouteMap"`
	} `json:"RoutingPolicy"`
}

type BGPType struct {
	BGPAsn                 int      `json:"BGPAsn"`
	RouterID               string   `json:"RouterID"`
	IPv4Network            []string `json:"IPv4Network"`
	EnableDefaultOriginate bool     `json:"EnableDefaultOriginate"`
	RoutePrefix            struct {
		MaxiPrefix  int    `json:"MaxiPrefix"`
		ErrorAction string `json:"ErrorAction"`
	} `json:"RoutePrefix"`
	IPv4Neighbor []struct {
		Description       string `json:"Description"`
		EnablePassword    bool   `json:"EnablePassword"`
		NeighborAsn       string `json:"NeighborAsn"`
		NeighborIPAddress string `json:"NeighborIPAddress"`
		PrefixList        []struct {
			Name      string `json:"Name"`
			Direction string `json:"Direction"`
		} `json:"PrefixList"`
		UpdateSource string `json:"UpdateSource"`
		Shutdown     bool   `json:"Shutdown"`
	} `json:"IPv4Neighbor"`
	PrefixListName []string `json:"PrefixListName"`
}

type StaticType struct {
	PrefixListName []string            `json:"PrefixListName"`
	RouteMapName   []string            `json:"RouteMapName"`
	Network        []StaticNetworkType `json:"Network"`
}

type StaticNetworkType struct {
	DstIPAddress string
	NextHop      string
	Name         string
}

type PrefixListType struct {
	Name   string `json:"Name"`
	Config []struct {
		Name      string `json:"Name"`
		Action    string `json:"Action"`
		Supernet  string `json:"Supernet"`
		IPAddress string `json:"IPAddress"`
		Operation string `json:"Operation"`
		Prefix    int    `json:"Prefix"`
	} `json:"Config"`
}

type RouteMapType struct {
	Name   string `json:"Name"`
	Config []struct {
		Index          int      `json:"Index"`
		Name           string   `json:"Name"`
		Action         string   `json:"Action"`
		PrefixListName []string `json:"PrefixListName"`
	} `json:"Config"`
}
