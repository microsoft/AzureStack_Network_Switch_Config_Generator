from pathlib import Path
import sys
import json
import pytest
import warnings

ROOT_DIR = Path(__file__).resolve().parent.parent
SRC_PATH = ROOT_DIR / "src"

# Add to Python path if not already
if str(SRC_PATH) not in sys.path:
    sys.path.insert(0, str(SRC_PATH))

from convertors import convert_switch_input_json
from loader import load_input_json

# Global test case base path
BASE_PATH = ROOT_DIR / "tests" / "test_cases"

def test_convert_switch_input_json():
    folder = "convert_switch_input_json"
    input_path = BASE_PATH / folder / "switch_input.json"
    expected_dir = BASE_PATH / folder / "expected_outputs"
    output_dir = BASE_PATH / folder / "generated_outputs"

    input_data = load_input_json(input_path)
    assert input_data is not None, "❌ Failed to load input JSON"

    # 🔄 Run conversion (writes files to output_dir)
    convert_switch_input_json(input_data, output_dir)

    # 🧪 Check each expected output file
    for name in ["tor1", "tor2", "bmc"]:
        expected_file = expected_dir / f"{name}.json"
        generated_file = output_dir / f"{name}.json"

        if not expected_file.exists():
            warnings.warn(f"[WARN] Expected file missing, skipping: {expected_file}", UserWarning)
            pytest.skip(f"[SKIP] No expected output for {name}")
            continue  # Not strictly needed, since skip halts the test

        assert generated_file.exists(), f"❌ Missing generated file: {generated_file}"

        with open(expected_file, "r", encoding="utf-8") as f:
            expected_data = json.load(f)
        with open(generated_file, "r", encoding="utf-8") as f:
            generated_data = json.load(f)

        assert expected_data == generated_data, f"❌ Mismatch in {name}.json"

    print("✅ All available outputs match expected JSON files.")