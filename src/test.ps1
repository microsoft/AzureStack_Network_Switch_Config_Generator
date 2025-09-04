python -m pip install -r requirements.txt

pyinstaller --onefile src/main.py `
  --name network_config_generator `
  --add-data "input/templates;input/templates"


pytest -s .\tests\test_generator.py -v
pytest -s .\tests\test_convertors.py -v