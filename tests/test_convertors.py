from pathlib import Path
import sys
import json
import pytest
import warnings

# === Path setup ===
ROOT_DIR = Path(__file__).resolve().parent.parent
SRC_PATH = ROOT_DIR / "src"
TEMPLATE_ROOT = ROOT_DIR / "input" / "templates"
TEST_CASES_ROOT = ROOT_DIR / "tests" / "test_cases"

# Add src to sys.path
if str(SRC_PATH) not in sys.path:
    sys.path.insert(0, str(SRC_PATH))

# === Imports from your project ===
from convertors.convertors_lab_switch_json import convert_switch_input_json
from loader import load_input_json

# === Step 1: Find test folders with *_input.json ===
def find_input_cases():
    input_cases = []
    for folder in TEST_CASES_ROOT.iterdir():
        if not folder.is_dir():
            continue
        if not folder.name.startswith("convert_"):
            print(f"[SKIP] Ignoring non-convert folder: {folder.name}")
            continue

        for input_file in folder.glob("*_input.json"):
            input_cases.append((folder.name, input_file))

    print(f"[INFO] Found {len(input_cases)} test input case(s).")
    return input_cases

# === Step 2: Core convertor test logic ===
def run_convert_and_compare(folder_name, input_file):
    folder_path = TEST_CASES_ROOT / folder_name
    expected_dir = folder_path / "expected_outputs"
    output_dir = folder_path / "generated_outputs"
    output_dir.mkdir(exist_ok=True)

    # Load input JSON
    input_data = load_input_json(input_file)
    assert input_data is not None, "❌ Failed to load input JSON"

    # Run convertor (writes JSON files into output_dir)
    convert_switch_input_json(input_data, output_dir)

    # Compare each output (tor1.json, tor2.json, bmc.json)
    for name in ["tor1", "tor2", "bmc"]:
        expected_file = expected_dir / f"{name}.json"
        generated_file = output_dir / f"{name}.json"

        if not expected_file.exists():
            warnings.warn(f"[WARN] Expected file missing, skipping: {expected_file}", UserWarning)
            pytest.skip(f"[SKIP] No expected output for {name}")
            continue

        assert generated_file.exists(), f"❌ Missing generated file: {generated_file}"

        with open(expected_file, "r", encoding="utf-8") as f:
            expected_data = json.load(f)
        with open(generated_file, "r", encoding="utf-8") as f:
            generated_data = json.load(f)

        assert expected_data == generated_data, f"❌ Mismatch in {name}.json"

    print(f"✅ Passed: {folder_name}")

# === Step 3: Parametrize test for pytest ===
ALL_INPUT_CASES = find_input_cases()

@pytest.mark.parametrize("input_case", ALL_INPUT_CASES, ids=lambda val: val[0])
def test_convert_switch_input_json(input_case):
    folder_name, input_file = input_case
    run_convert_and_compare(folder_name, input_file)