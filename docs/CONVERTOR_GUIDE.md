# Custom Input Format Guide

## 🎯 What is a Convertor?

If your switch data is in a different format than what the tool expects, you can create a **convertor** to transform it. Think of it as a translator between your format and the tool's format.

## 📝 Simple Example

**Your Input Format:**
```json
{
  "switches": [
    {
      "name": "leaf-01",
      "vendor": "cisco", 
      "os": "nxos",
      "vlans": "100,200,300"
    }
  ]
}
```

**Tool's Expected Format:**
```json
{
  "switch": {
    "hostname": "leaf-01",
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

## 🔧 How to Create Your Convertor

### Step 1: Place Your Convertor File
Put your convertor in: `convertors/my_custom_convertor.py`

### Step 2: Required Structure
Your convertor file must have this function:
```
convert_switch_input_json(input_data, output_directory)
```

### Step 3: What It Should Do
1. Read your input format
2. Transform it to the standard format
3. Save one JSON file per switch in the output directory

### Step 4: Use It
```bash
./network-config-generator --input_json your_data.json --convertor my_custom_convertor
```

## 📋 Standard Format Reference

Your convertor should create files with these sections:

### Required: Switch Info
```json
{
  "switch": {
    "hostname": "switch-name",    // Switch name
    "make": "cisco",              // cisco, dellemc, etc.  
    "firmware": "nxos"            // nxos, os10, etc.
  }
}
```

### Optional: VLANs
```json
{
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
