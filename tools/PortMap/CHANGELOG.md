# PortMap Tool - Changelog

## Version 1.0 Updates (October 2024)

### New Features

#### CSV Input Support
- **Added CSV input format** as an alternative to JSON
- Accepts two CSV files: `devices.csv` and `connections.csv`
- Flexible file naming patterns:
  - `{name}-devices.csv` + `{name}-connections.csv`
  - `devices.csv` + `connections.csv`
- Can specify either file as input—tool automatically finds companion file
- Full validation and error handling for CSV input
- Seamless conversion to internal format

#### Documentation Improvements
- **Quick Start section** moved to the top of README
- 5 common use cases with examples and expected outputs
- **CSV Input Format** dedicated section with:
  - Complete format specification
  - Required and optional columns
  - CSV vs JSON comparison table
  - Multiple usage examples
- Created `CSV_INPUT_EXAMPLES.md` with comprehensive examples
- Simplified Overview section for better accessibility
- Updated all sections to mention CSV support

### Files Changed

#### Modified Files
- `PortMap.ps1` - Added CSV input support
  - New function: `ConvertFrom-CsvToConfiguration`
  - Updated parameter validation to accept .csv files
  - Updated help documentation with CSV examples
  - Enhanced input format detection and processing
- `Test-PortMap.ps1` - Added CSV input tests
  - New function: `Test-CsvInputGeneration`
  - Tests for both devices.csv and connections.csv input
- `README.MD` - Restructured and expanded
  - Quick Start section at the beginning
  - CSV Input Format section with examples
  - Updated all relevant sections

#### New Files
- `sample-devices.csv` - Example devices configuration
- `sample-connections.csv` - Example connections configuration
- `CSV_INPUT_EXAMPLES.md` - Comprehensive CSV format guide
- `.gitignore` - Exclude generated output files

### Benefits

✅ **Easier for Network Engineers** - CSV format is familiar and editable in spreadsheet tools  
✅ **Quick Start** - Users can get started immediately with clear examples  
✅ **Flexible Input** - Both JSON and CSV supported  
✅ **Better Documentation** - Example-driven approach  
✅ **Backward Compatible** - All existing JSON workflows unchanged

### Usage Examples

#### JSON Input (Existing)
```powershell
.\PortMap.ps1 -InputFile "network-config.json" -OutputFormat Markdown
```

#### CSV Input (New)
```powershell
.\PortMap.ps1 -InputFile "sample-devices.csv" -OutputFormat Markdown
```

Both produce identical output formats (Markdown, CSV, or JSON).

### Migration Guide

**Existing JSON users:** No changes needed. All JSON functionality works exactly as before.

**New CSV users:** Create two CSV files:
1. `devices.csv` with device and port range information
2. `connections.csv` with connection mappings

See `CSV_INPUT_EXAMPLES.md` for complete format specification.

### Technical Details

- CSV parsing uses PowerShell's `Import-Csv` cmdlet
- Automatic device grouping from multiple port range rows
- Connection arrays built from CSV rows
- Validation ensures required fields are present
- Error messages guide users to correct issues

---

For detailed usage instructions, see [README.MD](README.MD).  
For CSV format examples, see [CSV_INPUT_EXAMPLES.md](CSV_INPUT_EXAMPLES.md).
