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
from convertors import convert_lab_switches
from loader import load_input_json

# === Helper function for better diff reporting ===
def find_json_differences(expected, actual, path=""):
    """Find specific differences between two JSON structures."""
    differences = []
    
    if type(expected) != type(actual):
        differences.append(f"Type mismatch at '{path}': expected {type(expected).__name__}, got {type(actual).__name__}")
        return differences
    
    if isinstance(expected, dict):
        all_keys = set(expected.keys()) | set(actual.keys())
        for key in all_keys:
            current_path = f"{path}.{key}" if path else key
            if key not in expected:
                differences.append(f"Unexpected key at '{current_path}': {actual[key]}")
            elif key not in actual:
                differences.append(f"Missing key at '{current_path}'")
            else:
                differences.extend(find_json_differences(expected[key], actual[key], current_path))
    
    elif isinstance(expected, list):
        if len(expected) != len(actual):
            differences.append(f"List length mismatch at '{path}': expected {len(expected)}, got {len(actual)}")
        
        for i in range(min(len(expected), len(actual))):
            current_path = f"{path}[{i}]" if path else f"[{i}]"
            differences.extend(find_json_differences(expected[i], actual[i], current_path))
    
    else:
        if expected != actual:
            differences.append(f"Value mismatch at '{path}': expected '{expected}', got '{actual}'")
    
    return differences

# === Step 1: Find test folders with *_input.json ===
def find_input_cases():
    input_cases = []
    for folder in TEST_CASES_ROOT.iterdir():
        if not folder.is_dir():
            continue
        if not folder.name.startswith("convert_"):
            continue

        for input_file in folder.glob("*_input.json"):
            input_cases.append((folder.name, input_file))

    return input_cases

# === Step 2: Core convertor test logic ===
def run_convert_and_compare(folder_name, input_file):
    folder_path = TEST_CASES_ROOT / folder_name
    expected_dir = folder_path / "expected_outputs"
    output_dir = folder_path / "generated_outputs"
    output_dir.mkdir(exist_ok=True)

    # Load input JSON
    input_data = load_input_json(input_file)
    assert input_data is not None, "Failed to load input JSON"

    # Run convertor (writes JSON files into output_dir)
    # Suppress print output during testing
    import io, sys
    old_stdout = sys.stdout
    sys.stdout = io.StringIO()
    try:
        convert_lab_switches(input_data, output_dir)
    finally:
        sys.stdout = old_stdout

    # Compare each generated output file with its corresponding expected file
    generated_files = list(output_dir.glob("*.json"))
    
    if not generated_files:
        pytest.skip(f"No output files generated for {folder_name}")
        return

    compared_files = 0
    missing_files = []
    
    for generated_file in generated_files:
        expected_file = expected_dir / generated_file.name

        if not expected_file.exists():
            missing_files.append(generated_file.name)
            continue

        with open(expected_file, "r", encoding="utf-8") as f:
            expected_data = json.load(f)
        with open(generated_file, "r", encoding="utf-8") as f:
            generated_data = json.load(f)

        # Better comparison with detailed error reporting
        if expected_data != generated_data:
            # Find and report specific differences
            differences = find_json_differences(expected_data, generated_data)
            error_msg = f"❌ {generated_file.name} - Differences found:"
            for diff in differences[:5]:  # Show first 5 differences
                error_msg += f"\n  • {diff}"
            if len(differences) > 5:
                error_msg += f"\n  ... and {len(differences) - 5} more differences"
            
            # Print the error for immediate visibility, then fail
            print(error_msg)
            pytest.fail(f"File comparison failed: {generated_file.name}", pytrace=False)
        
        compared_files += 1

    # Report results
    if compared_files == 0:
        if missing_files:
            skip_msg = f"All expected files missing: {', '.join(missing_files)}"
            print(f"⏭️  {folder_name}: SKIPPED - {skip_msg}")
            pytest.skip(skip_msg)
        else:
            skip_msg = "No files to compare"
            print(f"⏭️  {folder_name}: SKIPPED - {skip_msg}")
            pytest.skip(skip_msg)
    else:
        print(f"✅ {folder_name}: {compared_files} file(s) verified")
        if missing_files:
            print(f"⚠️  {folder_name}: {len(missing_files)} expected file(s) missing: {', '.join(missing_files)}")

# === Step 3: Parametrize test for pytest ===
ALL_INPUT_CASES = find_input_cases()

@pytest.mark.parametrize("input_case", ALL_INPUT_CASES, ids=lambda val: val[0])
def test_convert_switch_input_json(input_case):
    folder_name, input_file = input_case
    run_convert_and_compare(folder_name, input_file)