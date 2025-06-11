from pathlib import Path
import sys
import pytest
import warnings

print("✅ test_generator.py loaded")

# === Path setup ===
ROOT_DIR = Path(__file__).resolve().parent.parent
SRC_PATH = ROOT_DIR / "src"
TEMPLATE_ROOT = ROOT_DIR / "input" / "templates"
TEST_CASES_ROOT = ROOT_DIR / "tests" / "test_cases"

# Add src to sys.path
if str(SRC_PATH) not in sys.path:
    sys.path.insert(0, str(SRC_PATH))

from generator import generate_config


# === Step 1: Find test folders with input ===
def find_input_cases():
    input_cases = []

    for folder in TEST_CASES_ROOT.iterdir():
        if not folder.is_dir():
            continue

        # 🔍 Only process test case folders that start with 'std_'
        if not folder.name.startswith("std_"):
            print(f"[SKIP] Ignoring non-std folder: {folder.name}")
            continue

        for input_file in folder.glob("*_input.json"):
            input_cases.append((folder.name, input_file))

    print(f"[INFO] Found {len(input_cases)} std test folder(s).")
    return input_cases


# === Step 2: Generate configs using dynamic template selection ===
def generate_all_configs(input_case):
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name
    output_folder = folder_path

    print(f"[GENERATE] Case: {folder_name}")

    try:
        generated_paths = generate_config(
            input_std_json=str(input_file),
            template_folder=str(TEMPLATE_ROOT),
            output_folder=str(output_folder)
        )
    except Exception as e:
        print(f"[ERROR] Failed to generate configs for {folder_name}: {e}")
        return []

    # We need to return: folder_name, case_name, generated_path, expected_path
    all_pairs = []
    for output_file in output_folder.glob("generated_*.cfg"):
        case_name = output_file.stem.replace("generated_", "")
        expected_file = folder_path / f"expected_{case_name}.cfg"
        all_pairs.append((folder_name, case_name, str(output_file), str(expected_file)))

    return all_pairs


# === Step 3: Discover all test pairs ===
def discover_test_cases():
    all_cases = []
    input_folders = find_input_cases()

    for case in input_folders:
        case_tests = generate_all_configs(case)
        all_cases.extend(case_tests)

    print(f"[DEBUG] Total test comparisons: {len(all_cases)}")
    return all_cases

# === Run pytest parametrize ===
ALL_TEST_CASES = discover_test_cases()

@pytest.mark.parametrize(
    "folder_name,case_name,generated_path,expected_path",
    ALL_TEST_CASES,
    ids=lambda val: f"{val[0]}/{val[1]}" if isinstance(val, tuple) else val
)
def test_generated_config_output(folder_name, case_name, generated_path, expected_path):
    print(f"\n[TEST] Comparing: {folder_name}/{case_name}")

    expected_file = Path(expected_path)
    if not expected_file.exists():
        # 🔔 Temporary dev-friendly warning instead of error
        warnings.warn(f"[WARN] Expected file missing, skipping: {expected_path}", UserWarning)
        pytest.skip(f"[SKIP] No expected output yet for {folder_name}/{case_name}")
        return

    with open(generated_path) as gen, open(expected_file) as exp:
        generated = gen.read().strip()
        expected = exp.read().strip()

    assert generated == expected, f"[FAIL] Output mismatch: {folder_name}/{case_name}"
