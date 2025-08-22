# Switch Interface Template Guide

This guide shows you how to quickly create or modify switch interface templates for the Azure Stack Network Configuration Generator. **No programming experience is required** - just follow the examples below.

## Quick Start: Common Use Cases

### Use Case 1: I Have a New Switch Model
**Goal:** Create a template for a switch model that doesn't exist yet.

**Steps:**
1. **Find the template folder:** `input/switch_interface_templates/`
   - Look for existing templates in `cisco/` or `dellemc/` folders.
   - If your switch is from another vendor, create a new folder (e.g., `juniper/`).
2. **Copy a similar template:** Pick an existing `.json` file closest to your switch
3. **Rename the file:** Change filename to your exact switch model (e.g., `9336C-FX2.json`)
4. **Update 3 key things:**
   - Switch model name
   - Interface ranges (port counts)
   - Uplink port numbers

**Example:** Creating template for Cisco 9336C-FX2:
```bash
Copy: cisco/93180YC-FX.json → cisco/9336C-FX2.json
```

Then edit these parts:
```json
"model": "9336C-FX2",                    // Change model name
"end_intf": "1/32",                      // Change port count
"intf": "1/35"                           // Change uplink ports
```

### Use Case 2: Change Host Connection Ports
**Goal:** Servers connect to different ports than the template shows.

**Find this section in the file:**
```json
"start_intf": "1/1",                     // Change start port
"end_intf": "1/16"                       // Change end port
```

**Change to your ports:**
```json
"start_intf": "1/5",                     // Servers start at port 5
"end_intf": "1/24"                       // Servers end at port 24
```

### Use Case 3: Change Uplink Ports
**Goal:** Your uplink connections are on different ports.

**Find these sections:**
```json
"name": "P2P_Border1",
"intf": "1/48"                           // Change this port

"name": "P2P_Border2", 
"intf": "1/47"                           // Change this port
```

**Change to your uplink ports:**
```json
"intf": "1/52"                           // New uplink port
"intf": "1/51"                           // New uplink port
```

### Use Case 4: Add or Change VLANs
**Goal:** Use different VLAN numbers or add more VLANs.

**Find VLAN settings:**
```json
"native_vlan": "99",                     // Untagged VLAN
"tagged_vlans": "M,C,S"                  // Tagged VLANs
```

**Recommendation:** Use actual VLAN numbers instead of letters:
```json
"native_vlan": "100",                    // Your management VLAN
"tagged_vlans": "200,300,400"            // Your VLANs
```

> **Note:** M=Management, C=Compute, S=Storage. Using numbers is easier.

### Use Case 5: Change Port Channel Members
**Goal:** Different interfaces for link aggregation.

**Find port channel section:**
```json
"port_channels": [
  {
    "id": 50,
    "members": ["1/41", "1/42"]          // Change these interfaces
  }
]
```

**Change to your interfaces:**
```json
"members": ["1/49", "1/50", "1/51"]      // Your link aggregation ports
```

## Interface Naming by Vendor

**Cisco:** `1/48` (slot/port)  
**Dell EMC:** `1/1/48` (unit/slot/port)

Make sure you use the right format for your switch brand!

## Quick Validation

After making changes, check these basics:

### 1. File Syntax
- Every `{` has a matching `}`
- Every `[` has a matching `]`
- Commas `,` between items (except the last item)
- Text values in quotes `"like this"`

### 2. Interface Numbers
- Verify interface numbers exist on your switch
- Use correct format: Cisco `1/48` vs Dell `1/1/48`
- No interface used in multiple places

### 3. Test Your Changes
1. Save your template file
2. Run the configuration generator
3. Check output files in `test_output` folder

## Quick Troubleshooting

**Problem:** Generator shows errors  
**Solution:** Check JSON syntax using VS Code or online JSON validator

**Problem:** Wrong interface configuration generated  
**Solution:** Verify interface names match your switch documentation

**Problem:** Overlapping configurations  
**Solution:** Ensure same interface isn't configured in multiple places

---

## Detailed Reference (Advanced Users)

### Template File Structure

Templates are stored in JSON format in this directory structure:

```
input/switch_interface_templates/
├── cisco/
│   ├── 93108TC-FX3P.json
│   ├── 93180YC-FX.json
│   └── 93180YC-FX3.json
└── dellemc/
    └── S5248F-ON.json
```

### Complete Template Sections

Each template file contains these main sections:

#### 1. Switch Information
```json
{
    "make": "Cisco",           // Manufacturer name
    "model": "93180YC-FX",     // Exact switch model
    "type": "TOR",             // Switch type (Top of Rack)
```

#### 2. Interface Templates
Organized by deployment scenarios:

```json
"interface_templates": {
    "common": [...],           // Interfaces used in all scenarios
    "fully_converged": [...],  // Hyper-converged infrastructure
    "switched": [...],         // Traditional switched network
    "switchless": [...]        // Switchless configuration
}
```

#### 3. Port Channels
```json
"port_channels": [...]        // Link aggregation groups
```

### Complete Interface Properties

| Property | Description | Example Values |
|----------|-------------|----------------|
| `name` | Descriptive name | `"HyperConverged_To_Host"` |
| `type` | Interface type | `"Access"`, `"Trunk"`, `"L3"` |
| `description` | Interface description | `"Connection to compute host"` |
| `intf_type` | Physical interface type | `"Ethernet"`, `"loopback"` |
| `intf` | Single interface | `"1/48"` |
| `start_intf` | Range start | `"1/1"` |
| `end_intf` | Range end | `"1/16"` |
| `access_vlan` | Access port VLAN | `"100"` |
| `native_vlan` | Trunk native VLAN | `"99"` |
| `tagged_vlans` | Trunk tagged VLANs | `"100,200,300"` |
| `shutdown` | Disable interface | `true` or `false` |
| `ipv4` | IP address | `""` (empty for dynamic) |
| `service_policy` | QoS policy | `{"qos_input": "AZLOCAL-QOS-MAP"}` |

### Deployment Scenarios Explained

#### Common Interfaces
Used in all deployment types:
- **Unused**: Shutdown all interfaces initially
- **Loopback0**: Management loopback interface
- **P2P_Border1/Border2**: Uplink connections
- **Trunk_TO_BMC_SWITCH**: BMC management connection

#### Fully Converged
- **HyperConverged_To_Host**: Trunk carrying management (M), compute (C), and storage (S) traffic

#### Switched
- **Switched_Compute_To_Host**: Trunk for compute traffic
- **Switched_Storage_To_Host**: Trunk for storage traffic

#### Switchless
- **Switchless_Compute_To_Host**: Simplified trunk configuration

### Port Channel Details

Port channels (link aggregation) configuration:

```json
{
    "id": 50,                           // Channel group number
    "description": "P2P_IBGP",          // Description
    "type": "L3",                       // "L3" or "Trunk"
    "ipv4": "",                         // IP address (empty for dynamic)
    "members": ["1/41", "1/42"]         // Physical interfaces in the group
}
```

For trunk port channels, add:
```json
"native_vlan": "99",                    // Untagged VLAN
"tagged_vlans": "100,200"               // Tagged VLANs
```
