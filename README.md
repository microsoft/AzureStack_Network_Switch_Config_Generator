
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

#### Customized Input Format

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
├── input/                          # All device-related inputs
│   ├── standard_input.json         # Device variables (standard format)
│   └── templates/                  # Jinja2 templates by vendor/OS
│       ├── cisco/
│       │   └── nxos/
│       │       ├── system.j2
│       │       └── ...
│       └── dell/
│           └── os10/
│               └── system.j2
│
├── src/                            # Python logic
│   ├── __init__.py
│   ├── convert_inputxx.py          # Convert different input formats  to standard JSON
│   ├── generator.py                # Load JSON, select template, render
│   └── loader.py                   # Input reading, file handling
│
├── tests/                          # Automated tests
│   ├── test_generator.py           # Unit tests: render logic
│   ├── test_templates.py           # Template correctness tests
│   └── test_cases/                 # Integration tests (input + expected output)
│       ├── case_01/
│       │   ├── input.json
│       │   └── expected.config
│       ├── case_02/
│       └── ...
│
├── main.py                         # Entry-point for manual use
├── requirements.txt                # Python dependencies
├── .gitignore                      # Files to exclude from Git
├── README.md                       # Project overview
└── docs/                           # Optional: design docs, Mermaid, notes
    └── architecture.md
```

---                       |

---

## 🔧 Input Format (Example)

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

## 🛠️ Template Structure (Jinja2)

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

## 📦 Packaging Strategy

- Use [`pyproject.toml`](https://python-poetry.org/docs/pyproject/) for dependency management
- Use [PyInstaller](https://pyinstaller.org/) to create standalone executables
- Build artifacts in `dist/` for users without Python installed

---

## 🛠️ Why We Switched: Go Templates → Python + Jinja2

We initially used **Golang + Go Templates** to generate switch configurations. It worked, but we found some limitations as the project grew. Now, we’ve switched to **Python + Jinja2** for better flexibility and maintainability.

### 🔍 Comparison Table

| Feature                        | Go + Go Templates                   | Python + Jinja2                         |
|-------------------------------|-------------------------------------|-----------------------------------------|
| Templating Features           | Basic, minimal logic                | Powerful logic, filters, macros         |
| Community & Ecosystem         | Smaller for templates               | Large and well-supported                |
| Flexibility                   | Harder to scale with many vendors   | Easy to handle multi-vendor templates   |
| Config File Support           | Manual parsing needed               | Native support for JSON, YAML, TOML     |
| Customer Customization        | Needs Go rebuild                    | Just edit input files or templates      |
| Packaging                     | `go build` (simple binary)          | `pyinstaller` (self-contained app)      |

### ✅ Why Python + Jinja2 Is Better

- Templates are **easier to read and maintain**.
- Supports **multiple vendors** without duplicating logic.
- Accepts **human-readable input files** (JSON, YAML).
- No need to recompile — customers can **edit templates directly**.
- Larger ecosystem and more reusable tools.

This change helps us move faster, reduce complexity, and make the tool more user-friendly.

