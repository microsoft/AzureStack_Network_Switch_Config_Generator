Vlan
Interface
Routing
Global Setting

### SwitchLib Hierachy Sample

```
.
├── cisco
│   └── 9.3(9)
│       ├── 93180yc-fx
│       │   └── interface.json
│       ├── 9348gc-fxp
│       │   └── interface.json
│       └── template
│           ├── AllConfig.go.tmpl
│           ├── hostname.go.tmpl
│           ├── port.go.tmpl
│           ├── settings.go.tmpl
│           ├── stig.go.tmpl
│           └── vlan.go.tmpl
└── dellemc
    └── 10.5(3.4)
        ├── n3248te-on
        │   └── interface.json
        ├── s5248-on
        │   └── interface.json
        └── template
            ├── AllConfig.go.tmpl
            ├── hostname.go.tmpl
            ├── port.go.tmpl
            ├── settings.go.tmpl
            ├── stig.go.tmpl
            └── vlan.go.tmpl
```

#### User Input Template

```Go
type InputType struct {
	Version   string                 `json:"Version"`
	Settings  map[string]interface{} `json:"Settings"`
	IsNoBMC   bool                   `json:"IsNoBMC"`
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
```

#### Switch Framework JSON

Example: BGP Framework

```Go
type BGPType struct {
	BGPAsn                 string   `json:"BGPAsn"`
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
```

#### Switch Go Template

Example: BGP Template

```Go
{{define "bgp_prefix"}}
! bgp_prefix_list
{{ range .PrefixList -}}
{{ range .Config -}}
ip prefix-list {{.Name}} {{.Action}} {{.Supernet}}
{{end}}
{{end}}
{{end}}

{{define "bgp_routing"}}
! bgp.go.tmpl-bgp
router bgp {{.BGPAsn}}
  router-id {{.RouterID}}
  bestpath as-path multipath-relax
  log-neighbor-changes
  address-family ipv4 unicast
    maximum-paths 9
    {{- range .IPv4Network}}
    network {{.}}
    {{- end -}}
{{- /* Define variable before assign*/ -}}
{{$MaxiPrefix:= .RoutePrefix.MaxiPrefix}}
{{$ErrorAction:= .RoutePrefix.ErrorAction}}
  {{- range .IPv4Neighbor}}
  neighbor {{.NeighborIPAddress}}
    remote-as {{.NeighborAsn}}
    description {{.Description}}
    address-family ipv4 unicast
      maximum-prefix {{$MaxiPrefix}} {{$ErrorAction}}
      {{ range .PrefixList -}}
      prefix-list {{.Name}} {{.Direction}}
      {{end -}}
  {{end -}}
{{end}}
```
