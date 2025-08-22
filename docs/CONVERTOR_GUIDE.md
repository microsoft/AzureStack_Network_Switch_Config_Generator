# Converting Your Data Format

## 🤔 Do I Need This?

**You need a converter if:**
- Your switch data is in a different format than the tool expects
- You get an error like "conversion required"
- Your JSON structure looks different from the examples

**You DON'T need this if:**
- Your data already works with the tool
- You're using the provided example files

## � Simple Example

Let's say your data looks like this:

**Your Format:**
```json
{
  "devices": [
    {
      "name": "switch-01",
      "type": "cisco",
      "vlans": "100,200,300",
      "ports": "1-48"
    }
  ]
}
```

**Tool Needs:**
```json
{
  "switch": {
    "hostname": "switch-01",
    "make": "cisco",
    "firmware": "nxos"
  },
  "vlans": [
    {"vlan_id": 100, "name": "VLAN_100"},
    {"vlan_id": 200, "name": "VLAN_200"},
    {"vlan_id": 300, "name": "VLAN_300"}
  ],
  "interfaces": [],
  "port_channels": [],
  "bgp": {}
}
```

## �️ Create Your Converter (4 Steps)

### Step 1: Create the File
Create: `src/convertors/my_converter.py`

### Step 2: Write the Function
```python
def convert_switch_input_json(input_data, output_directory):
    """
    Convert your format to standard format
    
    Args:
        input_data: Your JSON data (as Python dict)
        output_directory: Where to save converted files
    """
    
    # Example: Convert each device
    for device in input_data["devices"]:
        
        # Build standard format
        standard_data = {
            "switch": {
                "hostname": device["name"],
                "make": device["type"],
                "firmware": "nxos"  # or detect from your data
            },
            "vlans": [],
            "interfaces": [],
            "port_channels": [],
            "bgp": {}
        }
        
        # Convert VLANs (example: "100,200,300" → list)
        if device.get("vlans"):
            vlan_ids = device["vlans"].split(",")
            for vlan_id in vlan_ids:
                standard_data["vlans"].append({
                    "vlan_id": int(vlan_id),
                    "name": f"VLAN_{vlan_id}"
                })
        
        # Save file for this switch
        import json
        import os
        filename = f"{device['name']}.json"
        filepath = os.path.join(output_directory, filename)
        
        with open(filepath, 'w') as f:
            json.dump(standard_data, f, indent=2)
```

### Step 3: Test It
```bash
python src/main.py --input_json your_data.json --convertor my_converter
```

### Step 4: Check Results
Look in your output folder for:
- `switch-01/` (folder with generated configs)
- `std_switch-01.json` (converted data for troubleshooting)

## 📋 Standard Format Reference

Your converter must create JSON with these sections:

### Required: Switch Info
```json
{
  "switch": {
    "hostname": "switch-name",
    "make": "cisco",        // cisco, dellemc, etc.
    "firmware": "nxos"      // nxos, os10, etc.
  }
}
```

### Optional: VLANs
```json
{
  "vlans": [
    {
      "vlan_id": 100,
      "name": "VLAN_100",
      "description": "Server VLAN"
    }
  ]
}
```

### Optional: Interfaces  
```json
{
  "interfaces": [
    {
      "name": "Ethernet1/1",
      "vlan": 100,
      "description": "Server port",
      "type": "access"     // access, trunk, etc.
    }
  ]
}
```

### Optional: BGP, Port Channels, etc.
```json
{
  "bgp": {},
  "port_channels": [],
  "prefix_lists": [],
  "qos": {}
}
```

## 🆘 Common Issues

**Error: "Module not found"**
- Make sure file is in `src/convertors/`
- Use dots in path: `convertors.my_converter`

**Error: "Function not found"**  
- Function must be named exactly: `convert_switch_input_json`

**No output files**
- Check file permissions in output directory
- Add print statements to debug your converter

## 💡 Tips

- **Start simple**: Convert just the basic switch info first
- **Use print()**: Add debug prints to see what's happening  
- **Check examples**: Look at existing converters in `src/convertors/`
- **Test incrementally**: Add one section at a time (VLANs, then interfaces, etc.)
  "vlans": [
    {
      "vlan_id": 100,             // VLAN number
      "name": "Production"        // VLAN name
    }
  ]
}
```

### Optional: Interfaces  
```json
{
  "interfaces": [
    {
      "name": "Ethernet1/1",      // Interface name
      "mode": "access",           // access or trunk
      "vlan": 100                 // VLAN for access ports
    }
  ]
}
```

## 🎯 Complete Working Example

Let's say you have this input:
```json
{
  "datacenter": {
    "leaf_switches": [
      {
        "hostname": "leaf-01",
        "brand": "cisco",
        "software": "nxos", 
        "access_vlans": [100, 200]
      }
    ]
  }
}
```

Your convertor (`convertors/datacenter_convertor.py`) would:

1. **Extract switch data** from `datacenter.leaf_switches`
2. **Map fields**: `brand` → `make`, `software` → `firmware`  
3. **Transform VLANs**: Convert `access_vlans` array to VLAN objects
4. **Save file**: Create `leaf-01.json` with standard format

## 🚀 Using Your Convertor

```bash
# Use your custom convertor
./network-config-generator --input_json datacenter.json --convertor datacenter_convertor

# Tool will:
# 1. Load your datacenter.json
# 2. Run datacenter_convertor to transform it
# 3. Generate switch configs using templates
```

## 💡 Tips

- **Start Simple**: Begin with just switch info and VLANs
- **Test Often**: Check that your output works with the tool
- **Use Examples**: Look at existing convertors for reference
- **One File Per Switch**: Each switch gets its own JSON file
- **Validate Data**: Check for required fields before processing

## 📚 Need Help?

- Look at `convertors/convertors_lab_switch_json.py` for a real example
- Check `tests/test_cases/` for sample inputs and outputs
- The tool will show errors if your format doesn't match

That's it! Convertors are just simple data transformers - nothing complex needed!
