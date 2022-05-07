[_TOC_]

To Do List:

- Separate [framework.json](framework.json) based on the diagram.
![Framework Diagram](frameworkSeparation.png)
- Define reference input json template.
- Define Vendor/Firmware/Template hierarchy.

## Template Paths

The Template files will match the device make, model and firmware.

```text
Template/Cisco/93180YC-FX/9.3/
   - index.json
   - NTP.tpl
   - SNMP.tpl
   - BGP.tpl
   ...
```

The index file will determine what files are called and in what order.  Once all the templates are processed, the index will be used to combine the data into a single file.  The exported file name will be based on the Hostname in the Input details.

### index.json

```JSON
{
    "Index": [
        "Username.tpl",
        "SNMP.tpl",
        "Interfaces.tpl",
        "BGP.tpl",
        "NTP.tpl"
    ]
}
```

## Input Details

The inital set of inputs will remain as simple as possible.  This will consist of a set of devices in the rack.

### Assumptions

1. The inital set of rack configurations will be fixed.  
2. All VLANs will be tagged to all defined host ports.  
3. The port definition will be based ont the framework files.
4. There will be 16 nodes per rack.
5. The network port connectivity will be based on a HUB layout.
6. The routing protocol will use BGP.

### Questions

1. 16 total nodes in a rack.  Should we outline this in the input?
2. Port utilization is not outlined. Ports 1-x for Compute, x-x for Storage
3. If BGP is being used, the AS numbers should be included in the input file. 

```JSON
{
    "Framework": "path/to/directory",
    "Device": [
        {
            "Make": "Cisco",
            "Type": "TOR1",
            "Hostname": "S31R28-TOR1",
            "Model": "93180YC-FX",
            "Firmware": 9.3
        },
        {
            "Make": "Cisco",
            "Type": "TOR2",
            "Hostname": "S31R28-TOR1",
            "Model": "93180YC-FX",
            "Firmware": 9.3
        },
        {
            "Make": "Dell",
            "Type": "BMC",
            "Hostname": "S31R28-BMC",
            "Model": "S3048",
            "Firmware": 9.2
        },
    ],
    "SwitchVlan": [
        {
            "ID": 2,
            "Type": "Unused",
            "Name": "UnusedPort",
            "Subnet": "",
        },
        {
            "ID": 7,
            "Type": "Compute",
            "Name": "HNVPANetwork",
            "Subnet": "10.10.100.0/24"
        },
        {
            "ID": 8,
            "Type": "Compute",
            "Name": "Management",
            "Subnet": "10.10.101.0/24"
        },
        {
            "ID": 25,
            "Type": "OOB",
            "Name": "BMC",
            "Subnet": "192.168.0.0/26"
        },
        {
            "ID": 15,
            "Type": "Tenant",
            "Name": "Tenant1",
            "Subnet": "10.100.0.0/24"
        },
        {
            "ID": 711,
            "Type": "Storage",
            "Name": "Storage1",
            "Subnet": ""
        },
        {
            "ID": 712,
            "Type": "Storage",
            "Name": "Storage2",
            "Subnet": ""
        }
    ]
}
```
