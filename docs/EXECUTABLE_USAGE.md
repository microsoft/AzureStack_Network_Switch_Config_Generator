# Using the Standalone Executable

## 📦 No Python Required!

The easiest way to use this tool is with our pre-built executables. No need to install Python or any dependencies!

## ⬇️ Download

1. Go to [Releases](../../releases) page
2. Download for your operating system:
   - **Windows**: `network_config_generator.exe`
   - **Linux**: `network_config_generator`

## 🚀 Quick Start (3 Steps)

### Step 1: Make it executable (Linux only)
```bash
chmod +x network_config_generator
```

### Step 2: Prepare your input file
- Any JSON file with your switch data works
- Tool automatically detects the format

### Step 3: Generate configs
```bash
# Windows
.\network_config_generator.exe --input_json your_switches.json --output_folder configs\

# Linux  
./network_config_generator --input_json your_switches.json --output_folder configs/
```

## 📋 What You Get

After running, you'll see folders like this:
```
configs/
├── switch-01/
│   ├── generated_interfaces
│   ├── generated_vlans  
│   ├── generated_bgp
│   └── generated_full_config
└── switch-02/
    ├── generated_interfaces
    └── ... (same files)
```

## ⚙️ Advanced Options

### Custom Data Format
If your data isn't in the standard format:
```bash
./network_config_generator --input_json custom_data.json --convertor my.custom.convertor
```

### Help
```bash
./network_config_generator --help
```

## 🆘 Need Help?

- **Issues?** → See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Custom data format?** → See [CONVERTOR_GUIDE.md](CONVERTOR_GUIDE.md)
- **Modify templates?** → See [TEMPLATE_GUIDE.md](TEMPLATE_GUIDE.md)
