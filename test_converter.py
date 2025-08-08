#!/usr/bin/env python3

import sys
import json
from pathlib import Path

# Add src to path
sys.path.insert(0, str(Path("src").resolve()))

from convertors.convertors_lab_switch_json import convert_switch_input_json

# Load the Dell lab input
with open("tests/test_cases/convert_lab_switch_input_json_dell_os10/lab_dell_os10_switch_input.json") as f:
    input_data = json.load(f)

# Convert and generate outputs
convert_switch_input_json(input_data, "test_output")

print("Conversion complete. Check test_output/ for generated files.")
