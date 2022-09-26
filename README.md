# Switch Configuration Generator

- [Switch Configuration Generator](#switch-configuration-generator)
  - [Project Design](#project-design)
    - [Overall Design](#overall-design)
      - [Logic Diagram](#logic-diagram)
      - [Resource Hierachy](#resource-hierachy)
    - [User Input Design](#user-input-design)
      - [Code Structure](#code-structure)
    - [Switch Framework and Template](#switch-framework-and-template)
    - [Switch JSON Object](#switch-json-object)
      - [Output Object Model](#output-object-model)
      - [IP Caculator Logic of Network Assignment](#ip-caculator-logic-of-network-assignment)
        - [Logic Diagram](#logic-diagram-1)
          - [Segment and Position Logic](#segment-and-position-logic)
    - [Switch Configuration Files](#switch-configuration-files)
      - [Logic Diagram](#logic-diagram-2)
      - [Template Structure](#template-structure)
  - [Concerns and Thoughts](#concerns-and-thoughts)
  - [MileStone Plan](#milestone-plan)
    - [~~Phase1 - POC Phase - 06/24~~](#phase1---poc-phase---0624)
    - [~~Phase2 Testing on Single Switch Validation Virtual (GNS3 Lab) - 08/19~~](#phase2-testing-on-single-switch-validation-virtual-gns3-lab---0819)
    - [~~Phase3 Testing on Multi-Switch Validation Virtual (GNS3 Lab) - 09/23~~](#phase3-testing-on-multi-switch-validation-virtual-gns3-lab---0923)
    - [Phase4 Physical Device Testing - 10/30](#phase4-physical-device-testing---1030)

## Project Design

### Overall Design

#### Logic Diagram

```mermaid
flowchart TD
    A[Switch Framework + Go Template]
    B[User Input Template]
    C(Generator Tool)
    D(Switch Output Object)
    E[Switch Object JSON Files]
    F[Switch Configuration Files]
    B --> C
    A --> C
    C --> D
    D -.-> |For Debug| E
    D --> |For Deploy| F
```

#### Resource Hierachy

Example:

```
.
├── input
│   ├── cisco
│   │   └── 93180yc-fx
│   │       └── 9.3
│   │           ├── framework
│   │           │   ├── interface.json
│   │           │   ├── ntp.json
│   │           │   └── routing.json
│   │           └── template
│   │               ├── allConfig.go.tmpl
│   │               ├── bgp.go.tmpl
│   │               ├── header.go.tmpl
│   │               ├── inBandPort.go.tmpl
│   │               └── vlan.go.tmpl
│   ├── dell
│   │   └── s3048-on
│   │       ├── os10
│   │       │   ├── framework
│   │       │   └── template
│   │       └── os9
│   └── input.json
├── output
│   ├── S31R28-TOR1.config
│   ├── S31R28-TOR1.json
│   ├── S31R28-TOR2.config
│   ├── S31R28-TOR2.json
│   └── tor1.conf
└── src
    ├── go.mod
    ├── go.sum
    ├── inbandport.go
    ├── ipcaculator.go
    ├── ipcaculator_test.go
    ├── main.go
    ├── output.go
    ├── routing.go
    ├── strcture.go
    ├── SwitchConfigGenerator
    ├── testcases
    │   ├── input1.json
    │   ├── input2.json
    │   ├── result1.json
    │   └── result2.json
    └── vlan.go
```

### User Input Design

#### Code Structure

- Device : Switches used in the Deployment
- Network : Required network variables user defined
- External : Optional network variables user defined

```Go
type InputType struct {
	External []struct {
		Type string   `json:"Type"`
		IP   []string `json:"IP"`
	} `json:"External"`
	Device []struct {
		Make                 string `json:"Make"`
		Type                 string `json:"Type"`
		Asn                  int    `json:"ASN"`
		Hostname             string `json:"Hostname"`
		Model                string `json:"Model"`
		Firmware             string `json:"Firmware"`
		GenerateDeviceConfig bool   `json:"GenerateDeviceConfig"`
	} `json:"Device"`
	Network []struct {
		VlanID           int           `json:"VlanID"`
		Type             string        `json:"Type"`
		Name             string        `json:"Name"`
		Subnet           string        `json:"Subnet"`
		Shutdown         bool          `json:"Shutdown"`
		SubnetAssignment []interface{} `json:"SubnetAssignment"`
	} `json:"Network"`
}
```

### Switch Framework and Template

This part is the core of the project. Each switch need to have paired `framework` and `template` files to be able generate configuration accordingly.

### Switch JSON Object

#### Output Object Model

```Go
type OutputType struct {
	Device struct {
		Make                 string `json:"Make"`
		Type                 string `json:"Type"`
		Asn                  int    `json:"ASN"`
		Hostname             string `json:"Hostname"`
		Model                string `json:"Model"`
		Firmware             string `json:"Firmware"`
		GenerateDeviceConfig bool   `json:"GenerateDeviceConfig"`
	} `json:"Device"`
	Port []struct {
		Port        string `json:"Port"`
		PortName    string `json:"PortName"`
		PortType    string `json:"PortType"`
		Description string `json:"Description"`
		Mtu         int    `json:"MTU"`
		Shutdown    bool   `json:"Shutdown"`
		IPAddress   string `json:"IPAddress"`
		UntagVlan   int    `json:"UntagVlan"`
		TagVlan     []int  `json:"TagVlan"`
	} `json:"Port"`
	Vlan []struct {
		VlanName  string `json:"VlanName"`
		VlanID    int    `json:"VlanID"`
		Type      string `json:"Type"`
		IPAddress string `json:"IPAddress"`
		Mtu       int    `json:"MTU"`
		Shutdown  bool   `json:"Shutdown"`
	} `json:"Vlan"`
	Routing struct {
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
			Static interface{} `json:"Static"`
		} `json:"Router"`
		PrefixList []struct {
			Index     int    `json:"Index"`
			Name      string `json:"Name"`
			Permit    bool   `json:"Permit"`
			Network   string `json:"Network"`
			Operation string `json:"Operation"`
			Prefix    int    `json:"Prefix"`
		} `json:"PrefixList"`
		RouteMap interface{} `json:"RouteMap"`
	} `json:"Routing"`
	Network []struct {
		VlanID       int           `json:"VlanID"`
		Type         string        `json:"Type"`
		Name         string        `json:"Name"`
		Subnet       string        `json:"Subnet"`
		Shutdown     bool          `json:"Shutdown"`
		IPAssignment []interface{} `json:"IPAssignment"`
	} `json:"Network"`
}
```

#### IP Caculator Logic of Network Assignment

##### Logic Diagram

```mermaid
flowchart TD
    A[InputObj.Network.Subnet]
    B[InputObj.Network.SubnetAssignment]
    C(All IPAddress in the Subnet)
    D(Position Index in the IP List)
    E[OutputObj.Network.IPAssignments]
    A -.-> |Create IP List| C
    B --> |Segment and Position| D
    C <-.-> |Index the IP List| D
    D --> |Generate| E
```

###### Segment and Position Logic

The IP Assignment will be sorted by assigned netmask in **descend order** and then assigned by assgined position.
Example:

- Given Input:

```json
{
  "VlanID": null,
  "Type": "IP",
  "Name": "SwitchMgmt",
  "Subnet": "192.168.1.0/28",
  "Shutdown": false,
  "SubnetAssignment": [
    {
      "Name": "P2P_TOR1_To_Border1",
      "Netmask": 31,
      "IPSize": 2,
      "IPAssignment": [
        {
          "Name": "TOR1_Loopback0",
          "Netmask": 32,
          "IPSize": 1,
          "IPAssignment": [
            {
              "Name": "TOR1_Loopback0",
              "Position": 0
            }
          ]
        },
        {
          "Name": "P2P_TOR1_To_Border1",
          "Netmask": 31,
          "IPSize": 2,
          "IPAssignment": [
            {
              "Name": "TOR1",
              "Position": 0
            },
            {
              "Name": "Border1",
              "Position": 1
            }
          ]
        }
      ]
    }
  ]
}
```

- Get Output:

```json
{
  "VlanID": 0,
  "Type": "IP",
  "Name": "SwitchMgmt",
  "Subnet": "192.168.1.0/28",
  "Shutdown": false,
  "IPAssignment": [
    {
      "Name": "P2P_TOR1_To_Border1/TOR1",
      "IPAddress": "192.168.1.0/31"
    },
    {
      "Name": "P2P_TOR1_To_Border1/Border1",
      "IPAddress": "192.168.1.1/31"
    },
    {
      "Name": "TOR1_Loopback0/TOR1_Loopback0",
      "IPAddress": "192.168.1.2/32"
    }
  ]
}
```

### Switch Configuration Files

Switch configuration is generated by using Go native package: [text/template](https://pkg.go.dev/text/template)

#### Logic Diagram

```mermaid
flowchart
    A[Switch Output Object]
    B(allConfig.go.tmpl)
    C[Final Switch Configuration]
    D(header.go.tmpl)
    E(vlan.go.tmpl)
    F(bgp.go.tmpl)
    G(xxx.go.tmpl)

    A <-.-> |Parse| D
    A <-.-> |Parse| E
    A <-.-> |Parse| F
    A <-.-> |Parse| G
    D --> |Merge| B
    E --> |Merge| B
    F --> |Merge| B
    G --> |Merge| B
    B --> C
```

#### Template Structure

| Config     | Template           | Source                       |
| ---------- | ------------------ | ---------------------------- |
| All Config | allConfig.go.tmpl  | All templates below          |
| Header     | header.go.tmpl     | OutputObj.Device             |
| VLAN       | vlan.go.tmpl       | OutputObj.Vlan               |
| InBandPort | inBandPort.go.tmpl | OutputObj.Port               |
| BGP        | bgp.go.tmpl        | OutputObj.Routing.Router.Bgp |

## Concerns and Thoughts

- Current functions heavily depend on given name.

  - Example: In input.json file, all vlan name has to be matched exact with related framework file.
  - Thoughts: We can define/hardcode what we want right now. However, use pre-defined portal/webpage? Bring more unecessary maintain or could reduce confusion?

- Due to different deploy pattern, the tool is very fragile and complex to maintain.

  - Thoughts: Do Not try to cover all the use cases, focus on the major(70%) customer deploy scenario, and make it reliable with well documentation. Then leave rest 30% manual modify if need at this moment.
  - So far only focus on Cisco NXOS, and the tool may break if generate Dell or others. The most fragile part are `Supernet` and `Routing`.

- Phase 1 Summary:
  - Redesign `routing.json` with prefix and routemap allocation, `input.json` with setting, `qos.json` [not in use], do we need dynamic or just static?
  - Tested in virutal lab, result looks good except the Qos feature not fully support in vImage.
```
TOR1# copy bootflash:///TOR1.config running-config
As per NIST requirements, the minimum RSA Key Size has to be 2048 in FIPS Mode
Generate RSA key with 2048 bits

Password prompt username is enabled.
After providing the required options in the username command, press enter.
User will be prompted for the username password and password will be hidden.
Note: Choosing password key in the same line while configuring user account, password will not be hidden.
XML interface to system may become unavailable since ssh is disabled
Note: To enable relay on any interface, please disable DHCP (v4/v6) on interfaces that have address assigned via DHCP (dynamic IP addressing).

2022 Aug  8 23:57:13 S31R28-TOR1 %$ VDC-1 %$ %SECURITYD-2-FEATURE_NXAPI_ENABLE: Feature nxapi is being enabled on HTTPS.
ERROR: Configuration failed, Interface does not exists or Ip address not configured on the interface.
Unable to perform the action due to incompatibility:  Module 1 returned status "ACLQOS_ERROR: Invalid Class id : Max class reached"
S31R28-TOR1(config)# system qos
S31R28-TOR1(config-sys-qos)#   service-policy type queuing output QOS_EGRESS_PORT
Unable to perform the action due to incompatibility:  Module 1 returned status "ACLQOS_ERROR: Invalid Class id : Max class reached"

S31R28-TOR1(config-sys-qos)#   service-policy type network-qos QOS_NETWORK^C
S31R28-TOR1(config-sys-qos)#



Edge port type (portfast) should only be enabled on ports connected to a single
 host. Connecting hubs, concentrators, switches, bridges, etc...  to this
 interface when edge port type (portfast) is enabled, can cause temporary bridging loops.
 Use with CAUTION


Copy complete, now saving to disk (please wait)...
Copy complete.
S31R28-TOR1# show version
Cisco Nexus Operating System (NX-OS) Software
TAC support: http://www.cisco.com/tac
Documents: http://www.cisco.com/en/US/products/ps9372/tsd_products_support_serie
s_home.html
Copyright (c) 2002-2021, Cisco Systems, Inc. All rights reserved.
The copyrights to certain works contained herein are owned by
other third parties and are used and distributed under license.
Some parts of this software are covered under the GNU Public
License. A copy of the license is available at
http://www.gnu.org/licenses/gpl.html.

Nexus 9000v is a demo version of the Nexus Operating System

Software
  BIOS: version
  NXOS: version 10.1(1)
  BIOS compile time:
  NXOS image file is: bootflash:///nxos64.10.1.1.bin
  NXOS compile time:  2/14/2021 22:00:00 [02/15/2021 07:39:04]

Hardware
  cisco Nexus9000 C9500v Chassis ("Supervisor Module")
   with 7832932 kB of memory.
  Processor Board ID 99IWXCECUI7
  Device name: S31R28-TOR1
  bootflash:    4287040 kB

Kernel uptime is 0 day(s), 0 hour(s), 22 minute(s), 44 second(s)

Last reset
  Reason: Unknown
  System version:
  Service:

plugin
  Core Plugin, Ethernet Plugin

Active Package(s):
```

## MileStone Plan

### ~~Phase1 - POC Phase - 06/24~~

  - Existing documentation
  - Unit Test
  - Optimize/Comment current code

### ~~Phase2 Testing on Single Switch Validation Virtual (GNS3 Lab) - 08/19~~

  - Function To Do List:
    - ~~`OutOfBandPort`~~ - port.go.tmpl-port + interface.json
    - ~~`Credential`~~ - stig.go.tmpl-stig_user + outputObj
    - ~~`Stig_Setting`~~ - stig.go.tmpl-stig_user
    - ~~`PrefixList`~~ - static.go.tmpl-static/bgp.go.tmpl-bgp + routing.json
    - ~~`NTP`~~ - settings.go.tmpl-ntp + outputObj
    - ~~`STP`~~ - stp.go.tmpl-inject_stp
    - ~~`VTY`~~ - default.go.tmpl-default_console_vty
    - ~~`Static`~~ - static.go.tmpl-static + routing.json
    - ~~`BGP`~~ - static/bgp.go.tmpl-bgp + routing.json
    - ~~`QOS`~~ - qos.go.tmpl-bgp + ?qos.json?
    - ~~`Syslog`~~ - settings.go.tmpl-inject_syslog + outputObj
  - ~~Complete `Unused` Port in `interface.json` framework~~

### ~~Phase3 Testing on Multi-Switch Validation Virtual (GNS3 Lab) - 09/23~~
  - Function To Do List: 
    - ~~`Port-Channel`~~
    - ~~`MLAG`~~
  - ~~Define BMC Framework~~
  - ~~BMC Cisco Switch Integration~~

### Phase4 Physical Device Testing - 10/30

  - Single Switch
  - Multi-Switch
  - CI/CD for Lab
  - Documentation
  - Lab team training