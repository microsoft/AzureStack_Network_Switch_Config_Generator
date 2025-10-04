# CSV Input Format Examples

This document provides complete examples of CSV input files for the PortMap tool.

## Simple 2-Switch Network Example

### sample-devices.csv

```csv
DeviceName,DeviceMake,DeviceModel,Location,PortRange,MediaType,Speed,Description
TOR-1,Cisco,93180YC-FX3,Rack A Position 42,1-48,SFP28,25G,25G interfaces supporting SFP28 transceivers
TOR-1,Cisco,93180YC-FX3,Rack A Position 42,49-52,QSFP28,100G,100G uplink interfaces supporting QSFP28 transceivers
TOR-2,Dell,S5248F-ON,Rack B Position 42,1-48,SFP28,25G,25G server interfaces
TOR-2,Dell,S5248F-ON,Rack B Position 42,49-54,QSFP28,100G,100G uplink and inter-switch interfaces
```

### sample-connections.csv

```csv
SourceDevice,SourcePorts,SourceMedia,DestinationDevice,DestinationPorts,DestinationMedia,ConnectionType,Notes
TOR-1,1-16,SFP28,Server-Nodes,0,SFP28,Server,Compute nodes 1-16
TOR-1,17-24,SFP28,Storage-Nodes,0,SFP28,Storage,Storage nodes 1-8
TOR-1,25,SFP28,LoadBalancer-1,0,SFP28,Network,Primary load balancer
TOR-1,49-50,QSFP28,SPINE-1,1-2,QSFP28,Uplink,Primary uplinks to spine
TOR-2,1-20,SFP28,Server-Nodes,1,SFP28,Server,Secondary NICs for compute nodes 1-20
TOR-2,49-50,QSFP28,SPINE-1,3-4,QSFP28,Uplink,Primary uplinks to spine
TOR-1,51,QSFP28,TOR-2,51,QSFP28,Inter-Switch,MLAG peer link
```

## Breakout Cable Example

For configurations with QSFP breakout cables (e.g., 1x100G port split into 4x25G ports):

### breakout-devices.csv

```csv
DeviceName,DeviceMake,DeviceModel,Location,PortRange,MediaType,Speed,Description
LEAF-1,Cisco,93180YC-FX3,Rack 1 U42,1-24,SFP28,25G,Standard 25G ports
LEAF-1,Cisco,93180YC-FX3,Rack 1 U42,25.1-25.4,QSFP_4x25G,25G,QSFP breakout cable - 4x25G interfaces
LEAF-1,Cisco,93180YC-FX3,Rack 1 U42,26.1-26.4,QSFP_4x25G,25G,QSFP breakout cable - 4x25G interfaces
LEAF-1,Cisco,93180YC-FX3,Rack 1 U42,33-34,QSFP28,100G,Uplink ports
```

### breakout-connections.csv

```csv
SourceDevice,SourcePorts,SourceMedia,DestinationDevice,DestinationPorts,DestinationMedia,ConnectionType,Notes
LEAF-1,1-4,SFP28,HOST-01,1-4,SFP28,Host,Standard connections
LEAF-1,25.1-25.4,QSFP_4x25G,HOST-02,1-4,SFP28,Host,Breakout cable to host
LEAF-1,26.1,QSFP_4x25G,HOST-03,1,SFP28,Host,Single breakout interface
LEAF-1,33,QSFP28,SPINE-01,1,QSFP28,Uplink,Uplink to spine
```

## Usage

Generate documentation from these CSV files:

```powershell
# Generate Markdown documentation
.\PortMap.ps1 -InputFile "sample-devices.csv" -OutputFormat Markdown

# Generate per-device CSV reports
.\PortMap.ps1 -InputFile "sample-connections.csv" -OutputFormat CSV

# Generate JSON for automation
.\PortMap.ps1 -InputFile "breakout-devices.csv" -OutputFormat JSON

# Include unused ports
.\PortMap.ps1 -InputFile "sample-devices.csv" -OutputFormat Markdown -ShowUnused
```

## CSV Format Guidelines

### Devices CSV

**Required Fields:**
- `DeviceName` - Must be unique and match exactly in connections CSV
- `DeviceMake` - Manufacturer name
- `DeviceModel` - Model number/name
- `PortRange` - Port range specification (see formats below)
- `MediaType` - Media type identifier
- `Speed` - Port speed (e.g., 10G, 25G, 40G, 100G)

**Optional Fields:**
- `Location` - Physical location description
- `Description` - Additional notes about the port range

**Port Range Formats:**
- Standard range: `1-48` (ports 1 through 48)
- Single port: `25` (just port 25)
- Breakout range: `25.1-25.4` (breakout interfaces 25.1, 25.2, 25.3, 25.4)
- Single breakout: `26.1` (single breakout interface)

### Connections CSV

**Required Fields:**
- `SourceDevice` - Must match a DeviceName from devices CSV
- `SourcePorts` - Port range on source device
- `SourceMedia` - Media type on source
- `DestinationDevice` - Destination device/host name
- `DestinationMedia` - Media type on destination
- `ConnectionType` - Type of connection (Server, Uplink, Storage, Network, Inter-Switch, etc.)

**Optional Fields:**
- `DestinationPorts` - Port range on destination (use "0" if not applicable)
- `Notes` - Additional connection notes

## Tips

1. **Editing**: Use Excel, Google Sheets, or any CSV editor
2. **Device Names**: Keep them consistent between files (case-sensitive)
3. **Port Ranges**: Use same format as you would in JSON
4. **Comments**: Use the Description/Notes fields for documentation
5. **Validation**: Run with `-Validate` flag to check before generating output

## Converting from JSON to CSV

If you have existing JSON files, you can manually create CSV files by extracting:
- Each device's portRanges → rows in devices.csv
- Each connection → rows in connections.csv

The CSV format is flatter but more spreadsheet-friendly than nested JSON.
