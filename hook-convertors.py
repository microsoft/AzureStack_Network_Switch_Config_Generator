# PyInstaller hook for convertors package
# This ensures all convertor modules are included in the bundle

from PyInstaller.utils.hooks import collect_submodules, collect_data_files

# Collect all submodules from the convertors package
hiddenimports = collect_submodules('convertors')

# Also explicitly add the BMC converter module
if 'convertors.convertors_bmc_switch_json' not in hiddenimports:
    hiddenimports.append('convertors.convertors_bmc_switch_json')

print(f"[HOOK] convertors hiddenimports: {hiddenimports}")
