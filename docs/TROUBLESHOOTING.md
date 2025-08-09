# Troubleshooting Guide

## 🚨 Common Issues and Solutions

### Input Format Issues

#### **Error: "Input is in lab format - conversion required"**
**Problem**: The tool detected that your input is not in standard format.

**Solutions**:
1. **Use the default convertor** (recommended):
   ```bash
   ./network-config-generator --input_json your_lab_input.json --output_folder configs/
   ```

2. **Specify a custom convertor**:
   ```bash
   ./network-config-generator --input_json your_input.json --convertor my.custom.convertor
   ```

3. **Check your input format**:
   - Standard format should have: `{"switch": {...}, "vlans": [...], ...}`
   - Lab format typically has: `{"Version": "...", "Description": "...", "InputData": {...}}`

---

#### **Error: "Failed to convert to standard format"**
**Problem**: The convertor failed to process your input.

**Solutions**:
1. **Check convertor exists**: Verify that the convertor module is in the `convertors/` directory

2. **Verify input file structure**:
   ```json
   {
     "Version": "1.0.0",
     "Description": "...",
     "InputData": {
       "Switches": [...],
       "Supernets": [...]
     }
   }
   ```

3. **Check for missing required fields** in your input data

4. **Enable debug mode** by adding verbose output to understand the conversion process

---

### Template Issues

#### **Error: "Template path not found"**
**Problem**: The tool cannot find templates for your switch type.

**Solution**: Check template directory structure:
```bash
ls -la input/jinja2_templates/{make}/{firmware}/
# Example: input/jinja2_templates/cisco/nxos/
```

**Expected structure**:
```
input/jinja2_templates/
├── cisco/
│   └── nxos/
│       ├── interfaces.j2
│       ├── vlans.j2
│       └── bgp.j2
└── dellemc/
    └── os10/
        └── *.j2
```

**Fix**: Create the missing template directory or check your switch metadata:
```json
{
  "switch": {
    "make": "cisco",    // Must match folder name
    "firmware": "nxos"  // Must match subfolder name
  }
}
```

---

#### **Error: "No templates found in template directory"**
**Problem**: Template directory exists but contains no `.j2` files.

**Solutions**:
1. **Check file extensions**: Templates must end with `.j2`
2. **Verify file permissions**: Ensure files are readable
3. **Check template content**: Ensure templates are valid Jinja2

---

#### **Error: Template rendering failed**
**Problem**: Jinja2 template has syntax errors or undefined variables.

**Debugging steps**:
1. **Check template syntax**:
   ```jinja2
   {# Good #}
   {% for vlan in vlans %}
   vlan {{ vlan.vlan_id }}
   {% endfor %}
   
   {# Bad - missing endfor #}
   {% for vlan in vlans %}
   vlan {{ vlan.vlan_id }}
   ```

2. **Check variable names**:
   ```jinja2
   {# Safe access #}
   {{ interface.mtu | default(1500) }}
   
   {# Check if defined #}
   {% if interface.mtu is defined %}
   mtu {{ interface.mtu }}
   {% endif %}
   ```

3. **Add debug output to templates**:
   ```jinja2
   {# Debug: Show available data #}
   <!-- DEBUG: Available VLANs -->
   {{ vlans | pprint }}
   ```

---

### File Permission Issues

#### **Error: Permission denied when writing files**
**Problem**: Output directory is not writable.

**Solutions**:
1. **Check directory permissions**:
   ```bash
   ls -ld output_folder/
   ```

2. **Create directory with proper permissions**:
   ```bash
   mkdir -p output_folder/
   chmod 755 output_folder/
   ```

3. **Run with appropriate user permissions**

---

### Executable Issues (Binary Distributions)

#### **Error: "network-config-generator: command not found"**
**Problem**: Executable not in PATH or not executable.

**Solutions**:
1. **Make executable** (Linux):
   ```bash
   chmod +x network-config-generator-linux-amd64
   ```

2. **Use full path**:
   ```bash
   ./network-config-generator-linux-amd64 --help
   ```

3. **Add to PATH**:
   ```bash
   export PATH=$PATH:/path/to/executable
   ```

---

#### **Windows: "The system cannot execute the specified program"**
**Problem**: Windows security or missing dependencies.

**Solutions**:
1. **Unblock the file**:
   - Right-click executable → Properties → General → Unblock

2. **Run as Administrator** if needed

3. **Check Windows Defender** - ensure file isn't quarantined

4. **Use PowerShell**:
   ```powershell
   .\network-config-generator-windows-amd64.exe --help
   ```

---

### JSON Format Issues

#### **Error: "Failed to parse JSON"**
**Problem**: Input JSON file has syntax errors.

**Solutions**:
1. **Validate JSON syntax**:
   ```bash
   python -m json.tool your_input.json
   ```

2. **Check for common JSON errors**:
   - Missing commas: `{"key1": "value1" "key2": "value2"}`  ❌
   - Trailing commas: `{"key": "value",}`  ❌ 
   - Single quotes: `{'key': 'value'}`  ❌
   - Comments: `{"key": "value", // comment}`  ❌

3. **Use a JSON validator** online or in your editor

---

#### **Error: "Input JSON was empty or failed to parse"**
**Problem**: JSON file is empty, corrupted, or has encoding issues.

**Solutions**:
1. **Check file size and content**:
   ```bash
   ls -l your_input.json
   head your_input.json
   ```

2. **Check file encoding**:
   ```bash
   file your_input.json
   ```

3. **Fix encoding if needed**:
   ```bash
   iconv -f windows-1252 -t utf-8 input.json > input_utf8.json
   ```

---

### Performance Issues

#### **Problem: Generation takes too long**
**Solutions**:

1. **Reduce template complexity**:
   - Avoid complex loops in templates
   - Use simpler Jinja2 expressions

2. **Check input data size**:
   - Large VLAN lists or interface configurations can slow processing

3. **Optimize templates**:
   ```jinja2
   {# Inefficient #}
   {% for interface in interfaces %}
     {% for vlan in vlans %}
       {# Complex nested logic #}
     {% endfor %}
   {% endfor %}
   
   {# Better #}
   {% for interface in interfaces %}
     {# Simple operations #}
   {% endfor %}
   ```

---

#### **Problem: Large output files**
**Solutions**:

1. **Remove debug output** from templates
2. **Optimize template whitespace**:
   ```jinja2
   {# Control whitespace #}
   interface {{ interface.name }}
   {%- if interface.description is defined %}
     description {{ interface.description }}
   {%- endif %}
   ```

---

### Multi-Switch Issues

#### **Problem: Only one switch config generated**
**Cause**: Input contains multiple switches but convertor only processes one.

**Solutions**:
1. **Check convertor implementation** - ensure it processes all switches
2. **Verify input format** - ensure all switches are properly defined
3. **Check output directory** - files might be overwriting each other

---

#### **Problem: Switch configs overwriting each other** 
**Cause**: Multiple switches have the same hostname.

**Solutions**:
1. **Ensure unique hostnames** in input data
2. **Check convertor logic** for hostname generation
3. **Use switch type/role in filename** if hostnames conflict

---

## 🔧 Debugging Techniques

### Enable Verbose Output

1. **Check tool output**: Run the tool with verbose logging to see processing steps

2. **Add debug to templates**:
   ```jinja2
   {# Show all available data #}
   <!-- DEBUG START -->
   {{ data | pprint }}
   <!-- DEBUG END -->
   ```

### Check Intermediate Files

1. **Examine standard format JSON**:
   ```bash
   # After conversion, check the generated standard files
   cat temp_converted/*.json | jq '.'
   ```

2. **Verify template discovery**:
   ```bash
   ls -la input/jinja2_templates/{make}/{firmware}/
   ```

### Test Components Separately

1. **Test convertor only**: Run conversion in isolation to check output format

2. **Test generator only**: Use pre-converted standard format files with the tool:
   ```bash
   ./network-config-generator --input_json standard_format.json --output_folder test/
   ```

### Validate Generated Configs

1. **Check config syntax** (vendor-specific):
   ```bash
   # Cisco NX-OS
   nxos-syntax-check generated_config.cfg
   
   # Dell OS10
   os10-syntax-check generated_config.cfg
   ```

2. **Compare with expected output**:
   ```bash
   diff expected_config.cfg generated_config.cfg
   ```

## 📞 Getting Help

### Before Reporting Issues

1. **Check this troubleshooting guide**
2. **Verify your input format** matches expected structure
3. **Test with sample data** from the repository
4. **Check template syntax** if using custom templates

### When Reporting Issues

Include the following information:

1. **Command used**:
   ```bash
   ./network-config-generator --input_json input.json --output_folder output/
   ```

2. **Error message** (full text)

3. **Input file structure** (sanitized):
   ```json
   {
     "switch": {
       "make": "cisco",
       "firmware": "nxos"
     }
   }
   ```

4. **Environment details**:
   - OS: Windows/Linux
   - Python version (if using Python)
   - Tool version

5. **Expected vs actual behavior**

### Useful Debug Commands

```bash
# Check tool version and help
./network-config-generator --help

# Test with minimal input
echo '{"switch":{"make":"cisco","firmware":"nxos","hostname":"test"},"vlans":[],"interfaces":[]}' > minimal.json
./network-config-generator --input_json minimal.json --output_folder debug/

# Validate JSON syntax
./network-config-generator --input_json your_input.json --output_folder debug/ --dry-run

# Check template structure
find input/jinja2_templates -name "*.j2" -type f

# Test with minimal input
echo '{"switch":{"make":"cisco","firmware":"nxos","hostname":"test"},"vlans":[],"interfaces":[]}' > minimal.json
./network-config-generator --input_json minimal.json --output_folder debug/
```

### Community Resources

- **GitHub Issues**: Report bugs and feature requests
- **GitHub Discussions**: Ask questions and share solutions
- **Documentation**: Check the `docs/` folder for detailed guides
- **Examples**: Look at `tests/test_cases/` for working examples

Remember: Most issues are related to input format, template paths, or JSON syntax. Double-check these areas first!
