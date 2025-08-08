# Network Config Generator - Binary Releases

## Download & Usage

### Quick Start
1. Download the appropriate executable for your OS from the [Releases](../../releases) page
2. Make it executable (Linux): `chmod +x network-config-generator`
3. Run: `./network-config-generator --help`

### Available Executables
- **Windows**: `network-config-generator-windows-amd64.exe`
- **Linux**: `network-config-generator-linux-amd64`

### Usage Examples

#### Convert lab format and generate configs
```bash
./network-config-generator --input_json lab_input.json --output_folder output/
```

#### Use standard format directly
```bash
./network-config-generator --input_json standard_input.json --output_folder output/
```

#### Use custom convertor
```bash
./network-config-generator --input_json lab_input.json --convertor my.custom.convertor
```

### Features
- ✅ Auto-detects input format (lab vs standard)
- ✅ Converts lab format to standard format automatically  
- ✅ Generates network switch configurations
- ✅ Supports custom convertors
- ✅ Multi-switch support
- ✅ Cross-platform executables
- ✅ No Python installation required

### Building from Source
If you prefer to build the executable yourself:

```bash
pip install pyinstaller
pip install -r requirements.txt
pyinstaller network_config_generator.spec
```

The executable will be created in the `dist/` directory.
