
# Customize Your Own Switch Library

The tool provides a framework for defining and customizing your switch library based on the `input.json` and `switch library` JSON files. This document is focuing on `switch library`, below is the high-level structure of default `switchLib`:


```
.
├── cisco
│   └── 9.3(9)
│       ├── 93180yc-fx
│       │   └── interface.json
│       │   └── bgp.json
│       │   └── static.json
│       ├── 9348gc-fxp
│       │   └── interface.json
│       └── template
│           ├── AllConfig.go.tmpl
│           ├── hostname.go.tmpl
│           ├── port.go.tmpl
│           ├── settings.go.tmpl
│           ├── stig.go.tmpl
│           └── vlan.go.tmpl
└── dellemc
    └── 10.5(3.4)
        ├── n3248te-on
        │   └── interface.json
        ├── s5248-on
        │   └── interface.json
        │   └── bgp.json
        └── template
            ├── AllConfig.go.tmpl
            ├── hostname.go.tmpl
            ├── port.go.tmpl
            ├── settings.go.tmpl
            ├── stig.go.tmpl
            └── vlan.go.tmpl
```

- You can modify any files within the switchLib directory to suit your requirements. However, in most cases, the files that typically need adjustment when adding new device models are interface.json and bgp.json.

- All you need to ensure is that the tool can locate the correct path to your resources.

 
### Notes:
- The files in the `template` folder are written using **Go Template**, which requires advanced knowledge to modify.
- Ensure you create a backup copy of the files before making any direct modifications.


# Examples
## Updating Port Assignment

If you are using a Cisco 93180YC-FX switch running `9.3(11)` as a Top-of-Rack (TOR) switch and want to assign port `1/23` to the Storage port group, follow these steps:

1. Navigate to the `interface.json` file located in the directory: `input/switchLib/cisco/9.3(11)/93180yc-fx/interface.json`

2. Open the file and add `1/23` to the appropriate list.

```
        {
            "Function": "Storage",
            "Port": [
                "1/17",
                "1/18",
                "1/19",
                "1/20",
                "1/21",
                "1/22",
				"1/23"
            ]
        }
```
3. Save the file after making the changes.

The tool will automatically read the updated `interface.json` and generate the new configuration based on your modifications.


## Updating Existing Switch with a New OS Image

If you upgraded the Cisco 93180YC-FX switch from `9.3(11)` to `10.4`, follow these steps to add support for the new version:

1. Copy the entire directory: `input/switchLib/cisco/9.3(11)`
2. Rename the copied directory to: `input/switchLib/cisco/10.4`
3. That's it! The new version is now supported.

### Notes:
1. Ensure that your `input.json` file is updated with the correct version (`10.4`) to reflect the changes.
2. In most cases, vendors do not significantly change configuration commands across versions, so there is no need to modify the configuration templates.
3. If upgrading across a major version gap, double-check with the vendor to confirm that the existing configurations remain compatible.

## Adding New Switch Model

If you are adding a new switch model, such as the Cisco 93190YC-FX switch running on `10.4`, follow these steps to add support for it:

1. Make sure direcotry `input/switchLib/cisco/10.4` existing.
2. Copy the entire directory: `input/switchLib/cisco/10.4/93180yc-fx`
3. Rename the copied directory to: `input/switchLib/cisco/10.4/93190yc-fx`
4. Update the json files (interface, bgp. etc) under the new folder accordingly.

