# Input JSON

Reference File: [sample_input.json](docs/sample_input.json)

The json format and struct could be refined and updated based on optimization and new requirement, so file issues if anything inaccurate.

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

| Key      | Value (Example)      | Comment                                                                         |
| -------- | -------------------- | ------------------------------------------------------------------------------- |
| Make     | String, "Cisco"      | Switch vendor name, has to be **exact** matched with switchLib folder name      |
| Model    | String, "93180YC-FX" | Switch model name, has to be **exact** matched with switchLib folder name       |
| Type     | String, "TOR1"       | Switch role. The whole section can be add or remove based on role assignment    |
| Hostname | String, "Cisco-TOR1" | Switch hostname. Hostname can be customized but has to be unique                |
| ASN      | Int, 65000           | BGP AS number, can be null/0 if no BGP                                          |
| Firmware | String, "9.3(9)"     | Switch firmware version, has to be **exact** matched with switchLib folder name |

### InputData.SwitchUplink

This value defines the routing between Border and TOR switches. Here are option values:

| Value  | Comment                                       |
| ------ | --------------------------------------------- |
| BGP    | Use BGP between Border and TOR                |
| Static | Use Static routing BGP between Border and TOR |

### InputData.HostConnectivity

This value defines the connection between TOR and MUX. Here are option values:

| Value  | Comment                                    |
| ------ | ------------------------------------------ |
| BGP    | Use BGP between TOR and MUX                |
| Static | Use Static routing BGP between TOR and MUX |
| L2     | Use pure L2 connection bewtween TOR to MUX |

### InputData.Supernets

This section defines switch Vlan and IP address:

| Key         | Value (Example)          | Comment                                            |
| ----------- | ------------------------ | -------------------------------------------------- |
| GroupID     | String, "TENANT"         | Section group id , easy for multi-section grouping |
| Name        | String, "Tenant201"      | Customized name, default is GroupID+VLANID         |
| VLANID      | Int, 201                 | Vlan ID. 0 if not a Vlan section                   |
| Description | String, "Vlan Tenant201" | Interface description.                             |

#### InputData.Supernets.IPv4

This sub-section defines switch Vlan and IP detail assignment, will only introduce the ones being used in the tool:

| Key        | Value (Example) | Comment                      |
| ---------- | --------------- | ---------------------------- |
| SwitchSVI  | Bool, true      | Create Vlan interface or not |
| Cidr       | Int, 24         | Subnet mask of the IPv4      |
| Assignment | List            | List of detail IP assignment |

### InputData.Setting

This section defines switch global setting variables:

| Key          | Value (Example) | Comment                 |
| ------------ | --------------- | ----------------------- |
| TimeServer   | List            | List of TimeServer IP   |
| SyslogServer | List            | List of SyslogServer IP |
| DNSForwarder | List            | List of DNSForwarder IP |
