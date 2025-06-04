# Quick Start

Common scenario to use the tool

## Prerequisits

### Download the Release Package

### Validate Essential files from the Package

Here are three essential files need to be used:

1. `SwitchConfigGenerator.exe` or `SwitchConfigGenerator` binary file depends on your runing enviroment.
2. `switchLib` folder, which defines the switch frameworks and templates to be used. More detail can be find more detail in [Switch_Library](./Switch_Library.md)
3. `input.json`, this file need to be updated by user. More detail can be find more detail in [User_Input_Json](./User_Input_Json.md)

## Get Start

### Create Your Onn `input.json`

Based on sample input JSON files, please update the values based on your enviroment, and here are the sections to be updated, and please use this doc as instruction [User_Input_Json](./User_Input_Json.md).

### Execute the Tool

Please make sure using admin priviliedge before execute the tool.

```powershell
PS C:\Downloads\switchConfigTool> ls

Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
d-----         2/28/2023   1:07 PM                Input
-a----         2/28/2023   1:25 PM        3724800 SwitchConfigGenerator.exe

PS C:\Downloads\switchConfigTool> .\SwitchConfigGenerator.exe -h
Usage of C:\Downloads\switchConfigTool\SwitchConfigGenerator.exe:
  -inputJsonFile string
        File path of switch deploy input.json (default "../input/cisco_sample_input1.json")
  -outputFolder string
        Folder path of switch configurations (default "../output")
  -password string
        Password for switch configuration
  -switchLib string
        Folder path of switch frameworks and templates (default "../input/switchLib")
  -username string
        Username for switch configuration

PS C:\Downloads\switchConfigTool> .\SwitchConfigGenerator.exe -switchLib .\Input\switchLibrary\ -inputJsonFile .\Input\testInput\cisco_bmc_bgp_input.json -outputFolder outputConifg

PS C:\Downloads\switchConfigTool\outputConifg> ls

Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----          3/4/2023   2:18 PM          14762 BMC.config
-a----          3/4/2023   2:18 PM          14082 BMC.json
-a----          3/4/2023   2:18 PM          34018 TOR1.config
-a----          3/4/2023   2:18 PM          27093 TOR1.json
-a----          3/4/2023   2:18 PM          34019 TOR2.config
-a----          3/4/2023   2:18 PM          27100 TOR2.json
```

#### If username and password are NOT provided

- username is `'azureadmin-'+ 5 * random characters`.
- password is 16 characters strong string includes `3 * UpperChar + 3 * Num + 3 * SpecialChar + 7 * LowerChar`.

Sample Configuration:

```text
username admin password 0 c2!cuh%xMC8*aRt9 role network-admin
username azureadmin-vlzxb password 0 c2!cuh%xMC8*aRt9 role network-admin
```

#### If username and password are provided

The configuration will use provided username/password in the configuration.

```powershell
PS C:\Downloads\switchConfigTool> .\SwitchConfigGenerator.exe -switchLib .\Input\switchLibrary\ -inputJsonFile .\Input\testInput\cisco_bmc_bgp_input.json -outputFolder outputConifg -username admintest -password admintest
```

Sample Configuration:

```text
username admin password 0 admintest role network-admin
username admintest password 0 admintest role network-admin
```
