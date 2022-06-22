package main

type InputType struct {
	External []struct {
		Type string   `json:"Type"`
		IP   []string `json:"IP"`
	} `json:"External"`
	Device  []DeviceType `json:"Device"`
	Network interface{}  `json:"Network"`
}

type DeviceType struct {
	Make                 string `json:"Make"`
	Type                 string `json:"Type"`
	Asn                  int    `json:"ASN"`
	Hostname             string `json:"Hostname"`
	Model                string `json:"Model"`
	Firmware             string `json:"Firmware"`
	GenerateDeviceConfig bool   `json:"GenerateDeviceConfig"`
	Username             string `json:"Username"`
	Password             string `json:"Password"`
}

type NetworkInputType struct {
	VlanID           int    `json:"VlanID"`
	Type             string `json:"Type"`
	Name             string `json:"Name"`
	Subnet           string `json:"Subnet"`
	Shutdown         bool   `json:"Shutdown"`
	SubnetAssignment []struct {
		Name         string               `json:"Name"`
		Netmask      int                  `json:"Netmask"`
		IPSize       int                  `json:"IPSize"`
		IPAssignment []NetworkInputIPItem `json:"IPAssignment"`
	} `json:"SubnetAssignment"`
}

type NetworkInputIPItem struct {
	Name     string `json:"Name"`
	Position int    `json:"Position"`
}

type NetworkOutputType struct {
	VlanID       int                   `json:"VlanID"`
	Type         string                `json:"Type"`
	Name         string                `json:"Name"`
	Subnet       string                `json:"Subnet"`
	Shutdown     bool                  `json:"Shutdown"`
	IPAssignment []NetworkOutputIPItem `json:"IPAssignment"`
}

type NetworkOutputIPItem struct {
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
		ID          int      `json:"ID"`
		Port        string   `json:"Port"`
		Speed       int      `json:"Speed"`
		Description string   `json:"Description"`
		Mtu         int      `json:"MTU"`
		PortType    string   `json:"PortType"`
		Shutdown    bool     `json:"Shutdown,omitempty"`
		IPAddress   string   `json:"IPAddress"`
		UntagVlan   string   `json:"UntagVlan"`
		TagVlan     []string `json:"TagVlan"`
	} `json:"Port"`
	OutOfBandPort []struct {
		ID             int    `json:"ID"`
		Name           string `json:"Name"`
		Description    string `json:"Description"`
		Mtu            int    `json:"MTU"`
		Type           string `json:"Type"`
		Network        string `json:"Network"`
		MgmtAssignment string `json:"MgmtAssignment"`
		Shutdown       bool   `json:"Shutdown"`
		NextHop        string `json:"NextHop"`
		Dhcp           bool   `json:"DHCP"`
	} `json:"OutOfBandPort"`
	Vlan []struct {
		Name              string `json:"Name"`
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
	PortChannel []struct {
		Name     string        `json:"Name"`
		ID       int           `json:"ID"`
		PortType string        `json:"PortType"`
		Settings []interface{} `json:"Settings"`
	} `json:"PortChannel"`
}

type OutputType struct {
	Device  DeviceType           `json:"Device"`
	Port    []PortType           `json:"Port"`
	Vlan    []VlanType           `json:"Vlan"`
	Routing *RoutingType         `json:"Routing"`
	Network *[]NetworkOutputType `json:"Network"`
}

type PortType struct {
	Port        string `json:"Port"`
	PortName    string `json:"PortName"`
	PortType    string `json:"PortType"`
	Description string `json:"Description"`
	Mtu         int    `json:"MTU"`
	Shutdown    bool   `json:"Shutdown"`
	IPAddress   string `json:"IPAddress"`
	UntagVlan   int    `json:"UntagVlan"`
	TagVlan     []int  `json:"TagVlan"`
}

type VlanType struct {
	VlanName  string `json:"VlanName"`
	VlanID    int    `json:"VlanID"`
	Type      string `json:"Type"`
	IPAddress string `json:"IPAddress"`
	Mtu       int    `json:"MTU"`
	Shutdown  bool   `json:"Shutdown"`
}

type RoutingType struct {
	Router struct {
		Bgp struct {
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
				RouteMap []struct {
					Name      string `json:"Name"`
					Direction string `json:"Direction"`
				} `json:"RouteMap"`
				UpdateSource string `json:"UpdateSource"`
				Shutdown     bool   `json:"Shutdown"`
			} `json:"IPv4Neighbor"`
		} `json:"BGP"`
		Static []struct {
			Name              string `json:"Name"`
			NetworkName       string `json:"NetworkName"`
			NetworkAssignment string `json:"NetworkAssignment"`
		} `json:"Static"`
	} `json:"Router"`
	PrefixList []struct {
		Index     int    `json:"Index"`
		Name      string `json:"Name"`
		Permit    bool   `json:"Permit"`
		Network   string `json:"Network"`
		Operation string `json:"Operation"`
		Prefix    int    `json:"Prefix"`
	} `json:"PrefixList"`
	RouteMap []struct {
		Index     int    `json:"Index"`
		Name      string `json:"Name"`
		Permit    bool   `json:"Permit"`
		Network   string `json:"Network"`
		Operation string `json:"Operation"`
		Prefix    int    `json:"Prefix"`
	} `json:"RouteMap"`
}
