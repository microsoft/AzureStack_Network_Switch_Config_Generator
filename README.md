
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
├── docs/
│   └── architecture.md             # Design documentation
├── input/
│   ├── standard_input.json         # Standard input file
│   └── templates/
│       └── cisco/
│           └── nxos/
│               ├── full_config.j2           # Merge all templates into one
│               ├── feature1.j2              # Default feature template
│               ├── feature2.j2
│               ├── system.j2
│               └── 10/                      # NX-OS version 10 specific templates
│                   ├── version_feature1.j2  # Versioned feature template
│                   └── version_system.j2    # Versioned system template
├── src/
│   ├── __init__.py
│   ├── main.py                     # Entry point for the tool
│   ├── convertor.py                # Converts various input formats
│   ├── generator.py                # Main generation logic
│   └── loader.py                   # Loads and parses input

├── tests/
│   ├── test_generator.py          # Unit tests for generator logic
│   ├── test_convertors.py         # Unit tests for input conversion
│   ├── test_cases/
│   │   ├── convert_switch_input_json/
│   │   │   └── convert_switch_input.json
│   │   ├── std_nxos_hyperconverged/
│   │   │   └── std_nxos_hyperconverged_input.json
│   │   ├── std_nxos_switched/
│   │   │   └── std_nxos_switched_input.json
│   │
├── requirements.txt               # Python dependencies

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
| `--input_std_json`   | **Required**. Path to the standard input JSON file.                      |
| `--template_folder`  | Optional. Path to the Jinja2 templates folder. <br>Default: `input/templates` |
| `--output_folder`    | Optional. Directory to save the generated configuration files. <br>Default: current working directory (`.`). |


### Option 1: Run with Python (for Python users)

If you're comfortable with Python, you can clone the repository and run the tool directly using the source code:

```bash
# Example
python src/main.py --input_std_json your_standard_input.json --template_folder templates/ --output_folder outputs
```

### Option 2: Use the Precompiled .exe (no Python needed)

A standalone executable (`.exe`) is available in the [Releases](https://github.com/YourRepo/releases) section.  
This version is compiled using **PyInstaller** and does **not require Python** or any additional dependencies to run.

```bash
# Example
.\network_config_generator.exe --input_std_json your_standard_input.json --template_folder templates/ --output_folder outputs
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


