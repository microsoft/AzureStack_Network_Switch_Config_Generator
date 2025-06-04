from pathlib import Path
from loader import load_input_json, load_template
import os

def generate_config(input_path, template_path, output_path):
    if not Path(input_path).exists():
        print(f"[ERROR] File not found: {input_path}")
        raise FileNotFoundError(f"Input JSON not found: {input_path}")

    data = load_input_json(input_path)
    if data is None:
        raise ValueError(f"Input JSON was empty or failed to parse: {input_path}")

    template_dir = str(Path(template_path).parent)
    template_file = Path(template_path).name
    template = load_template(template_dir, template_file)

    rendered = template.render(data)

    os.makedirs(Path(output_path).parent, exist_ok=True)
    with open(output_path, 'w') as f:
        f.write(rendered)

    print(f"Config generated: {output_path}")
