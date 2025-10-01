# IP Subnet Planning Tool - Quick Guide

## What This Tool Does

🎯 **Automatically divides your network into subnets** - no manual calculations needed!

- Takes a big network (like `192.168.1.0/24`)
- Splits it into smaller subnets based on your needs
- Shows you exactly what IP ranges to use
- Calculates everything automatically - no subnet math required!

---

## 🚀 Getting Started

### Step 1: Import the Module

```powershell
# Load the IP Management tool
Import-Module .\IPManagement.psm1
```

### Step 2: Get Help Anytime

```powershell
# See all available functions
Get-Command -Module IPManagement

# Get detailed help for any function
Get-Help New-SubnetPlanFromConfig -Examples
Get-Help New-SubnetPlanByHosts -Examples
Get-Help New-SubnetPlan -Examples

# Get full documentation
Get-Help New-SubnetPlanFromConfig -Full
```

💡 **Pro Tip:** Always run `Get-Help [FunctionName] -Examples` to see more examples and detailed usage!

---

## Three Ways to Use It (Pick One)

### 🏆 **Option 1: JSON Configuration** (RECOMMENDED)

**Best for:** Real projects with named subnets

**When to use:** When you want professional documentation and have multiple subnets with names/VLANs/descriptions.

```powershell
# 1. Create a simple JSON file (save as network.json):
{
  "network": "192.168.1.0/24",
  "subnets": [
    { "name": "Management", "vlan": 101, "hosts": 30 },
    { "name": "Users", "vlan": 102, "hosts": 100 },
    { "name": "Servers", "vlan": 103, "hosts": 20 }
  ]
}

# 2. Run the command (network is now defined in JSON!):
New-SubnetPlanFromConfig -ConfigPath "network.json"

# Alternative: Override network from command line
New-SubnetPlanFromConfig -Network "10.0.1.0/24" -ConfigPath "network.json"
```

**Sample Output:**

```text
Name       Vlan Subnet           Prefix Network       Broadcast     FirstHost     EndHost       UsableHosts TotalIPs
----       ---- ------           ------ -------       ---------     ---------     -------       ----------- --------
Users       102 192.168.1.0/25      25 192.168.1.0   192.168.1.127 192.168.1.1   192.168.1.126        126      128
Management  101 192.168.1.128/27    27 192.168.1.128 192.168.1.159 192.168.1.129 192.168.1.158         30       32
Servers     103 192.168.1.160/27    27 192.168.1.160 192.168.1.191 192.168.1.161 192.168.1.190         30       32
Available       192.168.1.192/26    26 192.168.1.192 192.168.1.255 192.168.1.193 192.168.1.254         62       64
```

**You get:** Professional table with subnet names, VLANs, IP ranges, and everything labeled clearly—plus the `TotalIPs` column so you can immediately see the total address count in each block.

#### 🔧 **JSON Configuration Flexibility**

**Two ways to specify subnet sizes in JSON:**

**Option A: By Host Count** (most common)

```json
{
  "network": "192.168.1.0/24",
  "subnets": [
    { "name": "Management", "vlan": 101, "hosts": 30 },
    { "name": "DMZ", "vlan": 102, "hosts": 50 }
  ]
}
```

**Option B: By CIDR Prefix** (for network experts)

```json
{
  "network": "192.168.1.0/24",
  "subnets": [
    { "name": "Management", "vlan": 101, "cidr": "27" },
    { "name": "DMZ", "vlan": 102, "cidr": "26" }
  ]
}
```

**Mixed Approach** (you can combine both in one file!)

```json
{
  "network": "192.168.1.0/24",
  "subnets": [
    { "name": "Management", "vlan": 101, "hosts": 30 },
    { "name": "DMZ", "vlan": 102, "cidr": "26" },
    { "name": "Servers", "vlan": 103, "hosts": 15 }
  ]
}
```

💡 **Pro Tip:** Use `hosts` when you know how many devices you need. Use `cidr` when you need exact subnet sizes.

#### 🆕 Named IP Assignments (per subnet)

Need specific hosts (gateway, VIPs, BMC, etc.) to land on predictable addresses? Add an `IPAssignments` array to any subnet and the tool will:

- Reserve hosts by position (starting at 1 = first usable IP)
- Show each named assignment inline with the subnet
- Highlight unused address ranges between assignments
- Include the broadcast address so nothing gets overlooked

```json
{
  "network": "10.0.0.0/24",
  "subnets": [
    {
      "name": "Mgmt",
      "vlan": 110,
      "cidr": "28",
      "IPAssignments": [
        { "Name": "Gateway", "Position": 1 },
        { "Name": "VMM", "Position": 3 }
      ]
    }
  ]
}
```

```powershell
$json = Get-Content '.\mgmt-subnet.json' -Raw
New-SubnetPlanFromConfig -JsonConfig $json
```

**Sample Output:**

```text
Name  Vlan Subnet         Label         IP                       TotalIPs Prefix Mask             Category
----  ---- ------         -----         --                       -------- ------ ----             --------
Mgmt  110  10.0.0.0/28    Network       10.0.0.0                        1 /28   255.255.255.240  Network
Mgmt  110  10.0.0.0/28    Gateway       10.0.0.1                        1 /28   255.255.255.240  Assignment
Mgmt  110  10.0.0.0/28    Unused Range  10.0.0.2 - 10.0.0.2             1 /28   255.255.255.240  Unused
Mgmt  110  10.0.0.0/28    VMM           10.0.0.3                        1 /28   255.255.255.240  Assignment
Mgmt  110  10.0.0.0/28    Unused Range  10.0.0.4 - 10.0.0.14           11 /28   255.255.255.240  Unused
Mgmt  110  10.0.0.0/28    Broadcast     10.0.0.15                       1 /28   255.255.255.240  Broadcast
Available     10.0.0.16/28 Available Range 10.0.0.17 - 10.0.0.30       14 /28   255.255.255.240  Available
```

🔎 **How positions work:** position `1` = first usable host, `2` = second usable host, and so on. Positions must be unique and stay within the usable host count of the subnet.

#### 🗂️ Export the Plan Automatically (New!)

Want clean artifacts for handoffs or automation? `New-SubnetPlanFromConfig` now saves the detailed plan directly to disk:

- `-ExportJsonPath` → writes the full detail table to JSON (perfect for tooling and APIs)
- `-ExportCsvPath` → CSV formatted output for spreadsheets and reporting
- `-ExportMarkdownPath` → Markdown table that drops straight into wikis and docs

You can use any or all of them in a single run. The cmdlet creates folders on the fly, so paths like `output\plans\` just work.

```powershell
# Save the plan while still seeing a formatted table in the console
New-SubnetPlanFromConfig -ConfigPath "network.json" `
  -ExportJsonPath ".\output\network-plan.json" `
  -ExportCsvPath ".\output\network-plan.csv" `
  -ExportMarkdownPath ".\output\network-plan.md"
```

```powershell
# Combine inline JSON with Markdown export
$json = Get-Content '.\mgmt-subnet.json' -Raw
New-SubnetPlanFromConfig -JsonConfig $json -ExportMarkdownPath '.\output\mgmt-plan.md'
```

✅ Console output still shows the pretty table (unless you pipe/assign it), while files capture the full detail for later.

#### 🌐 Multiple Primary Networks (New!)

Need to manage multiple independent networks (like separate compute and management networks)? The JSON configuration now supports an **array of network objects**, where each network is independently subnetted:

```json
[
  {
    "network": "192.168.1.0/24",
    "subnets": [
      { "name": "csu-edge-transport-compute", "vlan": "203", "cidr": "28" },
      { "name": "csu-exchange-compute", "vlan": "102", "cidr": "27" },
      { "name": "msu-compute", "vlan": "302", "cidr": "27" }
    ]
  },
  {
    "network": "10.50.1.0/24",
    "subnets": [
      { "name": "csu-edge-transport-management", "vlan": "101", "cidr": "27" },
      { "name": "msu-management", "vlan": "201", "cidr": "27" },
      { "name": "csu-exchange-management", "vlan": "301", "cidr": "26" }
    ]
  }
]
```

```powershell
# Process multiple networks from a file
New-SubnetPlanFromConfig -ConfigPath "multi-network.json"

# Or use inline JSON
$multiNetworkJson = Get-Content '.\multi-network.json' -Raw
New-SubnetPlanFromConfig -JsonConfig $multiNetworkJson
```

**Key features:**
- Each network is processed independently with its own subnets
- Supports all existing features (IPAssignments, VLAN tags, custom properties, etc.)
- Results are combined into a single output table showing all networks
- Fully backward compatible - single network objects still work as before

**Sample Output:**

```text
Subnet           Name                          Vlan Label           IP                            TotalIPs
------           ----                          ---- -----           --                            --------
192.168.1.0/27   csu-exchange-compute          102  Network         192.168.1.0                          1
192.168.1.0/27   csu-exchange-compute          102  Unused Range    192.168.1.1 - 192.168.1.30          30
192.168.1.0/27   csu-exchange-compute          102  Broadcast       192.168.1.31                         1
192.168.1.32/27  msu-compute                   302  Network         192.168.1.32                         1
...
10.50.1.0/26     csu-exchange-management       301  Network         10.50.1.0                            1
10.50.1.0/26     csu-exchange-management       301  Unused Range    10.50.1.1 - 10.50.1.62              62
10.50.1.0/26     csu-exchange-management       301  Broadcast       10.50.1.63                           1
...
```

💡 **Use case:** Perfect for environments where you need to separately manage compute networks (192.168.x.x) and management networks (10.x.x.x), each with their own independent subnetting.

---

### 🎯 **Option 2: Host Count** (SIMPLE)

**Best for:** Quick calculations

**When to use:** When you just know "I need X hosts" and want a fast answer.

```powershell
# "I need 2 subnets with 50 hosts each, and 1 subnet with 10 hosts"
New-SubnetPlanByHosts -Network "192.168.1.0/24" -HostRequirements @{ 50 = 2; 10 = 1 }
```

**Sample Output:**

```text
Name     Subnet           Prefix Network       Broadcast     FirstHost     EndHost       UsableHosts TotalIPs
----     ------           ------ -------       ---------     ---------     -------       ----------- --------
Assigned 192.168.1.0/26      26 192.168.1.0   192.168.1.63  192.168.1.1   192.168.1.62           62       64
Assigned 192.168.1.64/26     26 192.168.1.64  192.168.1.127 192.168.1.65  192.168.1.126          62       64
Assigned 192.168.1.128/28    28 192.168.1.128 192.168.1.143 192.168.1.129 192.168.1.142          14       16
Available 192.168.1.144/28   28 192.168.1.144 192.168.1.159 192.168.1.145 192.168.1.158          14       16
Available 192.168.1.160/27   27 192.168.1.160 192.168.1.191 192.168.1.161 192.168.1.190          30       32
Available 192.168.1.192/26   26 192.168.1.192 192.168.1.255 192.168.1.193 192.168.1.254          62       64
```

**You get:** Subnet ranges that fit your host requirements exactly.

---

### ⚙️ **Option 3: Technical Prefixes** (ADVANCED)

**Best for:** Network experts who know CIDR

**When to use:** When you know exactly what subnet sizes you want (like /26, /28).

```powershell
# "I need 2 /26 subnets and 3 /28 subnets"
New-SubnetPlan -Network "192.168.1.0/24" -PrefixRequirements @{ 26 = 2; 28 = 3 }
```

**Sample Output:**

```text
Name      Subnet           Prefix Network       Broadcast     FirstHost     EndHost       UsableHosts TotalIPs
----      ------           ------ -------       ---------     ---------     -------       ----------- --------
Assigned  192.168.1.0/26      26 192.168.1.0   192.168.1.63  192.168.1.1   192.168.1.62           62       64
Assigned  192.168.1.64/26     26 192.168.1.64  192.168.1.127 192.168.1.65  192.168.1.126          62       64
Assigned  192.168.1.128/28    28 192.168.1.128 192.168.1.143 192.168.1.129 192.168.1.142          14       16
Assigned  192.168.1.144/28    28 192.168.1.144 192.168.1.159 192.168.1.145 192.168.1.158          14       16
Assigned  192.168.1.160/28    28 192.168.1.160 192.168.1.175 192.168.1.161 192.168.1.174          14       16
Available 192.168.1.176/28    28 192.168.1.176 192.168.1.191 192.168.1.177 192.168.1.190          14       16
Available 192.168.1.192/26    26 192.168.1.192 192.168.1.255 192.168.1.193 192.168.1.254          62       64
```

**You get:** Precise control over subnet sizes with technical CIDR notation.

---

## 📊 JSON Output for APIs and Automation

All functions support JSON output using the `-AsJson` parameter - perfect for:

- **API Integration:** Consume data in web applications and REST APIs
- **Infrastructure as Code:** Use with Terraform, ARM templates, or Ansible
- **Data Export:** Save results to files for documentation or further processing
- **PowerShell Automation:** Pass structured data between scripts

### Quick Examples

```powershell
# Get JSON output for API consumption
$jsonResult = New-SubnetPlanByHosts -Network "192.168.1.0/24" -HostRequirements @{ 50 = 2; 10 = 3 } -AsJson

# Export to file for documentation or other tools
New-SubnetPlanFromConfig -ConfigPath "network.json" -AsJson | Out-File "subnet-plan.json"

# Parse JSON in PowerShell for further processing
$data = New-SubnetPlan -Network "10.0.0.0/22" -PrefixRequirements @{ 26 = 4 } -AsJson | ConvertFrom-Json
foreach ($subnet in $data) {
    Write-Host "Subnet: $($subnet.Subnet) has $($subnet.UsableHosts) usable hosts"
}
```

### Sample JSON Output

```json
[
  {
    "Name": "Assigned",
    "Subnet": "192.168.1.0/26",
    "Prefix": 26,
    "Network": "192.168.1.0",
    "Broadcast": "192.168.1.63",
    "FirstHost": "192.168.1.1",
    "EndHost": "192.168.1.62",
    "UsableHosts": 62,
    "TotalIPs": 64
  },
  {
    "Name": "Available",
    "Subnet": "192.168.1.64/26",
    "Prefix": 26,
    "Network": "192.168.1.64",
    "Broadcast": "192.168.1.127",
    "FirstHost": "192.168.1.65",
    "EndHost": "192.168.1.126",
    "UsableHosts": 62,
    "TotalIPs": 64
  }
]
```

💡 **Pro Tip:** JSON output maintains all the same data as table output but in a structured format perfect for automation!

🧮 **Quick glance:** The new `TotalIPs` field/column makes it easy to compare available address pools without doing mental math.
