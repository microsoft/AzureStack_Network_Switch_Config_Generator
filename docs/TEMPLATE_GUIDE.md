# Template Customization Guide

## 🎨 Overview

This guide shows you how to create and customize Jinja2 templates for generating network switch configurations. The template system is designed to be flexible, extensible, and vendor-agnostic.

## 📁 Template Structure

Templates are organized by vendor and firmware:

```
input/jinja2_templates/
├── cisco/
│   └── nxos/
│       ├── interfaces.j2
│       ├── vlans.j2
│       ├── bgp.j2
│       ├── prefix_lists.j2
│       └── qos.j2
├── dellemc/
│   └── os10/
│       ├── interfaces.j2
│       ├── vlans.j2
│       └── bgp.j2
└── your_vendor/            # Custom vendor support
    └── your_firmware/
        └── *.j2
```

## 🏷️ Template Discovery

The tool automatically discovers templates based on switch metadata:

```json
{
  "switch": {
    "make": "cisco",           # Maps to folder: cisco/
    "firmware": "nxos",        # Maps to folder: cisco/nxos/
    "model": "93180yc-fx",
    "hostname": "switch-1"
  }
}
```

**Template Path**: `input/jinja2_templates/cisco/nxos/*.j2`

## 📝 Basic Template Example

### Simple VLAN Template (`vlans.j2`)

```jinja2
{# vlans.j2 - Generate VLAN configuration #}
{% for vlan in vlans %}
vlan {{ vlan.vlan_id }}
  name {{ vlan.name }}
{% if vlan.interface is defined %}

interface vlan{{ vlan.vlan_id }}
  description {{ vlan.name }}
  ip address {{ vlan.interface.ip }}/{{ vlan.interface.cidr }}
{% if vlan.interface.mtu is defined %}
  mtu {{ vlan.interface.mtu }}
{% endif %}
{% if vlan.interface.redundancy is defined %}
  {{ vlan.interface.redundancy.type }} {{ vlan.interface.redundancy.group }}
  {{ vlan.interface.redundancy.type }} {{ vlan.interface.redundancy.group }} priority {{ vlan.interface.redundancy.priority }}
  {{ vlan.interface.redundancy.type }} {{ vlan.interface.redundancy.group }} ip {{ vlan.interface.redundancy.virtual_ip }}
{% endif %}
{% endif %}

{% endfor %}
```

### Output Example
```
vlan 100
  name Management

interface vlan100
  description Management
  ip address 192.168.1.1/24
  mtu 1500
  hsrp 100
  hsrp 100 priority 110
  hsrp 100 ip 192.168.1.254

vlan 200
  name Storage

interface vlan200
  description Storage
  ip address 10.0.0.1/24
  mtu 9000
```

## 🔧 Advanced Template Techniques

### 1. **Conditional Logic**

```jinja2
{# Interface configuration with conditions #}
{% for interface in interfaces %}
interface {{ interface.name }}
{% if interface.description is defined %}
  description {{ interface.description }}
{% endif %}
{% if interface.type == "Trunk" %}
  switchport mode trunk
{% if interface.native_vlan is defined %}
  switchport trunk native vlan {{ interface.native_vlan }}
{% endif %}
{% if interface.tagged_vlans is defined %}
  switchport trunk allowed vlan {{ interface.tagged_vlans }}
{% endif %}
{% elif interface.type == "Access" %}
  switchport mode access
  switchport access vlan {{ interface.vlan }}
{% elif interface.type == "L3" %}
  no switchport
{% if interface.ipv4 is defined and interface.ipv4 %}
  ip address {{ interface.ipv4 }}
{% endif %}
{% endif %}

{% endfor %}
```

### 2. **Loops and Filters**

```jinja2
{# BGP configuration with loops and filters #}
router bgp {{ bgp.asn }}
  router-id {{ bgp.router_id }}
  
  {# Advertise networks #}
{% for network in bgp.networks %}
{% if network %}
  network {{ network }}
{% endif %}
{% endfor %}

  {# Configure neighbors #}
{% for neighbor in bgp.neighbors %}
  neighbor {{ neighbor.ip }}
    remote-as {{ neighbor.remote_as }}
    description {{ neighbor.description }}
{% if neighbor.update_source is defined %}
    update-source {{ neighbor.update_source }}
{% endif %}
{% if neighbor.ebgp_multihop is defined %}
    ebgp-multihop {{ neighbor.ebgp_multihop }}
{% endif %}
    address-family ipv4 unicast
{% if neighbor.af_ipv4_unicast.prefix_list_in is defined %}
      prefix-list {{ neighbor.af_ipv4_unicast.prefix_list_in }} in
{% endif %}
{% if neighbor.af_ipv4_unicast.prefix_list_out is defined %}
      prefix-list {{ neighbor.af_ipv4_unicast.prefix_list_out }} out
{% endif %}

{% endfor %}
```

### 3. **Template Macros for Reusability**

```jinja2
{# macros.j2 - Reusable template components #}

{# Macro for interface common configuration #}
{% macro interface_common(interface) %}
interface {{ interface.name }}
{% if interface.description is defined %}
  description {{ interface.description }}
{% endif %}
{% if interface.mtu is defined %}
  mtu {{ interface.mtu }}
{% endif %}
{% endmacro %}

{# Macro for HSRP/VRRP configuration #}
{% macro redundancy_config(redundancy) %}
{% if redundancy.type == "hsrp" %}
  hsrp {{ redundancy.group }}
  hsrp {{ redundancy.group }} priority {{ redundancy.priority }}
  hsrp {{ redundancy.group }} ip {{ redundancy.virtual_ip }}
{% elif redundancy.type == "vrrp" %}
  vrrp {{ redundancy.group }}
  vrrp {{ redundancy.group }} priority {{ redundancy.priority }}
  vrrp {{ redundancy.group }} ip {{ redundancy.virtual_ip }}
{% endif %}
{% endmacro %}

{# Usage in main template #}
{% from 'macros.j2' import interface_common, redundancy_config %}

{% for vlan in vlans %}
{% if vlan.interface is defined %}
{{ interface_common(vlan.interface) }}
  ip address {{ vlan.interface.ip }}/{{ vlan.interface.cidr }}
{% if vlan.interface.redundancy is defined %}
{{ redundancy_config(vlan.interface.redundancy) }}
{% endif %}
{% endif %}
{% endfor %}
```

## 🏢 Vendor-Specific Customizations

### Cisco NX-OS Example

```jinja2
{# cisco/nxos/system.j2 #}
hostname {{ switch.hostname }}

feature bgp
feature interface-vlan
feature hsrp
feature lacp

ip domain-lookup
vdc {{ switch.hostname }} id 1

{# NX-OS specific features #}
feature nxapi
feature bash-shell
feature scp-server
```

### Dell EMC OS10 Example

```jinja2
{# dellemc/os10/system.j2 #}
hostname {{ switch.hostname }}

{# OS10 specific configuration #}
system-user linuxadmin password $6$rounds=656000$... role sysadmin 
aaa authorization exec default local

{# Enable required features #}
router bgp
interface vlan
```

## 🔄 Template Inheritance

Create base templates and extend them for specific vendors:

### Base Template (`base/interfaces.j2`)

```jinja2
{# Base interface template #}
{% for interface in interfaces %}
{% block interface_header %}
interface {{ interface.name }}
{% endblock %}

{% block interface_description %}
{% if interface.description is defined %}
  description {{ interface.description }}
{% endif %}
{% endblock %}

{% block interface_specific %}
{# Override in vendor-specific templates #}
{% endblock %}

{% block interface_footer %}
{# Common footer configuration #}
{% endblock %}

{% endfor %}
```

### Vendor-Specific Extension (`cisco/nxos/interfaces.j2`)

```jinja2
{# Cisco NX-OS specific interface template #}
{% extends "base/interfaces.j2" %}

{% block interface_specific %}
{% if interface.type == "Trunk" %}
  switchport mode trunk
{% if interface.native_vlan is defined %}
  switchport trunk native vlan {{ interface.native_vlan }}
{% endif %}
{% if interface.tagged_vlans is defined %}
  switchport trunk allowed vlan {{ interface.tagged_vlans }}
{% endif %}
{% elif interface.type == "Access" %}
  switchport mode access
  switchport access vlan {{ interface.vlan }}
{% elif interface.type == "L3" %}
  no switchport
{% if interface.ipv4 is defined and interface.ipv4 %}
  ip address {{ interface.ipv4 }}
{% endif %}
  no shutdown
{% endif %}
{% endblock %}
```

## 📊 Data Access Patterns

### Accessing Switch Data

```jinja2
{# Switch information #}
Hostname: {{ switch.hostname }}
Make: {{ switch.make }}
Model: {{ switch.model }}
Firmware: {{ switch.firmware }}
Version: {{ switch.version }}
```

### Working with Arrays

```jinja2
{# Loop through VLANs #}
{% for vlan in vlans %}
VLAN {{ vlan.vlan_id }}: {{ vlan.name }}
{% endfor %}

{# Filter arrays #}
{% for vlan in vlans if vlan.vlan_id > 100 %}
High VLAN: {{ vlan.vlan_id }}
{% endfor %}
```

### Nested Data Access

```jinja2
{# Access nested BGP data #}
BGP ASN: {{ bgp.asn }}
Router ID: {{ bgp.router_id }}

{% for neighbor in bgp.neighbors %}
Neighbor {{ neighbor.ip }}:
  AS: {{ neighbor.remote_as }}
  Description: {{ neighbor.description }}
{% if neighbor.af_ipv4_unicast is defined %}
  IPv4 Unicast:
{% for key, value in neighbor.af_ipv4_unicast.items() %}
    {{ key }}: {{ value }}
{% endfor %}
{% endif %}
{% endfor %}
```

## 🧪 Testing Templates

### Template Testing Approach

1. **Create test standard JSON files**
2. **Generate configs using your templates**
3. **Validate the output**

### Example Test Structure

```
tests/template_tests/
├── my_vendor_test/
│   ├── input/
│   │   └── test_standard_input.json
│   ├── templates/
│   │   └── my_vendor/
│   │       └── my_firmware/
│   │           ├── interfaces.j2
│   │           └── vlans.j2
│   └── expected_output/
│       ├── expected_interfaces.cfg
│       └── expected_vlans.cfg
```

### Template Testing Command

```bash
# Test your custom templates
python src/main.py \
  --input_json tests/template_tests/my_vendor_test/input/test_standard_input.json \
  --template_folder tests/template_tests/my_vendor_test/templates \
  --output_folder test_output/
```

## 🚀 Best Practices

### 1. **Use Meaningful Comments**
```jinja2
{# This section configures BGP neighbors for eBGP peering #}
{% for neighbor in bgp.neighbors if neighbor.remote_as != bgp.asn %}
```

### 2. **Handle Missing Data Gracefully**
```jinja2
{% if interface.mtu is defined and interface.mtu != 1500 %}
  mtu {{ interface.mtu }}
{% endif %}
```

### 3. **Use Consistent Formatting**
```jinja2
{# Good: Consistent indentation and spacing #}
interface {{ interface.name }}
  description {{ interface.description }}
  switchport mode {{ interface.mode }}

{# Bad: Inconsistent formatting #}
interface {{interface.name}}
description {{interface.description}}
 switchport mode {{interface.mode}}
```

### 4. **Leverage Jinja2 Filters**
```jinja2
{# String manipulation #}
hostname {{ switch.hostname | upper }}

{# List manipulation #}
{% for vlan in vlans | sort(attribute='vlan_id') %}

{# Default values #}
mtu {{ interface.mtu | default(1500) }}
```

### 5. **Modular Template Design**
Break large templates into smaller, manageable pieces:

```jinja2
{# main_config.j2 #}
{% include 'system.j2' %}
{% include 'interfaces.j2' %}
{% include 'vlans.j2' %}
{% include 'bgp.j2' %}
```

## 🔧 Troubleshooting Templates

### Common Issues

1. **Template Not Found**
   - Check template path: `input/jinja2_templates/{make}/{firmware}/`
   - Verify file extension: `.j2`

2. **Variable Not Defined Errors**
   ```jinja2
   {# Safe access with default #}
   {{ interface.mtu | default(1500) }}
   
   {# Check if defined #}
   {% if interface.mtu is defined %}
   ```

3. **Unexpected Output**
   - Use `{{ variable | pprint }}` for debugging
   - Check data types and structure

### Template Debugging

```jinja2
{# Debug template - shows all available data #}
<!-- DEBUG: Switch Data -->
{{ switch | pprint }}

<!-- DEBUG: VLANs Data -->
{{ vlans | pprint }}

<!-- DEBUG: Interfaces Data -->  
{{ interfaces | pprint }}
```

This template system provides powerful flexibility while maintaining simplicity. Start with basic templates and gradually add more sophisticated features as needed!
