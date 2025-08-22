# Using the Standalone Executable

## 📦 No Python Required!

The easiest way to use this tool is with our pre-built executables. No need to install Python or any dependencies!

## ⬇️ Download

1. Go to [Releases](../../releases) page
2. Download for your operating system:
   - **Windows**: `network_config_generator.exe`
   - **Linux**: `network_config_generator`

## 🚀 Quick Start

### Step 1: Prepare your input file
- Any JSON file with your switch data works
- Tool automatically detects the format

### Step 2: Generate configs
```bash
# Windows
.\network_config_generator.exe --input_json your_switches.json --output_folder configs\

# Linux
# Make it executable (Linux only)
chmod +x network_config_generator 
./network_config_generator --input_json your_switches.json --output_folder configs/
```


## 📋 What You Get

After running, if successful, you'll see:
- Individual feature configs (e.g., interfaces, VLANs, BGP)
- Merged full config (all switch features combined into one configuration file)
- Standard switch JSON (for debugging): Contains the normalized switch data used to generate configs, useful for troubleshooting or verifying input conversion.
  
```
configs/
├── switch-01/
│   ├── generated_interfaces
│   ├── generated_vlans  
│   ├── generated_bgp
│   └── generated_full_config
│   └── std_switch_json
└── switch-02/
    ├── generated_interfaces
    └── ... (same files)
```

## ⚙️ Advanced Options

To see all available options and usage details, run the executable with the `-h` or `--help` flag:

```bash
# Windows
.\network_config_generator.exe -h

# Linux
./network_config_generator -h
```

This will display a summary of command-line arguments, default values, and workflow steps.
```powershell
PS C:\Users\liunick\Downloads\config_test> .\network_config_generator.exe -h
usage: network_config_generator.exe [-h] --input_json INPUT_JSON [--template_folder TEMPLATE_FOLDER]
                                    [--output_folder OUTPUT_FOLDER] [--convertor CONVERTOR]

Network config generator - automatically detects input format and converts if needed, then generates configs.

options:
  -h, --help            show this help message and exit
  --input_json INPUT_JSON
                        Path to input JSON file (can be lab format or standard format)
  --template_folder TEMPLATE_FOLDER
                        Root folder containing vendor templates (default: input/jinja2_templates)
  --output_folder OUTPUT_FOLDER
                        Directory to save generated config files (default: current directory)
  --convertor CONVERTOR
                        Python module path for the convertor to use when input is not in standard format. Only used if
                        conversion is needed. (default: convertors.convertors_lab_switch_json)

Workflow: 1) Check if input is standard format 2) If not, convert using specified convertor 3) Generate config files
from standard format
```

## 🆘 Need Help?

- **Issues?** → See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Custom data format?** → See [CONVERTOR_GUIDE.md](CONVERTOR_GUIDE.md)
- **Modify templates?** → See [TEMPLATE_GUIDE.md](TEMPLATE_GUIDE.md)
