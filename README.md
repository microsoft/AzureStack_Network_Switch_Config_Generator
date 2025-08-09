
# Network Configuration Generation Tool

## 📘 Overview

This tool aims to generate vendor-specific network switch configurations (e.g., Cisco NX-OS, Dell OS10) using JSON input and Jinja2 templates. It supports optional packaging for environments where Python is not pre-installed.

---

## 🎯 Goals

- Support configuration generation for **multiple switch vendors**
- Allow users to define **input variables** in a structured JSON format
- Output readable **vendor-specific configuration files**
- Support both **source-code use and standalone executable**
- Ensure clean project structure and developer scalability

---

## 🧱 Design Architecture

### Overall Flow
```mermaid
flowchart LR
    A["Standard Input JSON<br/>(Network Variables)"]
    C["Jinja2 Template<br/>(Config Templates)"]
    E("Config Generation Tool<br/>(Python Script)")
    G(Generated Configuration)

    A -->|Load Variables| E
    C -->|Provide Templates| E
    E -->|Render Template with Variables| G

    subgraph User Input
        direction TB
        A
        C
    end
```

> **Note:**  
> The structure and format of the **JSON input file must remain fixed** to match the variables used in the Jinja2 templates, but you can safely update **values** as needed, either manually or programmatically.

#### Other User-Defined Input Support

To support a wide range of input data formats, the system allows users to define their own converters. These converters transform any non-standard input into a unified JSON structure. Sample converters are provided in the repository as references to help users get started.

```mermaid
flowchart LR
    U1["Non-Standard JSON Input"]
    U2["CSV Input"]
    U3["YAML Input"]
    U4["Other Format Input"]
    S1["Standard Input JSON"]

    U1 -->|convertor1.py| S1
    U2 -->|convertor2.py| S1
    U3 -->|convertor3.py| S1
    U4 -->|convertorx.py| S1

    subgraph "User-Defined Input Types"
        direction TB
        U1
        U2
        U3
        U4
    end

```
Each input type should be handled by a user-defined converter script (e.g., convertor1.py). These scripts are responsible for converting the input into the standardized JSON format. Example converter scripts are included in the repo to illustrate expected structure and behavior.


### Workflow Detail
```mermaid
flowchart LR
    A["Standard Input JSON"]
    B["Config Generation Tool"]
    T1["vlan.j2 Template"]
    T2["bgp.j2 Template"]
    T3["interface.j2 Template"]
    Tn["xxx.j2 Templates"]
    O1["VLAN Config"]
    O2["BGP Config"]
    O3["Interface Config"]
    O4["xxx Config"]
    OC["Combined Full Config"]

    A --> T1
    A --> T2
    A --> T3
    A --> Tn
    T1 --> B
    T2 --> B
    T3 --> B
    Tn --> B
    B --> O1
    B --> O2
    B --> O3
    B --> O4
    B --> OC

    subgraph "Jinja2 Templates"
        direction TB
        T1
        T2
        T3
        Tn
    end

    subgraph "Dedicated Config Output"
        direction TB
        O1
        O2
        O3
        O4
    end

    subgraph "Full Config Output"
        direction TB
        OC
    end
```


## 🗂️ Directory Structure

```plaintext
root/
├── .devcontainer/                  # Development container configuration
├── .github/                        # GitHub Actions workflows
├── build/                          # Build artifacts and executables
├── docs/                           # Documentation files
│   ├── CONVERTOR_GUIDE.md          # Guide for creating custom convertors
│   ├── EXECUTABLE_USAGE.md         # Standalone executable usage guide
│   ├── TEMPLATE_GUIDE.md           # Jinja2 template development guide
│   ├── TOOL_DESIGN.md              # Architecture and design documentation
│   └── TROUBLESHOOTING.md          # Common issues and solutions
├── input/                          # Input files and templates
│   ├── standard_input.json         # Standard format input example
│   ├── jinja2_templates/           # Jinja2 template files
│   │   ├── cisco/
│   │   │   └── nxos/               # Cisco NX-OS templates
│   │   │       ├── bgp.j2          # BGP configuration template
│   │   │       ├── full_config.j2  # Combined full configuration template
│   │   │       ├── interface.j2    # Interface configuration template
│   │   │       ├── login.j2        # Login/user configuration template
│   │   │       ├── port_channel.j2 # Port channel configuration template
│   │   │       ├── prefix_list.j2  # Prefix list configuration template
│   │   │       ├── qos.j2          # QoS configuration template
│   │   │       ├── system.j2       # System configuration template
│   │   │       └── vlan.j2         # VLAN configuration template
│   │   └── dellemc/                # Dell EMC templates (vendor-specific)
│   └── switch_interface_templates/ # Switch interface template configurations
│       ├── cisco/
│       └── dellemc/
├── src/                            # Source code
│   ├── main.py                     # Entry point for the tool
│   ├── generator.py                # Main generation logic
│   ├── loader.py                   # Input file loading and parsing
│   └── convertors/                 # Input format converters
│       └── convertors_lab_switch_json.py  # Lab format to standard JSON converter
├── tests/                          # Test files
│   ├── test_generator.py           # Unit tests for generator logic
│   ├── test_convertors.py          # Unit tests for input conversion
│   └── test_cases/                 # Test case data files
├── requirements.txt                # Python dependencies
├── network_config_generator.spec   # PyInstaller build specification
└── _old/                           # Legacy Go implementation (archived)

```

---

## 🔧 Input (Example)

### Standard Input JSON (Example)
```json
{
  "hostname": "tor-switch-1",
  "interfaces": [
    { "name": "Ethernet1/1", "vlan": 711, "description": "Compute1" },
    { "name": "Ethernet1/2", "vlan": 712, "description": "Storage1" }
  ],
  "vlans": [
    { "id": 711, "name": "Compute" },
    { "id": 712, "name": "Storage" }
  ],
  "bgp": {
    "asn": 65001,
    "router_id": "192.168.0.1",
    "neighbors": [
      { "ip": "192.168.0.2", "remote_as": 65002 }
    ]
  }
}
```

---

### Input Jinja2 Template (Example)

Example: `templates/nxos/bgp.j2`

```jinja2
router bgp {{ bgp.asn }}
  router-id {{ bgp.router_id }}
{% for neighbor in bgp.neighbors %}
  neighbor {{ neighbor.ip }}
    remote-as {{ neighbor.remote_as }}
{% endfor %}
```

---

## 🚀 Quick Start

### 📥 Input Parameters

The tool accepts the following input parameters:

| Parameter            | Description                                                                 |
|----------------------|-----------------------------------------------------------------------------|
| `--input_json`       | **Required**. Path to the input JSON file (lab format or standard format). |
| `--output_folder`    | **Required**. Directory to save the generated configuration files.         |
| `--convertor`        | Optional. Custom convertor module (e.g., `my.custom.convertor`).          |

### Option 1: Run with Python (for Python users)

If you're comfortable with Python, you can clone the repository and run the tool directly using the source code:

```bash
# Auto-detect format and generate configs
python src/main.py --input_json your_input.json --output_folder outputs/

# Use with lab format (will be auto-converted)
python src/main.py --input_json lab_input.json --output_folder outputs/

# Use with standard format (direct processing)
python src/main.py --input_json standard_input.json --output_folder outputs/

# Use custom convertor for specialized formats
python src/main.py --input_json custom_input.json --output_folder outputs/ --convertor my.custom.convertor
```

### Option 2: Use the Precompiled Executable (no Python needed)

Standalone executables are available in the [Releases](https://github.com/microsoft/AzureStack_Network_Switch_Config_Generator/releases) section.  
These are compiled using **PyInstaller** and do **not require Python** or any additional dependencies.

```bash
# Windows
.\network-config-generator-windows-amd64.exe --input_json your_input.json --output_folder outputs\

# Linux
./network-config-generator-linux-amd64 --input_json your_input.json --output_folder outputs/
```

---

## Customization and Extensibility

### Templates

### Convertors

### Test Cases


---

## 🚀 What’s Improved in This Version

We’ve significantly redesigned the tool architecture to make it more modular, maintainable, and user-friendly. While the original version used **Go templates**, the new version is built with **Python + Jinja2** and brings several key improvements:

### ✅ Key Enhancements

- **Modular Output Generation**  
  Instead of producing a single full configuration, the tool now supports generating individual configuration sections (e.g., VLANs, interfaces, BGP) based on your input needs—making review, debugging, and reuse much easier.

- **No Rebuild Required for Changes**  
  All logic is now in editable Python and Jinja2 templates—no compilation needed. You can easily update the templates or logic without rebuilding any binary.

- **Cleaner and More Flexible Logic**  
  We’ve removed hardcoded rules and added more structured parsing logic, which makes the tool easier to extend, test, and adapt to different network designs or vendors.

- **Open to Contributions**  
  The new structure makes it easy for contributors to add new templates or enhance existing ones. If you support a new vendor or configuration style, you can simply submit a template pull request.

This upgrade is not just a technology shift—it’s a foundation for faster iteration, better collaboration, and easier maintenance.


