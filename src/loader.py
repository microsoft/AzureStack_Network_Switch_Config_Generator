# loader.py
import json
import os

def load_json(filepath, verbose=False):
    """
    Load and parse a JSON file safely.

    Args:
        filepath (str): Path to the JSON file.
        verbose (bool): Print extra info for debugging.

    Returns:
        dict or list: Parsed JSON data, or None if there's an error.
    """
    if not os.path.exists(filepath):
        print(f"[ERROR] File not found: {filepath}")
        return None

    try:
        with open(filepath, "r", encoding="utf-8") as file:
            data = json.load(file)
            if verbose:
                print(f"[✓] Loaded JSON from: {filepath}")
                if isinstance(data, dict):
                    print(f"     Top-level keys: {list(data.keys())}")
            return data
    except json.JSONDecodeError as e:
        print(f"[ERROR] Failed to parse JSON ({filepath}): {e}")
        return None
    except Exception as e:
        print(f"[ERROR] Unexpected error while loading {filepath}: {e}")
        return None


def pretty_print_json(data, output_path):
    """
    Pretty-print JSON data to a file.

    Args:
        data (dict or list): The JSON object to save.
        output_path (str): Path to the output file.
    """
    try:
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        with open(output_path, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)
        print(f"[✓] Pretty-printed JSON saved to: {output_path}")
    except Exception as e:
        print(f"[ERROR] Failed to write pretty JSON: {e}")
