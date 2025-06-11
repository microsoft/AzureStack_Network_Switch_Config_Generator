import argparse
from pathlib import Path
import sys
from generator import generate_config
from loader import get_real_path  # Only used for PyInstaller-packed assets

def main():
    parser = argparse.ArgumentParser(description="Run config generator with auto-detected templates.")

    parser.add_argument("--input_std_json", required=True,
                        help="Path to input standard JSON file (external, not bundled)")

    parser.add_argument("--template_folder", default="input/templates",
                        help="Root folder containing vendor templates (default: input/templates)")

    parser.add_argument("--output_folder", default=".",
        help="Directory to save generated config files (default: current directory)"
    )


    args = parser.parse_args()

    # ✅ Only use .resolve() for user-provided files
    input_std_json_path = Path(args.input_std_json).resolve()
    output_folder_path = Path(args.output_folder).resolve()
    template_folder_arg = Path(args.template_folder)

    # Only use get_real_path if user did NOT override default
    if args.template_folder == parser.get_default('template_folder'):
        template_folder = get_real_path(template_folder_arg)
    else:
        template_folder = template_folder_arg.resolve()


    print(f"🧾 Input JSON File:     {input_std_json_path}")
    print(f"🧩 Template Folder:     {template_folder}")
    print(f"📁 Output Directory:    {output_folder_path}")

    # === Validation ===
    if not input_std_json_path.exists():
        print(f"[ERROR] Input file not found: {input_std_json_path}")
        sys.exit(1)

    if not template_folder.exists():
        print(f"[ERROR] Template folder not found: {template_folder}")
        sys.exit(1)

    output_folder_path.mkdir(parents=True, exist_ok=True)

    # === Run generation ===
    try:
        generate_config(
            input_std_json=str(input_std_json_path),
            template_folder=str(template_folder),
            output_folder=str(output_folder_path)
        )
        print("✅ Configs generated successfully!")
    except Exception as e:
        print(f"❌ Failed to generate configs: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
