# Using the Standalone Executable

## 📦 No Python Required!

The easiest way to use this tool is with our pre-built executables. No need to install Python or any dependencies!

## ⬇️ Download

1. Go to [Releases](../../releases) page
2. Download for your operating system:
   - **Windows**: `network_config_generator.exe`
   - **Linux**: `network_config_generator`

## 🚀 Quick Start


To see all available options and usage details, run the executable with the `-h` or `--help` flag:

```powershell
# Windows
.\network_config_generator.exe -h

# Linux
./network_config_generator -h
```

This will display a summary of command-line arguments, default values, and workflow steps.
```powershell
> .\network_config_generator.exe -h
usage: network_config_generator.exe [-h] --input_json INPUT_JSON [--template_folder TEMPLATE_FOLDER]
                                    [--output_folder OUTPUT_FOLDER] [--convertor CONVERTOR]

options:
  -h, --help            show this help message and exit
  --input_json INPUT_JSON
                        Path to input JSON file (lab or standard format)
  --template_folder TEMPLATE_FOLDER
                        Folder containing Jinja2 templates (default: input/jinja2_templates)
  --output_folder OUTPUT_FOLDER
                        Directory to save generated configs (default: same directory as input file)
  --convertor CONVERTOR
                        Convertor to use for non-standard input formats (default:
                        convertors.convertors_lab_switch_json)

Examples:
  network_config_generator.exe --input_json input/standard_input.json --output_folder output/
  network_config_generator.exe --input_json my_lab_input.json --output_folder configs/ --convertor lab
```

![Network Config Generator Demo](./media/network-config-generator01.gif)


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


## 🆘 Need Help?

- **Issues?** → See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Custom data format?** → See [CONVERTOR_GUIDE.md](CONVERTOR_GUIDE.md)
- **Modify templates?** → See [TEMPLATE_GUIDE.md](TEMPLATE_GUIDE.md)
