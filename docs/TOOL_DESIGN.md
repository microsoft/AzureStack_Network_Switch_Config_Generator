# Tool Architecture & Design

## � Who Should Read This

**Developers and Contributors** who want to:
- Understand how the tool works internally
- Contribute new features
- Debug complex issues
- Extend the architecture

**Skip this if you just want to use the tool** - check out [EXECUTABLE_USAGE.md](EXECUTABLE_USAGE.md) instead.

## 🏗️ High-Level Architecture

The tool follows a simple 4-stage pipeline:

```
Input JSON → Convert → Generate → Output Configs
    │           │         │           │
    │           │         │           └─ Individual .cfg files
    │           │         └─ Jinja2 templates + JSON data  
    │           └─ Standard JSON format
    └─ Any JSON format
```

**What each stage does:**

1. **Input**: Accept any JSON format
2. **Convert**: Transform to standard format (if needed)
3. **Generate**: Apply Jinja2 templates to create configs
4. **Output**: Save configuration files

## 📂 Code Structure

```
src/
├── main.py           # Entry point - CLI interface
├── loader.py         # Load JSON files (handles PyInstaller)
├── generator.py      # Generate configs from templates
└── convertors/       # Convert different input formats
    └── *.py
```

### Key Files Explained

**`main.py`** - The orchestrator
- Parses command line arguments
- Detects input format (standard vs lab vs custom)
- Calls converter if needed
- Calls generator for each switch
- Manages output directories

**`loader.py`** - File handling
- Loads JSON files (with PyInstaller compatibility)
- Handles resource paths for bundled executables

**`generator.py`** - Template engine
- Finds appropriate templates for switch vendor/firmware
- Renders Jinja2 templates with switch data
- Saves generated configuration files

**`convertors/*.py`** - Format converters
- Each converter handles a specific input format
- Must implement: `convert_switch_input_json(data, output_dir)`
- Dynamically loaded at runtime

## 🔄 Data Flow

### 1. Input Processing
```python
# main.py
data = load_input_json(input_file)
is_standard = is_standard_format(data)

if not is_standard:
    # Convert to standard format
    convert_to_standard_format(input_file, temp_dir, converter)
```

### 2. Format Detection
```python
def is_standard_format(data):
    # Standard format has these keys
    standard_keys = {"switch", "vlans", "interfaces"}
    
    # Lab format has these keys  
    lab_keys = {"Version", "Description", "InputData"}
    
    return has_standard_keys and not has_lab_keys
```

### 3. Dynamic Converter Loading
```python
def load_convertor(module_path):
    # Import module at runtime
    module = importlib.import_module(module_path)
    
    # Get the conversion function
    return module.convert_switch_input_json
```

### 4. Template Discovery
```python
# generator.py
def find_templates(switch_data):
    make = switch_data["switch"]["make"]        # "cisco"
    firmware = switch_data["switch"]["firmware"] # "nxos"
    
    template_dir = f"templates/{make}/{firmware}/"
    return glob(f"{template_dir}/*.j2")
```

### 5. Config Generation
```python
def generate_config(template_file, switch_data):
    template = jinja2_env.get_template(template_file)
    return template.render(**switch_data)
```

## 🔧 Key Design Principles

### 1. **Modular Architecture**
- Each component has a single responsibility
- Easy to test individual parts
- Can swap out components (different converters, template engines)

### 2. **Dynamic Loading**
- Converters loaded at runtime based on user choice
- Templates discovered automatically based on switch metadata
- No hardcoded vendor/format support

### 3. **Format Agnostic Input**
- Tool doesn't care about input format
- Auto-detection for common formats
- User-extensible through custom converters

### 4. **Template-Driven Output**
- All configuration logic in templates, not code
- Easy to add new vendors/formats
- Non-developers can customize output

### 5. **Clean Separation**
- Input processing ↔ Business logic ↔ Output generation
- Data flows in one direction
- Each stage validates its input

## 🧪 Testing Strategy

### Unit Tests
```python
# Test individual functions
test_is_standard_format()
test_load_convertor()
test_template_discovery()
```

### Integration Tests  
```python
# Test end-to-end flows
test_lab_format_conversion()
test_standard_format_direct()
test_custom_convertor()
```

### Template Tests
```python
# Validate template output
test_vlan_template_cisco()
test_bgp_template_dellemc()
```

## 🚀 Extension Points

### Adding New Input Formats
1. Create `convertors/my_format.py`
2. Implement `convert_switch_input_json(data, output_dir)`
3. Use with `--convertor my_format`

### Adding New Vendors
1. Create `templates/vendor/firmware/` directory
2. Add `.j2` template files
3. Set switch data: `"make": "vendor", "firmware": "firmware"`

### Adding New Features
1. **New CLI options**: Modify `main.py` argument parser
2. **New data fields**: Update template variable documentation
3. **New output formats**: Modify template rendering in `generator.py`

## 🐛 Debugging Tips

### Enable Debug Mode
```python
# Add to main.py
import logging
logging.basicConfig(level=logging.DEBUG)
```

### Trace Data Flow
```python
# Add debug prints at each stage
print(f"Input data keys: {data.keys()}")
print(f"Detected format: {is_standard_format(data)}")
print(f"Template path: {template_path}")
```

### Test Components Individually
```python
# Test converter without full pipeline
from convertors.my_converter import convert_switch_input_json
convert_switch_input_json(test_data, "debug_output/")
```

## 📈 Performance Considerations

### Large Input Files
- Converters process data in memory
- Consider streaming for very large datasets
- Template rendering is fast (cached compilation)

### Many Switches
- Each switch processed independently
- Could parallelize generation phase
- Template compilation cached automatically

### Executable Size
- PyInstaller bundles everything
- Exclude unnecessary dependencies
- Consider separate executables per platform
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
