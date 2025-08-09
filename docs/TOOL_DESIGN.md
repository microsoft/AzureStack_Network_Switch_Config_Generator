# Tool Design

## 🏗️ Architecture Overview

The Network Switch Config Generator follows a modular, pipeline-based architecture that separates concerns and enables extensibility.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌──────────────────┐
│   Input Layer   │───▶│ Conversion Layer │───▶│ Generation Layer│───▶│   Output Layer   │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └──────────────────┘
        │                       │                       │                       │
        │                       │                       │                       │
    ┌───▼────┐               ┌───▼────┐               ┌───▼────┐               ┌───▼────┐
    │Lab JSON│               │Standard│               │Template│               │.cfg    │
    │Std JSON│               │  JSON  │               │Engine  │               │Files   │
    │CSV/YAML│               │        │               │        │               │        │
    └────────┘               └────────┘               └────────┘               └────────┘
```

## 🎯 Core Components

### 1. **Input Processing**
- **Purpose**: Accept various input formats and normalize them
- **Supported Formats**:
  - **Standard Format**: Direct JSON with switch, VLAN, and interface definitions
  - **Lab Format**: Structured format with metadata and nested switch data
  - **Custom Formats**: Extensible through custom convertors

### 2. **Format Detection**
- **Automatic Detection**: Tool automatically identifies input format
  - **Standard Format**: Contains keys like `switch`, `vlans`, `interfaces`
  - **Lab Format**: Contains keys like `Version`, `Description`, `InputData`
- **Manual Override**: Users can specify custom convertors for specialized formats

### 3. **Conversion System**
- **Dynamic Loading**: Convertors are loaded at runtime using Python's `importlib`
- **Standard Interface**: All convertors must implement `convert_switch_input_json(data, output_dir)`
- **Extensible**: Users can create custom convertors for any input format

```python
def load_convertor(convertor_module_path):
    """
    Dynamically loads convertor modules:
    - convertors.convertors_lab_switch_json (default)
    - my.custom.convertor (user-defined)
    """
```

### 4. **Template Engine (`generator.py`)**
- **Pluggable Architecture**: Convertors can be dynamically loaded
- **Default Convertor**: Handles lab format to standard format conversion
- **Custom Convertors**: Users can specify alternative conversion logic
- **Error Handling**: Graceful failure with helpful error messages

### 4. **Template Engine**
- **Jinja2-Based**: Uses industry-standard templating with powerful features
- **Vendor Support**: Organized by manufacturer and firmware version
- **Template Discovery**: Automatically finds templates based on switch metadata
- **Rich Context**: Provides comprehensive data context to templates

### 5. **Configuration Generation**
- **Multi-Switch Support**: Generates configs for multiple switches from single input
- **File Organization**: Each switch gets its own configuration file
- **Naming Convention**: Files named by switch hostname or specified pattern

## 📁 Data Flow

### Input Processing Pipeline
```
User Input → Format Detection → Conversion (if needed) → Validation → Generation
```

### Template Resolution
```
Switch Metadata → Template Directory → Template Discovery → Rendering → Output Files
```

Example:
```
{"make": "cisco", "firmware": "nxos"} 
    ↓
input/jinja2_templates/cisco/nxos/
    ↓
[interfaces.j2, vlans.j2, bgp.j2, ...]
    ↓
[generated_interfaces.cfg, generated_vlans.cfg, ...]
```

## 🔌 Extensibility Points

### 1. **Custom Convertors**
Create convertors for any input format by implementing the required interface and placing them in the `convertors/` directory.

**Usage**:
```bash
./network-config-generator --input_json my_input.json --convertor my.custom.convertor
```

### 2. **Custom Templates**
Add support for new vendors or versions:

```
input/jinja2_templates/
├── cisco/
│   └── nxos/
├── dellemc/
│   └── os10/
└── my_vendor/          # Custom vendor
    └── my_os/          # Custom OS
        ├── interfaces.j2
        ├── vlans.j2
        └── bgp.j2
```

### 3. **Template Inheritance**
Templates support Jinja2 inheritance for code reuse:

```jinja2
{# base_interface.j2 #}
interface {{ interface.name }}
  description {{ interface.description }}
{% block interface_config %}{% endblock %}

{# cisco_interface.j2 #}
{% extends "base_interface.j2" %}
{% block interface_config %}
  switchport mode {{ interface.mode }}
  switchport access vlan {{ interface.vlan }}
{% endblock %}
```

## 🏛️ Standard Format Schema

The tool uses a standardized JSON schema for internal processing:

```json
{
  "switch": {
    "make": "cisco|dellemc|...",
    "firmware": "nxos|os10|...",
    "model": "93180yc-fx|s5248f-on|...",
    "hostname": "switch-hostname",
    "type": "TOR1|TOR2|BORDER|...",
    "version": "10.3.4"
  },
  "vlans": [
    {
      "vlan_id": 100,
      "name": "Management",
      "interface": {
        "ip": "192.168.1.1",
        "cidr": 24,
        "mtu": 9216,
        "redundancy": {
          "type": "hsrp|vrrp",
          "group": 100,
          "priority": 150,
          "virtual_ip": "192.168.1.254"
        }
      }
    }
  ],
  "interfaces": [...],
  "port_channels": [...],
  "bgp": {...},
  "prefix_lists": {...},
  "qos": {...}
}
```

## 🔄 Workflow States

### State 1: Input Detection
```python
if is_standard_format(data):
    # Skip conversion, use input directly
else:
    # Use convertor to transform to standard format
```

### State 2: Multi-Switch Handling
```python
for std_file in standard_format_files:
    # Create output directory for switch
    # Generate configs from templates
    # Copy standard JSON for troubleshooting
```

### State 3: Template Processing
```python
for template_file in template_directory:
    # Load Jinja2 template
    # Render with standard JSON data
    # Write to output file
```

## 🚀 Performance Considerations

### Parallel Processing
- Multiple switches processed sequentially (simplicity over complexity)
- Templates for each switch processed in parallel-safe manner
- File I/O optimized for large deployments

### Memory Efficiency
- Streaming JSON processing for large inputs
- Template caching within single execution
- Temporary files cleaned up automatically

### Error Recovery
- Continue processing other switches if one fails
- Detailed error reporting for debugging
- Graceful degradation on template errors

## 🔧 Configuration Management

### Environment Variables
```bash
TEMPLATE_PATH=/custom/templates/    # Override default template location
CONVERTOR_PATH=my.custom.convertor  # Default convertor override
DEBUG=1                            # Enable debug logging
```

### Config File Support (Future Enhancement)
```yaml
# config.yaml
default_template_folder: input/jinja2_templates
default_convertor: convertors.convertors_lab_switch_json
output_settings:
  create_subdirectories: true
  copy_standard_json: true
  file_permissions: 644
```

## 🧪 Testing Architecture

### Unit Tests
- `test_convertors.py`: Test conversion logic
- `test_generator.py`: Test template processing
- `test_main.py`: Test CLI interface

### Integration Tests
- End-to-end workflows with sample data
- Template rendering validation
- Multi-switch scenarios

### Test Data Management
```
tests/test_cases/
├── convert_lab_*/              # Conversion test cases
│   ├── lab_input.json
│   └── expected_outputs/
├── std_*/                      # Generation test cases
│   ├── std_input.json
│   └── expected_*.cfg
```

This architecture provides a solid foundation for network configuration generation while maintaining flexibility for future enhancements and customizations.
