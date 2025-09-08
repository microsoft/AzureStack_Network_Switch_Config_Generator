# PortMap Tool

## Overview

The PortMap tool is a comprehensive network documentation utility designed to generate professional-grade documentation for physical port assignments and cable connections across network devices. It creates detailed documentation artifacts that describe the complete physical cabling configuration of network switches, serving as essential documentation for network design, implementation, and maintenance.

## Purpose

This tool solves critical documentation challenges in network design by:

- **Comprehensive Physical Documentation**: Creating detailed artifacts of switch physical cable configurations with complete metadata
- **Sequential Port Mapping**: Providing logically ordered port documentation with connected and unused ports in numerical sequence
- **Advanced Port Analysis**: Tracking port utilization, media types, connection status, and availability
- **Professional Output Generation**: Producing multiple output formats with intelligent file naming and embedded metadata
- **Integration Ready**: Generating outputs that seamlessly integrate with switch configuration generators and network management tools

## Features

### Core Functionality

- **Physical-to-Physical Mapping**: Maps source device ports to destination device ports with comprehensive media type and connection information
- **Sequential Port Display**: Shows all ports in numerical order (1, 2, 3...) with connected and unused ports properly integrated
- **Advanced Port Utilization Analysis**: Identifies and categorizes used, unused, and available ports across all devices
- **Device Inventory Management**: Maintains detailed device make/model information alongside comprehensive port mappings
- **Multiple Output Formats**: Supports JSON, CSV, and Markdown formats with consistent formatting and metadata

### Advanced Features

- **Intelligent File Naming**: Automatic generation of descriptive output filenames including device information, manufacturers, and parameters
- **Incremental File Management**: Automatic incremental numbering to prevent file overwrites
- **Embedded Metadata**: Complete input parameters, generation details, and device context included in all outputs
- **Sequential Port Ordering**: All ports displayed in logical numerical sequence regardless of connection status
- **Professional Formatting**: Network engineering-standard documentation format with clear status indicators

### Output Formats

#### 1. Markdown Tables

Professional network documentation format with per-device organization:

```markdown
### TOR-1

**Device:** TOR-1 (Cisco 93180YC-FX3)
**Total Ports:** 52
**Connected Ports:** 28
**Unused Ports:** 24

| Port | Media | Status     | Destination Device | Destination Port | Destination Media | Type    | Notes                 |
| ---- | ----- | ---------- | ------------------ | ---------------- | ----------------- | ------- | --------------------- |
| 1    | SFP28 | Active     | Server-Nodes       | 0                | SFP28             | Server  | Compute nodes 1-16    |
| 2    | SFP28 | Active     | Server-Nodes       | 0                | SFP28             | Server  | Compute nodes 1-16    |
| 25   | SFP28 | Active     | LoadBalancer-1     | 0                | SFP28             | Network | Primary load balancer |
| 26   | SFP28 | **Unused** | -                  | -                | -                 | -       | Available             |
| 27   | SFP28 | **Unused** | -                  | -                | -                 | -       | Available             |
```

#### 2. CSV Format - Individual Files Per Device

Creates individual CSV files for each device with comprehensive metadata and structured data:

**File Structure:**

- **One CSV file per device** with device-specific naming
- **Metadata section** with device details, port counts, and generation information
- **Sequential port data** with connected and unused ports in numerical order
- **Complete connection details** for each port including media types and destinations

**Example Files Generated:**

- `sample-network-config-Cisco-93180YC-FX3-1dev-portmap.csv` (TOR-1 device)
- `sample-network-config-Dell-S5248F-ON-1dev-portmap.csv` (TOR-2 device)  
- `sample-network-config-Cisco-N9K-C93360YC-FX2-1dev-portmap.csv` (SPINE-1 device)

**CSV Structure:**

```csv
"Section","DeviceName","Make","Model","Location","TotalPorts","UsedPorts","UnusedPorts"...
"Metadata","TOR-1","Cisco","93180YC-FX3","Rack A, Position 42","52","28","24"...
"---","---","---","---","---","---","---","---"...
"Port","TOR-1","Cisco","93180YC-FX3","Rack A, Position 42","","","","","","","","1","SFP28","Active","Server-Nodes","0","SFP28","Server","Compute nodes 1-16"
"Port","TOR-1","Cisco","93180YC-FX3","Rack A, Position 42","","","","","","","","26","SFP28","Unused","","","","","Available"
```

#### 3. JSON Format - Consolidated Output

Machine-readable format with complete device and connection information for integration with configuration generators and other network tools. JSON output creates a single consolidated file with all devices.

## Target Audience

- **Network Engineers**: Primary users for designing and documenting physical network layouts
- **System Administrators**: For understanding and maintaining network infrastructure
- **Documentation Teams**: For creating network design artifacts and compliance documentation

## Input Requirements

### Device Configuration JSON

The tool requires a JSON file describing devices at a high level:

```json
{
  "devices": [
    {
      "deviceName": "TOR-1",
      "deviceMake": "Cisco",
      "deviceModel": "93180YC-FX3",
      "portRanges": [
        {
          "range": "1-48",
          "mediaType": "SFP28",
          "speed": "25G",
          "description": "25G interfaces supporting SFP28"
        },
        {
          "range": "49-52",
          "mediaType": "QSFP28",
          "speed": "100G",
          "description": "100G interfaces supporting QSFP28"
        }
      ]
    }
  ],
  "connections": [
    {
      "sourceDevice": "TOR-1",
      "sourcePorts": "1-9",
      "sourceMedia": "SFP28",
      "destinationDevice": "Node1-Node9",
      "destinationPorts": "0",
      "destinationMedia": "SFP28"
    },
    {
      "sourceDevice": "TOR-1",
      "sourcePorts": "16",
      "sourceMedia": "SFP28",
      "destinationDevice": "LoadBalancer1",
      "destinationPorts": "0",
      "destinationMedia": "SFP28"
    }
  ]
}
```

## Installation & Dependencies

### Prerequisites

- PowerShell 5.1 or PowerShell Core 7.0+ (cross-platform compatible)
- JSON configuration file describing your network devices and connections
- Write permissions to the output directory

### Dependencies

- No external PowerShell modules required
- Uses built-in PowerShell JSON and CSV cmdlets
- Leverages .NET Framework/Core collections for optimal performance

### Key Features & Recent Improvements

#### Sequential Port Integration

When `-ShowUnused` is specified, the tool displays all ports in perfect numerical sequence (1, 2, 3, ..., N) with connected and unused ports seamlessly integrated rather than in separate sections.

#### Intelligent File Naming

Automatically generates descriptive filenames incorporating:

- Input configuration name
- Device manufacturers and models
- Device count and tool identifier
- Automatic incremental numbering

#### Professional Markdown Output

Per-device organization with:

- Device summary information (make, model, port counts)
- Unified port tables showing connected and unused ports
- Clear status indicators and availability information
- Network engineering documentation standards

#### Comprehensive Metadata

All output formats include complete generation context:

- Tool version and generation timestamp
- Input parameters and processing details
- Device inventory and utilization statistics
- User and system information for audit trails

### CSV Output Behavior

The CSV output format has unique behavior compared to other formats:

#### Individual Files Per Device

- **One CSV file per device** rather than a single consolidated file
- **Device-specific filenames** automatically generated for each device
- **Self-contained data** - each CSV includes complete device context and metadata

#### File Naming Pattern

CSV files follow this naming convention:

```text
{input-config-name}-{manufacturer}-{model}-1dev-portmap[-increment].csv
```

#### Examples

For a 3-device configuration, CSV output creates:

- `sample-config-Cisco-93180YC-FX3-1dev-portmap.csv` (TOR-1)
- `sample-config-Dell-S5248F-ON-1dev-portmap.csv` (TOR-2)  
- `sample-config-Cisco-N9K-C93360YC-FX2-1dev-portmap.csv` (SPINE-1)

#### CSV vs Other Formats

| Format   | Output Type | Files Created | File Naming |
|----------|-------------|---------------|-------------|
| Markdown | Consolidated | 1 file | Multi-device filename |
| JSON     | Consolidated | 1 file | Multi-device filename |
| CSV      | Individual | 1 per device | Single-device filename |

## Usage

### Basic Usage

```powershell
# Generate Markdown documentation with automatic file naming
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat Markdown

# Generate Markdown with unused ports included (sequential ordering)
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat Markdown -ShowUnused

# Generate individual CSV files per device with comprehensive metadata
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat CSV

# Generate individual CSV files per device including unused ports  
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat CSV -ShowUnused

# Generate single JSON file for integration with other tools
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat JSON

# Note: CSV format creates individual files per device automatically
# -OutputFile parameter is ignored for CSV format as it creates multiple files
```

### Command Line Parameters

- `-InputFile`: Path to the input JSON configuration file (required)
- `-OutputFormat`: Output format - Markdown, CSV, or JSON (required)
- `-OutputFile`: Custom output file path (applies to Markdown/JSON only; CSV creates individual files automatically)
- `-ShowUnused`: Include unused ports in sequential order within connection tables/files
- `-LogLevel`: Logging verbosity - Info, Warning, Error (default: Info)

### Automatic File Naming

The tool generates intelligent output filenames that include:

- Input filename (without extension)
- Device manufacturers (e.g., "Cisco-Dell")
- Device models (abbreviated for readability)
- Device count
- Tool identifier ("portmap")
- Incremental numbering to prevent overwrites

**Example**: `sample-network-config-Cisco-Dell-93180YC-FX3-N9K-C93360YC-FX2-S5248F-ON-3dev-portmap.md`

## Integration with AzureStack Network Switch Config Generator

The PortMap tool integrates seamlessly with the main project:

1. **JSON Output**: Can be consumed by the main switch configuration generator
2. **Template Integration**: Port mapping data can inform Jinja2 template generation
3. **Validation**: Ensures physical design matches logical configuration requirements
4. **Documentation**: Provides physical layer documentation alongside generated configurations

## Use Cases

### 1. Network Design Documentation

Create comprehensive physical layer documentation for network designs with professional formatting:

```powershell
.\PortMap.ps1 -InputFile "datacenter-design.json" -OutputFormat Markdown
# Generates: datacenter-design-[manufacturers]-[models]-[count]dev-portmap.md
```

### 2. Port Utilization Analysis

Analyze complete port usage including unused ports in sequential order:

```powershell
.\PortMap.ps1 -InputFile "current-config.json" -OutputFormat Markdown -ShowUnused
# Shows all ports 1,2,3... with connected and unused ports integrated
```

### 3. Configuration Generator Input

Generate structured JSON data for switch configuration tools:

```powershell
.\PortMap.ps1 -InputFile "rack-layout.json" -OutputFormat JSON
# Creates machine-readable data with complete device and connection metadata
```

### 4. Device-Specific Spreadsheet Analysis

Export individual CSV files per device for detailed analysis and reporting:

```powershell
.\PortMap.ps1 -InputFile "network-inventory.json" -OutputFormat CSV -ShowUnused
# Creates one CSV file per device with complete port data and metadata
# Example outputs: TOR-1.csv, TOR-2.csv, SPINE-1.csv
```

### 5. Network Auditing

Document existing network configurations for compliance and auditing:

```powershell
.\PortMap.ps1 -InputFile "production-network.json" -OutputFormat Markdown -ShowUnused
# Professional documentation showing all ports, usage status, and availability
```

## Output Examples

### Sample Output Files

**Generated Filenames** (with automatic naming):

**Markdown Output (Consolidated):**

- `sample-network-config-Cisco-Dell-93180YC-FX3-N9K-C93360YC-FX2-S5248F-ON-3dev-portmap.md`

**JSON Output (Consolidated):**

- `datacenter-layout-Cisco-Dell-93180YC-FX3-N9K-C93360YC-FX2-S5248F-ON-3dev-portmap.json`

**CSV Output (Individual Files Per Device):**

- `sample-network-config-Cisco-93180YC-FX3-1dev-portmap.csv` (TOR-1 device)
- `sample-network-config-Dell-S5248F-ON-1dev-portmap.csv` (TOR-2 device)  
- `sample-network-config-Cisco-N9K-C93360YC-FX2-1dev-portmap.csv` (SPINE-1 device)

### Device Summary Table (All Formats)

| Device Name | Make  | Model        | Total Ports | Used Ports | Unused Ports |
| ----------- | ----- | ------------ | ----------- | ---------- | ------------ |
| TOR-1       | Cisco | 93180YC-FX3  | 52          | 28         | 24           |
| TOR-2       | Dell  | S5248F-ON    | 54          | 22         | 32           |
| SPINE-1     | Cisco | N9K-C93360YC | 96          | 0          | 96           |

### Sequential Port Mapping (Markdown with -ShowUnused)

```markdown
### TOR-1

**Device:** TOR-1 (Cisco 93180YC-FX3)
**Total Ports:** 52
**Connected Ports:** 28
**Unused Ports:** 24

| Port | Media  | Status     | Destination Device | Destination Port | Destination Media | Type         | Notes                    |
| ---- | ------ | ---------- | ------------------ | ---------------- | ----------------- | ------------ | ------------------------ |
| 1    | SFP28  | Active     | Server-Nodes       | 0                | SFP28             | Server       | Compute nodes 1-16       |
| 2    | SFP28  | Active     | Server-Nodes       | 0                | SFP28             | Server       | Compute nodes 1-16       |
| 25   | SFP28  | Active     | LoadBalancer-1     | 0                | SFP28             | Network      | Primary load balancer    |
| 26   | SFP28  | **Unused** | -                  | -                | -                 | -            | Available                |
| 27   | SFP28  | **Unused** | -                  | -                | -                 | -            | Available                |
| 49   | QSFP28 | Active     | SPINE-1            | 1                | QSFP28            | Uplink       | Primary uplinks to spine |
| 50   | QSFP28 | Active     | SPINE-1            | 2                | QSFP28            | Uplink       | Primary uplinks to spine |
| 51   | QSFP28 | Active     | TOR-2              | 51               | QSFP28            | Inter-Switch | MLAG peer link           |
| 52   | QSFP28 | **Unused** | -                  | -                | -                 | -            | Available                |
```

### Embedded Metadata (All Formats)

All outputs include comprehensive metadata:

- Generation timestamp and tool version
- Input file name and parameters used
- Device count, manufacturers, and models
- Generated by user and computer information
- Processing summary with port utilization statistics

## Error Handling & Logging

The tool includes comprehensive error handling and logging:

### Error Detection

- Invalid JSON input files with detailed syntax error reporting
- Missing or malformed device configurations
- Port range conflicts and overlapping assignments
- Media type mismatches between source and destination
- Invalid output format specifications
- File permission and path validation issues

### Logging Features

- Configurable logging levels (Info, Warning, Error)
- Detailed processing information with timestamps
- File operation status and output location details
- Processing statistics and validation results
- Comprehensive error messages with resolution guidance

### Robust Operation

- Graceful handling of missing optional fields
- Safe array and object property access (PowerShell best practices)
- Automatic creation of output directories
- Incremental file naming to prevent overwrites
- Input validation with clear error messages

## Testing & Validation

### Automated Testing

The tool includes comprehensive automated testing:

```powershell
# Run the included test suite
.\Test-PortMap.ps1

# Manual testing with sample configuration
.\PortMap.ps1 -InputFile "sample-network-config.json" -OutputFormat Markdown -ShowUnused
```

### Validation Features

- JSON schema validation for input configurations
- Port range overlap detection
- Media type compatibility checking
- Output file integrity verification
- Processing statistics validation

### Sample Configurations

The tool includes `sample-network-config.json` demonstrating:

- Multiple device types (Cisco and Dell switches)
- Different port configurations (SFP28, QSFP28)
- Various connection types (server, storage, uplink, inter-switch)
- Mixed utilized and unused ports

## Contributing

This tool is part of the AzureStack Network Switch Config Generator project. Please follow the project's contribution guidelines and coding standards.

### Development Guidelines

- Maintain PowerShell best practices and StrictMode compatibility
- Include comprehensive comment-based help for all functions
- Add appropriate error handling and logging
- Update tests for new functionality
- Follow consistent code formatting and naming conventions

## Support

For issues and questions:

1. Check the main project's troubleshooting guide
2. Review input JSON format requirements
3. Validate device configuration syntax
4. Consult the project's support documentation

## Version History

### **v1.0** (Current)

**Core Features:**

- Complete port mapping functionality with device inventory
- Multiple output formats: Markdown, CSV, and JSON with consistent formatting
- Professional documentation generation with embedded metadata

**Advanced Features:**

- Sequential port ordering with connected and unused ports integrated
- Intelligent automatic file naming with device and parameter information
- Incremental file management to prevent overwrites
- Comprehensive metadata embedding in all output formats
- Robust error handling and configurable logging

**PowerShell Best Practices:**

- StrictMode compatibility and safe property access
- Comment-based help with detailed parameter documentation
- Proper error handling with try/catch blocks
- Efficient data structures and processing algorithms
- Cross-platform PowerShell 5.1+ and Core 7.0+ support

**Professional Output:**

- Network engineering standard documentation format
- Per-device organization with clear status indicators
- Complete processing statistics and generation details
- Integration-ready JSON output with full device context

---

**Note**: This tool is designed to work alongside the main AzureStack Network Switch Config Generator and follows the same architectural patterns and coding standards.
