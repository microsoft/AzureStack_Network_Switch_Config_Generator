from pathlib import Path
import sys
import pytest
import warnings
from collections import defaultdict
import re

# === Test result tracking ===
case_results = defaultdict(lambda: {"passed": 0, "skipped": 0, "failed": 0, "files": []})
case_summary_shown = set()

# === Path setup ===
ROOT_DIR = Path(__file__).resolve().parent.parent
SRC_PATH = ROOT_DIR / "src"
TEMPLATE_ROOT = ROOT_DIR / "input" / "jinja2_templates"
TEST_CASES_ROOT = ROOT_DIR / "tests" / "test_cases"

# Add src to sys.path
if str(SRC_PATH) not in sys.path:
    sys.path.insert(0, str(SRC_PATH))

from generator import generate_config

# === Helper function for better diff reporting ===
def find_text_differences(expected, actual, max_lines=10):
    """Find specific differences between two text files."""
    expected_lines = expected.split('\n')
    actual_lines = actual.split('\n')
    differences = []
    
    # Check line count difference
    if len(expected_lines) != len(actual_lines):
        differences.append(f"Line count mismatch: expected {len(expected_lines)}, got {len(actual_lines)}")
    
    # Compare line by line
    max_compare = min(len(expected_lines), len(actual_lines))
    for i in range(max_compare):
        if expected_lines[i] != actual_lines[i]:
            differences.append(f"Line {i+1}: expected '{expected_lines[i][:50]}...', got '{actual_lines[i][:50]}...'")
            if len(differences) >= max_lines:
                break
    
    return differences


def validate_config_syntax(content, file_type="cfg"):
    """Validate configuration file syntax."""
    errors = []
    lines = content.split('\n')
    
    if file_type == "cfg":
        # Basic validation for Cisco/Dell config files
        bracket_stack = []
        for line_num, line in enumerate(lines, 1):
            stripped = line.strip()
            
            # Skip comments and empty lines
            if not stripped or stripped.startswith('!') or stripped.startswith('#'):
                continue
            
            # Check for unbalanced brackets/braces
            for char in stripped:
                if char == '{':
                    bracket_stack.append('{')
                elif char == '}':
                    if not bracket_stack or bracket_stack[-1] != '{':
                        errors.append(f"Line {line_num}: Unbalanced closing brace")
                    else:
                        bracket_stack.pop()
        
        if bracket_stack:
            errors.append(f"Unclosed braces: {len(bracket_stack)} opening brace(s) without closing")
        
        # Check for common configuration errors
        if 'interface' in content.lower():
            # Basic interface validation
            interface_lines = [l for l in lines if 'interface' in l.lower()]
            if interface_lines and len(interface_lines) < 2:
                warnings.warn(f"Config may be incomplete: found {len(interface_lines)} interface definition(s)", UserWarning)
    
    return errors


def validate_generated_structure(output_folder):
    """Validate that required output files were generated."""
    output_folder = Path(output_folder)
    expected_file_types = [".cfg"]
    generated_files = []
    
    for file_type in expected_file_types:
        generated_files.extend(output_folder.glob(f"*{file_type}"))
    
    return len(generated_files) > 0, len(generated_files)



# === Step 1: Find test folders with input ===
def find_input_cases():
    input_cases = []

    for folder in TEST_CASES_ROOT.iterdir():
        if not folder.is_dir():
            continue

        # Only process test case folders that start with 'std_'
        if not folder.name.startswith("std_"):
            continue

        for input_file in folder.glob("*_input.json"):
            input_cases.append((folder.name, input_file))

    return input_cases


# === Step 4: Render each template (compact single-line logging) ===
def generate_all_configs(input_case):
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name
    output_folder = folder_path

    # Suppress print output during generation
    import io, sys
    old_stdout = sys.stdout
    sys.stdout = io.StringIO()
    try:
        generate_config(
            input_std_json=str(input_file),
            template_folder=str(TEMPLATE_ROOT),
            output_folder=str(output_folder)
        )
    except Exception as e:
        sys.stdout = old_stdout
        # Only log critical errors, suppress template path warnings during test discovery
        if "Template path not found" not in str(e):
            warnings.warn(f"Failed to generate configs for {folder_name}: {e}")
        return []
    finally:
        sys.stdout = old_stdout

    # Gather all output files and map to their expected counterparts
    return [
        (
            folder_name,
            output_file.stem.replace("generated_", ""),
            str(output_file),
            str(folder_path / f"expected_{output_file.stem.replace('generated_', '')}.cfg")
        )
        for output_file in sorted(output_folder.glob("generated_*.cfg"))
    ]

def discover_test_cases():
    all_cases = []
    input_folders = find_input_cases()

    for case in input_folders:
        case_tests = generate_all_configs(case)
        all_cases.extend(case_tests)

    return all_cases

# === Run pytest parametrize ===
ALL_TEST_CASES = discover_test_cases()

@pytest.mark.parametrize(
    "folder_name,case_name,generated_path,expected_path",
    ALL_TEST_CASES,
    ids=lambda val: f"{val[0]}/{val[1]}" if isinstance(val, tuple) and len(val) >= 2 else str(val)
)
def test_generated_config_output(folder_name, case_name, generated_path, expected_path):
    expected_file = Path(expected_path)
    test_name = f"{folder_name}/{case_name}"
    
    if not expected_file.exists():
        case_results[folder_name]["skipped"] += 1
        case_results[folder_name]["files"].append(f"⏭️  {case_name}")
        
        # Show case summary when the last file for this case is processed
        show_case_summary_if_complete(folder_name)
        pytest.skip("Expected file missing")
        return

    with open(generated_path) as gen, open(expected_file) as exp:
        generated = gen.read().strip()
        expected = exp.read().strip()

    # Validate syntax before comparison
    syntax_errors = validate_config_syntax(generated, "cfg")
    if syntax_errors:
        case_results[folder_name]["failed"] += 1
        case_results[folder_name]["files"].append(f"❌ {case_name} (syntax)")
        error_msg = f"❌ {test_name} - Syntax errors:\n"
        for err in syntax_errors[:3]:
            error_msg += f"  • {err}\n"
        print(error_msg)
        show_case_summary_if_complete(folder_name)
        pytest.fail(f"Config syntax validation failed", pytrace=False)

    # Content comparison
    if expected != generated:
        case_results[folder_name]["failed"] += 1
        case_results[folder_name]["files"].append(f"❌ {case_name}")
        
        # Find and report specific differences
        differences = find_text_differences(expected, generated)
        error_msg = f"❌ {test_name} - Differences found:"
        for diff in differences[:3]:  # Show first 3 differences for config files
            error_msg += f"\n  • {diff}"
        if len(differences) > 3:
            error_msg += f"\n  ... and {len(differences) - 3} more differences"
        
        # Print the error for immediate visibility, then fail
        print(error_msg)
        show_case_summary_if_complete(folder_name)
        pytest.fail(f"Config comparison failed", pytrace=False)

    case_results[folder_name]["passed"] += 1
    case_results[folder_name]["files"].append(f"✅ {case_name}")
    show_case_summary_if_complete(folder_name)

def show_case_summary_if_complete(folder_name):
    """Show summary for a test case if all its files have been processed"""
    if folder_name in case_summary_shown:
        return
        
    results = case_results[folder_name]
    total_files = results["passed"] + results["skipped"] + results["failed"]
    expected_total = get_expected_file_count(folder_name)
    
    if total_files >= expected_total:
        case_summary_shown.add(folder_name)
        status_icon = "✅" if results["failed"] == 0 else "❌"
        print(f"\n{status_icon} {folder_name}: {results['passed']} passed, {results['skipped']} skipped, {results['failed']} failed")

def get_expected_file_count(folder_name):
    """Get the expected number of test files for a case"""
    # Count the number of generated files to determine expected count
    folder_path = TEST_CASES_ROOT / folder_name
    return len(list(folder_path.glob("generated_*.cfg")))


# === Additional validation tests ===
@pytest.mark.parametrize("input_case", find_input_cases(), ids=lambda val: f"config_syntax_{val[0]}")
def test_config_syntax_validation(input_case):
    """Verify that generated configuration files have valid syntax."""
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name
    
    # Generate configs
    config_files = list(folder_path.glob("generated_*.cfg"))
    
    if not config_files:
        pytest.skip(f"No generated configs for {folder_name}")
    
    # Validate syntax of each config
    for config_file in config_files:
        with open(config_file, "r", encoding="utf-8") as f:
            content = f.read()
        
        errors = validate_config_syntax(content, "cfg")
        assert not errors, f"Syntax errors in {config_file.name}:\n" + "\n".join(f"  • {e}" for e in errors)


@pytest.mark.parametrize("input_case", find_input_cases(), ids=lambda val: f"output_generation_{val[0]}")
def test_output_file_generation(input_case):
    """Verify that output files are generated correctly."""
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name
    
    has_output, file_count = validate_generated_structure(folder_path)
    assert has_output, f"No output files generated for {folder_name}"
    assert file_count > 0, f"Expected at least 1 output file, got {file_count}"


def pytest_sessionfinish(session, exitstatus):
    """Print final summary"""
    if case_results:
        print(f"\n📊 Generator Test Summary:")
        total_passed = sum(r["passed"] for r in case_results.values())
        total_skipped = sum(r["skipped"] for r in case_results.values())
        total_failed = sum(r["failed"] for r in case_results.values())
        total_cases = len(case_results)
        print(f"📈 {total_cases} test cases: {total_passed} passed, {total_skipped} skipped, {total_failed} failed")
        
        if total_failed > 0:
            print("\n❌ Failed tests:")
            for folder_name, results in sorted(case_results.items()):
                if results["failed"] > 0:
                    print(f"   • {folder_name}: {results['failed']} failure(s)")

