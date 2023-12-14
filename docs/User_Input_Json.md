# User Input Definition JSON

Reference File: [sample_input_files](/src/test/testInput)

The json format and structure could be refined and updated based on optimization and new requirement, so file issues if anything inaccurate.

## Structure

### Global

| Key         | Value (Example)            | Comment          |
| ----------- | -------------------------- | ---------------- |
| Version     | String, ("1.0.0")          | Reserved         |
| Description | String, ("Input Template") | Reserved         |
| InputData   | Object                     | Used in the tool |

### InputData.MainEnvData

| Key         | Value (Example) | Comment                                          |
| ----------- | --------------- | ------------------------------------------------ |
| MainEnvData | Object          | NOT used in this tool, used for sever deployment |

### InputData.Switches

This section is used for switch framework *switchLib*selection and the directory path is the key, which has to be **extact** matched with input values (case insensitive).

| Key      | Value (Example)      | Comment                                                                                                                                |
| -------- | -------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| Make     | String, "Cisco"      | Switch vendor name, has to be **exact** matched with switchLib folder name, case sensitive                                             |
| Model    | String, "93180YC-FX" | Switch model name, has to be **exact** matched with switchLib folder name, case sensitive                                              |
| Type     | String, "TOR1"       | Switch role. Type includes ["BORDER1","BORDER2","TOR1","TOR2","BMC","MUX"], but the whole unit can be removed based on deployment need |
| Hostname | String, "Cisco-TOR1" | Switch hostname. Hostname can be customized but has to be unique                                                                       |
| ASN      | Int, 65000           | BGP AS number, can be null/0 if no BGP used                                                                                            |
| Firmware | String, "9.3(9)"     | Switch firmware version, has to be **exact** matched with switchLib folder name                                                        |

### InputData.DeploymentPattern

This value defines the deployment pattern, and value will be used for generating configuration. Here are option values:

| Value          | Comment                              |
| -------------- | ------------------------------------ |
| Hyperconverged | Storage switched and fully converged |
| Switched       | Storage switched but non-converged   |
| Switchless     | Storage switchless                   |

Note: Check [this link](https://learn.microsoft.com/en-us/azure-stack/hci/plan/choose-network-pattern) for more deployment detail.

### InputData.SwitchUplink

This value defines the routing between Border and TOR switches, and value will be used for generating configuration. Here are option values:

| Value  | Comment                                       |
| ------ | --------------------------------------------- |
| BGP    | Use BGP between Border and TOR                |
| Static | Use Static routing BGP between Border and TOR |

### InputData.HostConnectivity

This value defines the connection between TOR and MUX. Here are option values:

| Value  | Comment                                                      |
| ------ | ------------------------------------------------------------ |
| BGP    | Use BGP between TOR and MUX                                  |
| Static | [To Be Developed] Use Static routing BGP between TOR and MUX |
| L2     | [To Be Developed] Use pure L2 connection bewtween TOR to MUX |

### InputData.Supernets

This section defines switch Vlan and IP address:

| Key       | Value (Example)  | Comment                                              |
| --------- | ---------------- | ---------------------------------------------------- |
| GroupName | String, "TENANT" | Section group name , easy for multi-section grouping |

| Description | String, "Vlan Tenant201" | Interface description. |

#### InputData.Supernets.IPv4

This sub-section defines switch Vlan and IP detail assignment, will only introduce the ones being used in the tool:

| Key         | Value (Example)          | Comment                                                |
| ----------- | ------------------------ | ------------------------------------------------------ |
| Name        | String, "Tenant201"      | Customized name, default is GroupID+VLANID             |
| NetworkType | String, "TENANT"         | Secondary group name , easy for multi-section grouping |
| VLANID      | Int, 201                 | Vlan ID. 0 if not a Vlan section                       |
| Cidr        | Int, 30                  | Network Classless Inter-Domain Routing                 |
| Subnet      | String, "10.69.176.0/24" | Network subnet                                         |
| Netmask     | String, "255.255.255.0"  | Network mask                                           |
| Gateway     | String, "10.69.176.1"    | Network gateway                                        |
| SwitchSVI   | Bool, true               | Create Vlan interface or not                           |
| Cidr        | Int, 24                  | Subnet mask of the IPv4                                |
| Assignment  | List                     | List of detail IP assignment                           |

#### InputData.Supernets.IPv4.Assignment

This sub-section defines the detail IP assignment

| Key  | Value (Example)       | Comment                                |
| ---- | --------------------- | -------------------------------------- |
| Name | String, "TOR1"        | Specific network assignment name       |
| IP   | String, "10.69.176.2" | Specific network assignment ip address |

### InputData.Setting

This section defines switch global setting variables:

| Key          | Value (Example) | Comment                 |
| ------------ | --------------- | ----------------------- |
| TimeServer   | List            | List of TimeServer IP   |
| SyslogServer | List            | List of SyslogServer IP |
| DNSForwarder | List            | List of DNSForwarder IP |

### InputData.WANSIM

This section defines WANSIM VM network variables to generate WANSIM `netplan` and `frr` configuration.

| Key             | Value (Example)       | Comment                                              |
| --------------- | --------------------- | ---------------------------------------------------- |
| Hostname        | String,"rack1-wansim" | Hostname of WANSIM (Not being used, can replace DNS) |
| Loopback        | Object                | Loopback as Tunnel Source IP and BGP Router ID       |
| GRE1            | Object                | Tunnel Variables with TOR1                           |
| GRE2            | Object                | Tunnel Variables with TOR2                           |
| BGP             | Object                | BGP Peer with TORs via GRE Tunnels                   |
| RerouteNetworks | Object                | List of GroupName of Supernets                       |

#### InputData.WANSIM.Loopback

This section defines WANSIM VM Loopback setting variables, which will be used for two GRE Tunnel Source IP, so has to be advertised and unique in the network:

| Key       | Value (Example)           | Comment                         |
| --------- | ------------------------- | ------------------------------- |
| IP        | string, "10.10.32.129"    | VM Assigned Loopback IP Address |
| IPNetwork | string, "10.10.32.129/32" | VM Assigned Loopback IPNetwork  |
| Subnet    | string,"10.10.32.128/25"  | Loopback Subnet                 |

#### InputData.WANSIM.GRE1

This section defines the GRE IP Information between WANSIM VM and TOR1, these IPs are all private for the Tunnel so can be reused:

| Key       | Value  | Example      | Comment                                     |
| --------- | ------ | ------------ | ------------------------------------------- |
| Name      | string | "TOR1"       | Reserved, not being used.                   |
| LocalIP   | string | "2.1.1.0"    | WANSIM VM GRE1 Tunnel Local IP              |
| RemoteIP  | string | "2.1.1.1"    | WANSIM VM GRE1 Tunnel Remote IP (TOR1 side) |
| IPNetwork | string | "2.1.1.0/31" | GRE1 Tunnel Subnet                          |

#### InputData.WANSIM.GRE2

This section defines the GRE IP Information between WANSIM VM and TOR2, these IPs are all private for the Tunnel so can be reused:

| Key       | Value  | Example      | Comment                                     |
| --------- | ------ | ------------ | ------------------------------------------- |
| Name      | string | "TOR2"       | Reserved, not being used.                   |
| LocalIP   | string | "2.1.1.2"    | WANSIM VM GRE2 Tunnel Local IP              |
| RemoteIP  | string | "2.1.1.3"    | WANSIM VM GRE2 Tunnel Remote IP (TOR2 side) |
| IPNetwork | string | "2.1.1.2/31" | GRE2 Tunnel Subnet                          |

#### InputData.WANSIM.BGP

This section defines the BGP information to generate BGP in FRR config:

| Key      | Value | Example | Comment                        |
| -------- | ----- | ------- | ------------------------------ |
| LocalASN | int   | 65003   | Reserved, not being used.      |
| IPv4Nbr  | List  |         | WANSIM VM GRE2 Tunnel Local IP |

##### InputData.WANSIM.BGP.IPv4Nbr

Because the TOR switches information already includes in the file, so only need to put the uplink switches if any to peer with WANSIM.

| Key               | Value  | Example      | Comment                                   |
| ----------------- | ------ | ------------ | ----------------------------------------- |
| NeighborAsn       | int    | 65001        | Uplink Switch ASN                         |
| NeighborIPAddress | string | "10.10.36.2" | Uplink Switch IP Address to Peer BGP      |
| Description       | string | "To_Uplink1" | BGP Neighbor Description                  |
| EbgpMultiHop      | int    | 8            | EBGP Multihop Value                       |
| UpdateSource      | string | "eth0"       | Source Interface of WANSIM VM to Peer BGP |

#### InputData.WANSIM.RerouteNetworks

This section defines the network need to be redirected into WANSIM VM:

| Key             | Value | Example                     | Comment                                        |
| --------------- | ----- | --------------------------- | ---------------------------------------------- |
| RerouteNetworks | List  | ["Infrastructure","TENANT"] | List of GroupName defined in Supernets section |

## Q&A

#### Are there any examples to understand the mapping between definition and configuration.

Use this definition JSON as example: [s46r21-definition.json](/src/test/testInput/s46r21-definition.json)

The tool generates two parts of configuration:

- [Azure Stack Switch Configuration](/src/test/testOutput/s46r21-definition/)

  - TOR Switches Configuration.
  - BMC Switch Configuration if has.

- [WANSIM VM Configuration](/src/test/testOutput/s46r21-definition/wansim_config/) which peer with Azure Stack Switch. [Click Here](<(https://github.com/microsoft/AzureStackWANSimulator)>) to understand WANSIM Feature.
  - VM netplan YAML File, which defines all the interfaces with IP.
  - FRR configuration: daemons + frr.conf.
  - Default network profile rules: 1Gbit Download and Upload Bandwidth.
  - Bash script to add WANSIM related log into syslog.
