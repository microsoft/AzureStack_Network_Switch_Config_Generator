from pathlib import Path
from loader import load_input_json, load_template
import os

def generate_config(input_std_json, template_folder, output_folder):
    input_std_json_path = Path(input_std_json).resolve()
    template_folder_path = Path(template_folder)  # already resolved by main.py
    output_folder_path = Path(output_folder).resolve()

    # ✅ Step 1: Validate input JSON
    if not input_std_json_path.exists():
        raise FileNotFoundError(f"[ERROR] Input JSON not found: {input_std_json_path}")

    data = load_input_json(str(input_std_json_path))
    if data is None:
        raise ValueError(f"[ERROR] Input JSON was empty or failed to parse: {input_std_json_path}")

    # ✅ Step 2: Extract metadata
    try:
        make = data["switch"]["make"].lower()
        firmware = data["switch"]["firmware"].lower()
        version = data["switch"].get("version", "").lower()  # Optional, not currently used
    except KeyError as e:
        raise ValueError(f"[ERROR] Missing expected switch metadata: {e}")

    # ✅ Step 3: Resolve template subfolder
    template_dir = template_folder_path / make / firmware
    if not template_dir.exists():
        raise FileNotFoundError(f"[ERROR] Template path not found: {template_dir}")

    template_files = list(template_dir.glob("*.j2"))
    if not template_files:
        raise FileNotFoundError(f"[WARN] No templates found in: {template_dir}")

    # ✅ Step 4: Render each template
    os.makedirs(output_folder_path, exist_ok=True)

    for template_path in template_files:
        template = load_template(str(template_dir), template_path.name)
        rendered = template.render(data)

        output_file = output_folder_path / f"generated_{template_path.stem}.cfg"
        with open(output_file, "w", encoding="utf-8") as f:
            f.write(rendered)

        print(f"[✓] Generated: {output_file}")
