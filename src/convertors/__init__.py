"""
Convertors Package - Static Registry Pattern

This package contains convertor modules that transform various input formats
into the standardized JSON format used by the network config generator.

All convertors are statically imported here so PyInstaller can detect them
without needing hooks or hidden imports.

Available convertors:
- convertors_lab_switch_json: Converts lab-style input JSON to standard format
- convertors_bmc_switch_json: Converts BMC switch definitions to standard format
"""

# Static imports - PyInstaller can detect these automatically
from .convertors_lab_switch_json import convert_switch_input_json as convert_lab_switches
from .convertors_bmc_switch_json import convert_bmc_switches

# Registry mapping convertor names to functions
# This allows main.py to use registry lookup instead of dynamic imports
CONVERTORS = {
    # Primary convertor names (for backward compatibility)
    'convertors.convertors_lab_switch_json': convert_lab_switches,
    'convertors.convertors_bmc_switch_json': convert_bmc_switches,
    
    # Short aliases for convenience
    'lab': convert_lab_switches,
    'bmc': convert_bmc_switches,
}

# Export public API
__all__ = [
    'CONVERTORS',
    'convert_lab_switches',
    'convert_bmc_switches',
]
