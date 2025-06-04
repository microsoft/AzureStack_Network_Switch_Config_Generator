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

# TO Do List

## Lab Input

### Switch

- Remove generateConfig/BMC flag, use provide data as true. 1 Tor? No BMC? Use Class to control.
- Add Mux element
- Leave hostname since the cloud section is not the dependcy. Template variable and render based on logic.
- Border connect routing will be global defined and use as string. Like "bgp","static"?
- Only Border for now? If >1 switches, use BGP for peer? iBGP port-channel50 for interface.
- Firmware? no vlan2 for unused port. Only generate config for ports with configuration?

### Supernets

- Name and VLAN ID, which can be the unique key? Network Name or ID? I prefer name which is readable but need to all lower case and unique meaning. VLAN ID can be customized input based on user define.
- GRE? Network input?
- Logic Tenant?
- Gateway? Assignment? Separate segment template?

### Others

- Multi-Cloud? Customer provide or provide incremental logic? Leave an extenable design.
- NTP information in cloud? Hostname will read from it. Define once in template.
- Version, Description will be option for json input track only.
- Config can be modular generated. Like VLAN, BGP, Static.

## To DO

- All key will be capitalize to unify to avoid typo.
- Modular all the configuration. Parse while coding, block by block.
  - Hostname
  - Fibs
  - Vlan

### Design

- Create Global Device Map, which flexiable to adapt to different pattern.
- Clearly separate functions and template by feature, so easier to generate certain config.
