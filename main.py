# main.py
import argparse
import json
import os

from src.loader import load_json, pretty_print_json

# Path to your JSON file
json_path = "C:\\WORK\\REPO\\Github\\AzureStack_Network_Switch_Config_Generator\\input\\standard_input.json"

# Load the JSON
data = load_json(json_path, verbose=True)

# If loaded successfully, overwrite it with pretty format
if data:
    pretty_print_json(data, json_path)
