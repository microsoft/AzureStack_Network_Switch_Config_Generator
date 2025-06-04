from pathlib import Path
import sys
import pytest

print("✅ test_generator.py loaded")

# === Path setup ===
ROOT_DIR = Path(__file__).resolve().parent.parent
SRC_PATH = ROOT_DIR / "src"
TEMPLATE_DIR = ROOT_DIR / "input" / "templates" / "cisco" / "nxos"
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



# === Step 2: Generate configs from templates ===
def generate_all_configs(input_case):
    folder_name, input_file = input_case
    folder_path = TEST_CASES_ROOT / folder_name

    generated_files = []

    if not TEMPLATE_DIR.exists():
        print(f"[WARN] Template directory not found: {TEMPLATE_DIR}")
        return generated_files

    templates = list(TEMPLATE_DIR.glob("*.j2"))
    if not templates:
        print(f"[WARN] No templates found in {TEMPLATE_DIR}")
        return generated_files

    for template_path in templates:
        if not template_path.exists():
            print(f"[WARN] Template not found: {template_path}")
            continue

        case_name = template_path.stem
        output_cfg = folder_path / f"generated_{case_name}.cfg"

        print(f"[GENERATE] {folder_name}/{case_name}")
        generate_config(
            input_path=str(input_file),
            template_path=str(template_path),
            output_path=str(output_cfg)
        )

        generated_files.append((folder_name, case_name, output_cfg, template_path))

    return generated_files


# === Step 3: Pair generated files with expected ===
def discover_test_cases():
    all_cases = []
    input_folders = find_input_cases()

    for case in input_folders:
        generated = generate_all_configs(case)

        for folder_name, case_name, generated_file, template_path in generated:
            expected_file = TEST_CASES_ROOT / folder_name / f"expected_{case_name}.cfg"

            all_cases.append((
                folder_name,
                case_name,
                str(generated_file),
                str(expected_file)
            ))

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

    if not Path(expected_path).exists():
        print(f"[SKIP] No expected file found: {expected_path}")
        pytest.skip(f"Expected output not available: {expected_path}")

    with open(generated_path) as gen, open(expected_path) as exp:
        generated = gen.read().strip()
        expected = exp.read().strip()

    assert generated == expected, f"[FAIL] Output mismatch: {folder_name}/{case_name}"
