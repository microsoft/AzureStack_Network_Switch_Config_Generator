from pathlib import Path
import sys
import json
import pytest
import warnings
from collections import defaultdict

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

# === Test result tracking ===
convertor_results = defaultdict(lambda: {"passed": 0, "skipped": 0, "failed": 0, "errors": []})

# === Helper function for better diff reporting ===
def find_json_differences(expected, actual, path="", max_diff=10):
    """Find specific differences between two JSON structures."""
    differences = []
    
    if type(expected) != type(actual):
        differences.append(f"Type mismatch at '{path}': expected {type(expected).__name__}, got {type(actual).__name__}")
        return differences
    
    if isinstance(expected, dict):
        all_keys = set(expected.keys()) | set(actual.keys())
        for key in sorted(all_keys):
            current_path = f"{path}.{key}" if path else key
            if key not in expected:
                differences.append(f"Unexpected key at '{current_path}': {actual[key]}")
            elif key not in actual:
                differences.append(f"Missing key at '{current_path}'")
            else:
                differences.extend(find_json_differences(expected[key], actual[key], current_path))
            
            if len(differences) >= max_diff:
                break
    
    elif isinstance(expected, list):
        if len(expected) != len(actual):
            differences.append(f"List length mismatch at '{path}': expected {len(expected)}, got {len(actual)}")
        
        for i in range(min(len(expected), len(actual))):
            current_path = f"{path}[{i}]" if path else f"[{i}]"
            differences.extend(find_json_differences(expected[i], actual[i], current_path))
            
            if len(differences) >= max_diff:
                break
    
    else:
        if expected != actual:
            differences.append(f"Value mismatch at '{path}': expected '{expected}', got '{actual}'")
    
    return differences[:max_diff]


def validate_json_structure(data, required_fields=None):
    """Validate that JSON has expected structure."""
    errors = []
    
    if not isinstance(data, dict):
        errors.append(f"Expected dict, got {type(data).__name__}")
        return errors
    
    if required_fields:
        for field in required_fields:
            if field not in data:
                errors.append(f"Missing required field: '{field}'")
    
    return errors


def validate_lab_format(data):
    """Validate that input JSON is in lab format."""
    required_keys = ["Version", "Description", "InputData"]
    return validate_json_structure(data, required_keys)


def validate_standard_format(data):
    """Validate that converted JSON is in standard format."""
    required_keys = ["switch"]
    missing = []
    
    for key in required_keys:
        if key not in data:
            missing.append(f"Missing required key: '{key}'")
    
    # Validate switch object has essential properties
    if "switch" in data and isinstance(data["switch"], dict):
        switch_required = ["make", "firmware"]
        for key in switch_required:
            if key not in data["switch"]:
                missing.append(f"Missing in switch: '{key}'")
    
    return missing


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

    # Validate input format
    format_errors = validate_lab_format(input_data)
    if format_errors:
        error_msg = f"Invalid lab format in {folder_name}:\n" + "\n".join(f"  • {e}" for e in format_errors)
        pytest.fail(error_msg, pytrace=False)

    # Run convertor (writes JSON files into output_dir)
    # Suppress print output during testing
    import io, sys
    old_stdout = sys.stdout
    sys.stdout = io.StringIO()
    try:
        convert_lab_switches(input_data, output_dir)
    except Exception as e:
        sys.stdout = old_stdout
        convertor_results[folder_name]["failed"] += 1
        convertor_results[folder_name]["errors"].append(str(e))
        pytest.fail(f"Convertor failed: {e}", pytrace=False)
    finally:
        sys.stdout = old_stdout

    # Compare each generated output file with its corresponding expected file
    generated_files = list(output_dir.glob("*.json"))
    
    if not generated_files:
        convertor_results[folder_name]["skipped"] += 1
        pytest.skip(f"No output files generated for {folder_name}")
        return

    compared_files = 0
    missing_files = []
    format_issues = []
    
    for generated_file in sorted(generated_files):
        expected_file = expected_dir / generated_file.name

        if not expected_file.exists():
            missing_files.append(generated_file.name)
            continue

        with open(expected_file, "r", encoding="utf-8") as f:
            expected_data = json.load(f)
        with open(generated_file, "r", encoding="utf-8") as f:
            generated_data = json.load(f)

        # Validate generated format
        validation_errors = validate_standard_format(generated_data)
        if validation_errors:
            format_issues.extend([f"{generated_file.name}: {e}" for e in validation_errors])

        # Content comparison
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
            convertor_results[folder_name]["failed"] += 1
            convertor_results[folder_name]["errors"].append(generated_file.name)
            pytest.fail(f"File comparison failed: {generated_file.name}", pytrace=False)
        
        compared_files += 1
        convertor_results[folder_name]["passed"] += 1

    # Report format issues as warnings if content matched
    if format_issues and compared_files > 0:
        for issue in format_issues:
            print(f"⚠️  {issue}")

    # Report results
    if compared_files == 0:
        if missing_files:
            skip_msg = f"All expected files missing: {', '.join(missing_files)}"
            print(f"⏭️  {folder_name}: SKIPPED - {skip_msg}")
            convertor_results[folder_name]["skipped"] += 1
            pytest.skip(skip_msg)
        else:
            skip_msg = "No files to compare"
            print(f"⏭️  {folder_name}: SKIPPED - {skip_msg}")
            convertor_results[folder_name]["skipped"] += 1
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


# === Additional validation tests ===
@pytest.mark.parametrize("input_case", ALL_INPUT_CASES, ids=lambda val: f"input_format_{val[0]}")
def test_input_format_validation(input_case):
    """Verify that input files conform to expected lab format."""
    folder_name, input_file = input_case
    
    input_data = load_input_json(input_file)
    assert input_data is not None, f"Failed to load {input_file}"
    
    # Validate lab format structure
    errors = validate_lab_format(input_data)
    assert not errors, f"Format validation failed for {folder_name}:\n" + "\n".join(f"  • {e}" for e in errors)


@pytest.mark.parametrize("input_case", ALL_INPUT_CASES, ids=lambda val: f"output_format_{val[0]}")
def test_output_format_validation(input_case):
    """Verify that generated output conforms to standard format."""
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name
    output_dir = folder_path / "generated_outputs"
    
    if not output_dir.exists() or not list(output_dir.glob("*.json")):
        pytest.skip(f"No generated outputs for {folder_name}")
    
    # Validate each generated file
    for output_file in output_dir.glob("*.json"):
        with open(output_file, "r", encoding="utf-8") as f:
            output_data = json.load(f)
        
        errors = validate_standard_format(output_data)
        assert not errors, f"Format validation failed for {output_file.name}:\n" + "\n".join(f"  • {e}" for e in errors)


def pytest_sessionfinish(session, exitstatus):
    """Print final summary of convertor tests."""
    if convertor_results:
        print(f"\n📊 Convertor Test Summary:")
        total_passed = sum(r["passed"] for r in convertor_results.values())
        total_skipped = sum(r["skipped"] for r in convertor_results.values())
        total_failed = sum(r["failed"] for r in convertor_results.values())
        total_cases = len(convertor_results)
        
        print(f"📈 {total_cases} convertor cases: {total_passed} passed, {total_skipped} skipped, {total_failed} failed")
        
        if total_failed > 0:
            print("\n❌ Failed cases:")
            for case_name, results in sorted(convertor_results.items()):
                if results["failed"] > 0:
                    print(f"   • {case_name}: {results['failed']} failure(s)")
                    for error in results["errors"][:3]:
                        print(f"     - {error}")
