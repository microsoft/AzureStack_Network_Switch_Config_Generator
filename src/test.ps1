pyinstaller --onefile src/main.py `
  --name network_config_generator `
  --add-data "input/templates;input/templates"


.\network_config_generator.exe --input_std_json std_nxos_hyperconverged_input.json --output_folder .