function New-SubnetPlanByHosts {
  <#
  .SYNOPSIS
      Creates an optimized IP subnet allocation plan based on host count requirements.

  .DESCRIPTION
      This user-friendly function allows network engineers to specify subnet requirements by the number
      of hosts needed rather than complex CIDR prefix lengths. The function automatically calculates the
      optimal subnet sizes and allocates them efficiently across the available IP space.
      
      Key features:
      - Input by host count - more intuitive for network planning
      - Automatic calculation of optimal subnet sizes
      - Uses largest-first allocation strategy for maximum efficiency
      - Shows remaining space as available blocks for future use
      - Validates that requirements can be met within the parent network
      
  .PARAMETER Network
      The parent network in CIDR notation (e.g., "192.168.1.0/24").
      Must be a valid IPv4 network address with prefix length.

  .PARAMETER HostRequirements
      Hashtable specifying subnet requirements where keys represent the number of hosts needed
      and values indicate how many subnets of that size are required.
      The function automatically calculates the smallest subnet size that can accommodate each host count.
      Example: @{ 25 = 2; 11 = 3; 5 = 1 } means "2 subnets with 25 hosts, 3 subnets with 11 hosts, 1 subnet with 5 hosts"

  .EXAMPLE
      New-SubnetPlanByHosts -Network "192.168.1.0/24" -HostRequirements @{ 25 = 2; 11 = 2 }
      
      Creates 2 subnets that can accommodate 25 hosts each (/27 subnets with 30 usable hosts)
      and 2 subnets that can accommodate 11 hosts each (/28 subnets with 14 usable hosts).

  .EXAMPLE
      New-SubnetPlanByHosts -Network "10.0.0.0/22" -HostRequirements @{ 100 = 1; 50 = 2; 10 = 5 }
      
      Creates subnets for: 1×100 hosts (/25=126 hosts), 2×50 hosts (/26=62 hosts each), 5×10 hosts (/28=14 hosts each).

  .EXAMPLE
      New-SubnetPlanByHosts -Network "172.16.0.0/24" -HostRequirements @{ 60 = 1; 25 = 1; 10 = 2; 5 = 3 }
      
      Mixed requirements: 60 hosts (/26), 25 hosts (/27), 2×10 hosts (/28 each), 3×5 hosts (/29 each).

  .EXAMPLE
      New-SubnetPlanByHosts -Network "192.168.1.0/24" -HostRequirements @{ 25 = 2; 11 = 2 } -AsJson
      
      Returns the subnet plan in JSON format for API consumption or further processing.

  .OUTPUTS
      Array of PSCustomObject containing detailed subnet information including
      network addresses, broadcast addresses, usable host ranges, and allocation status.
      When -AsJson is specified, returns JSON string representation of the data.

  .NOTES
      Host count calculation:
      - The function finds the smallest subnet size that can accommodate the requested hosts
      - /30 subnets = 2 usable hosts, /29 = 6 hosts, /28 = 14 hosts, /27 = 30 hosts, etc.
      - Always allocates enough space to meet or exceed the host requirement
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory = $true, HelpMessage = "IPv4 network in CIDR format (e.g., '192.168.1.0/24')")]
    [ValidatePattern('^(\d{1,3}\.){3}\d{1,3}/([1-2]?[0-9]|3[0-2])$')]
    [string]$Network,
    
    [Parameter(Mandatory = $true, HelpMessage = "Hashtable of host counts and required quantities")]
    [ValidateNotNullOrEmpty()]
    [hashtable]$HostRequirements,
    
    [Parameter(Mandatory = $false, HelpMessage = "Output results in JSON format")]
    [switch]$AsJson
  )

  # Begin processing with input validation and parameter setup
  Write-Verbose "Starting subnet planning by host count for network: $Network"
  
  # Function to calculate the minimum prefix length needed for a given number of hosts
  function Get-PrefixForHosts {
    [CmdletBinding()]
    param([int]$HostCount)
    
    if ($HostCount -le 0) {
      throw "Host count must be greater than 0"
    }
    
    # Need to account for network and broadcast addresses (except for /31 and /32)
    # For most subnets: usable hosts = total addresses - 2
    # We need: 2^(32-prefix) - 2 >= HostCount
    # So: 2^(32-prefix) >= HostCount + 2
    # Therefore: 32-prefix >= log2(HostCount + 2)
    # So: prefix <= 32 - log2(HostCount + 2)
    
    $requiredAddresses = $HostCount + 2  # Add network and broadcast
    $bitsNeeded = [Math]::Ceiling([Math]::Log($requiredAddresses, 2))
    $prefix = 32 - $bitsNeeded
    
    # Ensure prefix is within valid range
    if ($prefix -lt 1) { $prefix = 1 }
    if ($prefix -gt 30) { $prefix = 30 }  # /31 and /32 are special cases
    
    return $prefix
  }
  
  # Convert host counts to prefix lengths
  $PrefixRequirements = @{}
  
  Write-Verbose "Converting host requirements to subnet sizes:"
  foreach ($hostCount in $HostRequirements.Keys) {
    $requiredSubnets = $HostRequirements[$hostCount]
    $optimalPrefix = Get-PrefixForHosts -HostCount $hostCount
    $actualHosts = [Math]::Pow(2, 32 - $optimalPrefix) - 2
    
    Write-Verbose "  $requiredSubnets × $hostCount hosts → /$optimalPrefix subnets ($actualHosts usable hosts each)"
    
    # Add to prefix counts (handle multiple host counts that map to same prefix)
    if ($PrefixRequirements.ContainsKey($optimalPrefix)) {
      $PrefixRequirements[$optimalPrefix] += $requiredSubnets
    }
    else {
      $PrefixRequirements[$optimalPrefix] = $requiredSubnets
    }
  }
  
  # Call the optimized subnet planning function with calculated prefix lengths
  Write-Verbose "Calling New-SubnetPlan with calculated prefix requirements"
  $result = New-SubnetPlan -Network $Network -PrefixRequirements $PrefixRequirements
  
  if ($AsJson) {
    return $result | ConvertTo-Json -Depth 10
  }
  else {
    return $result
  }
}

function New-SubnetPlan {
  <#
  .SYNOPSIS
      Creates an optimized IPv4 subnet allocation plan with efficient space utilization.

  .DESCRIPTION
      This function subdivides an IPv4 parent network into multiple subnets according to specified
      prefix length requirements. It employs a largest-first allocation strategy to maximize space
      efficiency and minimize fragmentation.
      
      Key features:
      - Optimized allocation algorithm ensures minimal IP address waste
      - Largest-first strategy prevents fragmentation
      - Shows both assigned and available network blocks
      - Comprehensive validation of all subnet requirements
      - Professional output format suitable for network documentation
      
  .PARAMETER Network
      The parent network in CIDR notation (e.g., "192.168.1.0/24").
      Must be a valid IPv4 network address with prefix length.

  .PARAMETER PrefixRequirements
      Hashtable specifying subnet requirements where keys represent prefix lengths
      and values indicate how many subnets of that prefix are required.
      Example: @{ 26 = 2; 28 = 3 } means "2 subnets with /26 prefix, 3 subnets with /28 prefix"

  .EXAMPLE
      New-SubnetPlan -Network "192.168.1.0/24" -PrefixRequirements @{ 26 = 2; 28 = 1 }
      
      Creates 2 /26 subnets (62 hosts each) and 1 /28 subnet (14 hosts).

  .EXAMPLE
      New-SubnetPlan -Network "10.0.1.0/24" -PrefixRequirements @{ 27 = 2; 28 = 1 }
      
      Output:
      Network: 10.0.1.0/27, Broadcast: 10.0.1.31, Hosts: 10.0.1.1-10.0.1.30, Category: Assigned
      Network: 10.0.1.32/27, Broadcast: 10.0.1.63, Hosts: 10.0.1.33-10.0.1.62, Category: Assigned
      Network: 10.0.1.64/28, Broadcast: 10.0.1.79, Hosts: 10.0.1.65-10.0.1.78, Category: Assigned
      Network: 10.0.1.80/28, Broadcast: 10.0.1.95, Hosts: 10.0.1.81-10.0.1.94, Category: Available

  .EXAMPLE
      New-SubnetPlan -Network "192.168.1.0/24" -PrefixRequirements @{ 28 = 3; 26 = 1; 29 = 2 }
      
      Mixed subnet sizes: 1×/26 (62 hosts), 3×/28 (14 hosts each), 2×/29 (6 hosts each).

  .EXAMPLE
      New-SubnetPlan -Network "192.168.1.0/24" -PrefixRequirements @{ 26 = 2; 28 = 1 } -AsJson
      
      Returns the subnet plan in JSON format suitable for APIs or automated processing.

  .OUTPUTS
      Array of PSCustomObject containing comprehensive subnet information including
      network addresses, broadcast addresses, usable host ranges, and allocation status.
      When -AsJson is specified, returns JSON string representation of the data.

  .NOTES
      The function uses a greedy bin-packing algorithm to achieve optimal space utilization.
      All addresses are automatically calculated using binary arithmetic for accuracy.
  #>
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory = $true, HelpMessage = "IPv4 network in CIDR format (e.g., '192.168.1.0/24')")]
    [ValidatePattern('^(\d{1,3}\.){3}\d{1,3}/([1-2]?[0-9]|3[0-2])$')]
    [string]$Network,
    
    [Parameter(Mandatory = $true, HelpMessage = "Hashtable of prefix lengths and required quantities")]
    [ValidateNotNullOrEmpty()]
    [hashtable]$PrefixRequirements,
    
    [Parameter(Mandatory = $false, HelpMessage = "Output results in JSON format")]
    [switch]$AsJson
  )

  # Begin processing with input validation and parameter setup
  Write-Verbose "Starting subnet allocation for network: $Network"
  
  # Internal helper functions for IP address manipulation and subnet calculations
  
  # Converts IPv4 dotted decimal notation to 32-bit unsigned integer
  function Convert-IpToInt {
    [CmdletBinding()]
    param([string]$Ip)
    
    $octets = $Ip.Split('.') | ForEach-Object { [uint32]$_ }
    return ($octets[0] -shl 24) -bor ($octets[1] -shl 16) -bor ($octets[2] -shl 8) -bor $octets[3]
  }
  
  # Converts 32-bit unsigned integer to IPv4 dotted decimal notation
  function Convert-IntToIp {
    [CmdletBinding()]
    param([uint32]$Value)
    
    $octet1 = ($Value -shr 24) -band 0xFF
    $octet2 = ($Value -shr 16) -band 0xFF
    $octet3 = ($Value -shr 8) -band 0xFF
    $octet4 = $Value -band 0xFF
    return "$octet1.$octet2.$octet3.$octet4"
  }
  
  # Calculates the total number of addresses in a subnet based on prefix length
  function Get-BlockSize {
    [CmdletBinding()]
    param([int]$Prefix)
    
    if ($Prefix -lt 0 -or $Prefix -gt 32) {
      throw "Invalid prefix length: $Prefix. Must be between 0 and 32."
    }
    return [uint32][Math]::Pow(2, 32 - $Prefix)
  }
  
  # Determines the network base address for a given IP and prefix length
  function Get-NetworkBase {
    [CmdletBinding()]
    param([uint32]$IpInt, [int]$Prefix)
    
    if ($Prefix -eq 0) {
      return [uint32]0
    }
    $hostBits = 32 - $Prefix
    $mask = [uint32]([Math]::Pow(2, 32) - [Math]::Pow(2, $hostBits))
    return ($IpInt -band $mask)
  }
  
  # Creates a comprehensive subnet information object
  function New-Block {
    [CmdletBinding()]
    param([uint32]$StartInt, [int]$Prefix, [string]$Category)
    
    $size = Get-BlockSize -Prefix $Prefix
    $network = $StartInt
    $broadcast = $StartInt + $size - 1
    
    # Calculate usable address range (excluding network and broadcast for subnets > /30)
    $usableStart = if ($size -gt 2) { $network + 1 } else { $null }
    $usableEnd = if ($size -gt 2) { $broadcast - 1 } else { $null }
    $usableCount = [int]([Math]::Max(0, $size - 2))
    
    return [PSCustomObject]@{
      Name        = $Category
      Prefix      = $Prefix
      Subnet      = "$(Convert-IntToIp -Value $network)/$Prefix"
      Network     = Convert-IntToIp -Value $network
      Broadcast   = Convert-IntToIp -Value $broadcast
      FirstHost   = if ($usableStart) { Convert-IntToIp -Value $usableStart } else { $null }
      EndHost     = if ($usableEnd) { Convert-IntToIp -Value $usableEnd } else { $null }
      UsableHosts = $usableCount
      TotalIPs    = $size
    }
  }
  
  # Finds the highest power of 2 that is less than or equal to the given number
  function Get-HighestPowerOfTwoLE {
    [CmdletBinding()]
    param([uint32]$Number)
    
    if ($Number -eq 0) { return 0 }
    $power = [Math]::Floor([Math]::Log([double]$Number, 2))
    return [uint32][Math]::Pow(2, $power)
  }
  
  # Implements greedy bin-packing algorithm for optimal space utilization
  function Invoke-GreedyPacking {
    [CmdletBinding()]
    param([uint32]$StartInt, [uint32]$RemainingSpace)
    
    $blocks = @()
    $currentPosition = $StartInt
    $remainingBytes = $RemainingSpace
    
    Write-Verbose "Starting greedy packing at position $(Convert-IntToIp -Value $StartInt) with $RemainingSpace addresses"
    
    while ($remainingBytes -gt 0) {
      # Calculate alignment boundary (largest power of 2 that divides current position)
      $alignmentSize = [uint32]($currentPosition -band (-$currentPosition))
      if ($alignmentSize -eq 0) { $alignmentSize = 1 }
      
      # Find optimal block size considering both alignment and remaining space
      $maxSpaceSize = Get-HighestPowerOfTwoLE -Number $remainingBytes
      $optimalSize = [uint32]([Math]::Min($alignmentSize, $maxSpaceSize))
      
      # Calculate corresponding prefix length
      $prefix = 32 - [int][Math]::Log([double]$optimalSize, 2)
      
      # Create available block and update counters
      $blocks += New-Block -StartInt $currentPosition -Prefix $prefix -Category 'Available'
      $currentPosition += $optimalSize
      $remainingBytes -= $optimalSize
      
      Write-Verbose "Allocated available block: $(Convert-IntToIp -Value ($currentPosition - $optimalSize))/$prefix"
    }
    
    return $blocks
  }
  # Main processing logic begins here
  try {
    # Parse and validate the supernet input
    $parts = $Network.Split('/')
    $baseIpAddress = $parts[0]
    $parentPrefix = [int]$parts[1]
    
    # Convert base IP to integer and validate it's a proper network address
    $baseInt = Convert-IpToInt -Ip $baseIpAddress
    $parentNetworkBase = Get-NetworkBase -IpInt $baseInt -Prefix $parentPrefix
    
    if ($parentNetworkBase -ne $baseInt) {
      throw "The IP address '$baseIpAddress' is not the network address for the network '$Network'. Expected: $(Convert-IntToIp -Value $parentNetworkBase)/$parentPrefix"
    }
    
    $parentNetworkSize = Get-BlockSize -Prefix $parentPrefix
    Write-Verbose "Parent network validated: $Network contains $parentNetworkSize addresses"
    
    # Validate child subnet requirements and calculate total space needed
    $totalRequiredSpace = [uint32]0
    foreach ($prefixLength in $PrefixRequirements.Keys) {
      if ($prefixLength -lt $parentPrefix) {
        throw "Child prefix /$prefixLength cannot be larger than parent prefix /$parentPrefix"
      }
      if ($prefixLength -lt 0 -or $prefixLength -gt 32) {
        throw "Invalid prefix length: /$prefixLength. Must be between /$($parentPrefix + 1) and /32"
      }
      
      $subnetSize = Get-BlockSize -Prefix $prefixLength
      $requiredCount = [uint32]$PrefixRequirements[$prefixLength]
      $totalRequiredSpace += $requiredCount * $subnetSize
      
      Write-Verbose "Requirement: $requiredCount × /$prefixLength subnets ($subnetSize addresses each)"
    }
    
    if ($totalRequiredSpace -gt $parentNetworkSize) {
      $excessSpace = $totalRequiredSpace - $parentNetworkSize
      throw "Required address space ($totalRequiredSpace addresses) exceeds supernet capacity ($parentNetworkSize addresses) by $excessSpace addresses"
    }
    
    Write-Verbose "Space validation passed: $totalRequiredSpace/$parentNetworkSize addresses required"
    
    # Determine subnet allocation order using optimal largest-first strategy
    Write-Verbose "Using LargestFirst allocation strategy for optimal space utilization"
    $orderedPrefixes = $PrefixRequirements.Keys | Sort-Object { [int]$_ }
    
    # Build ordered list of assigned subnets for allocation
    $allocationQueue = @()
    foreach ($prefix in $orderedPrefixes) {
      $count = [int]$PrefixRequirements[$prefix]
      for ($i = 0; $i -lt $count; $i++) { 
        $allocationQueue += [int]$prefix 
      }
    }
    
    Write-Verbose "Allocation queue prepared with $($allocationQueue.Count) subnets"
    
    # Allocate assigned subnets sequentially
    $allocatedSubnets = @()
    $currentPosition = [uint32]$parentNetworkBase
    
    foreach ($prefixLength in $allocationQueue) {
      $subnetSize = Get-BlockSize -Prefix $prefixLength
      
      # Verify alignment requirements are met
      if (($currentPosition % $subnetSize) -ne 0) {
        throw "Alignment error: Cannot place /$prefixLength subnet at address $(Convert-IntToIp -Value $currentPosition). Address must be aligned to $subnetSize-byte boundary."
      }
      
      # Create and record the assigned subnet
      $allocatedSubnets += New-Block -StartInt $currentPosition -Prefix $prefixLength -Category 'Assigned'
      $currentPosition += $subnetSize
      
      Write-Verbose "Allocated assigned subnet: $(Convert-IntToIp -Value ($currentPosition - $subnetSize))/$prefixLength"
    }
    
    # Calculate and allocate remaining address space using greedy packing
    $usedSpace = [uint32]($currentPosition - $parentNetworkBase)
    $remainingSpace = $parentNetworkSize - $usedSpace
    
    if ($remainingSpace -gt 0) {
      Write-Verbose "Processing remaining space: $remainingSpace addresses starting at $(Convert-IntToIp -Value $currentPosition)"
      $availableSubnets = Invoke-GreedyPacking -StartInt $currentPosition -RemainingSpace $remainingSpace
      $allocatedSubnets += $availableSubnets
    }
    else {
      Write-Verbose "No remaining space to allocate - supernet fully utilized"
    }
    
    # Create optimized output for display
    $summaryView = $allocatedSubnets | Select-Object Name, Subnet, Prefix, Network, Broadcast, FirstHost, EndHost, UsableHosts, TotalIPs
    
    Write-Verbose "Subnet allocation completed successfully. Total subnets created: $($allocatedSubnets.Count)"
    
    # Handle JSON output format
    if ($AsJson) {
      return $summaryView | ConvertTo-Json -Depth 10
    }
    
    # Check if output is being captured (assigned to variable) or piped
    if ($MyInvocation.Line -match '\$\w+\s*=' -or $MyInvocation.Line -match '\|') {
      # Return data objects when captured or piped
      return $summaryView
    }
    else {
      # Display clean table when run directly with proper spacing
      $tableOutput = $summaryView | Format-Table -AutoSize | Out-String
      Write-Host $tableOutput.TrimEnd()
      Write-Host ""  # Add blank line for better readability
    }
  }
  catch {
    Write-Error "Failed to process CIDR subdivision: $($_.Exception.Message)"
    throw
  }
}

function New-SubnetPlanFromConfig {
  <#
  .SYNOPSIS
      Creates an optimized IP subnet allocation plan from JSON configuration with named subnets.

  .DESCRIPTION
      This function reads network requirements from JSON configuration and generates a subnet plan
      with descriptive names for each subnet. Perfect for infrastructure planning and documentation.
      
      Key features:
      - JSON-based configuration for easy management and version control
      - Named subnets with descriptions for better documentation
      - Enhanced output showing subnet names and purposes
      - Supports both file-based and direct JSON string input
      - Network definition can be included in JSON configuration for single-file setup
      - Supports both host count requirements and direct CIDR prefix specification
      - Allows named IP assignments within each subnet for deterministic host placement
      - Supports multiple primary networks in a single configuration (array of network objects)
      - Validates configuration and provides detailed error messages
      
  .PARAMETER Network
      The parent network in CIDR notation (e.g., "192.168.1.0/24").
      Optional when network is defined in JSON configuration.
      If specified, overrides any network definition in the JSON.

  .PARAMETER ConfigPath
      Path to JSON configuration file containing subnet requirements and optionally the network definition.
      Cannot be used together with -JsonConfig parameter.

  .PARAMETER JsonConfig
      JSON string containing subnet requirements and optionally the network definition.
      Cannot be used together with -ConfigPath parameter.

  .EXAMPLE
      # Create a configuration file with network defined in JSON using host requirements
      $config = @{
        network = "192.168.1.0/24"
        subnets = @(
          @{ name = "Management"; description = "Network management and monitoring"; hosts = 15 }
          @{ name = "Production"; description = "Production web servers"; hosts = 50 }
          @{ name = "Database"; description = "Database cluster"; hosts = 8 }
          @{ name = "Backup"; description = "Backup and storage network"; hosts = 5 }
        )
      }
      $config | ConvertTo-Json -Depth 3 | Out-File "network-config.json"
      
      # Simple usage - network defined in JSON
      New-SubnetPlanFromConfig -ConfigPath "network-config.json"
      
      # Or override the network from JSON
      New-SubnetPlanFromConfig -Network "10.0.1.0/24" -ConfigPath "network-config.json"

  .EXAMPLE
      # Create a configuration file using CIDR prefixes directly
      $configWithCidr = @{
        network = "192.168.1.0/24"
        subnets = @(
          @{ name = "csu-edge-transport-compute"; vlan = "203"; cidr = "28" }
          @{ name = "csu-exchange-compute"; vlan = "102"; cidr = "27" }
          @{ name = "msu-compute"; vlan = "302"; cidr = "27" }
          @{ name = "csu-edge-transport-management"; vlan = "101"; cidr = "27" }
        )
      }
      $configWithCidr | ConvertTo-Json -Depth 3 | Out-File "network-cidr-config.json"
      
      New-SubnetPlanFromConfig -ConfigPath "network-cidr-config.json"

  .EXAMPLE
      # Use direct JSON string with network definition
      $jsonConfig = @'
      {
        "network": "10.0.0.0/22",
        "subnets": [
          { "name": "MGMT", "description": "Management VLAN", "hosts": 25 },
          { "name": "DMZ", "description": "Demilitarized Zone", "hosts": 10 },
          { "name": "LAN", "description": "Internal LAN", "hosts": 100 }
        ]
      }
      '@
      
      New-SubnetPlanFromConfig -JsonConfig $jsonConfig

  .EXAMPLE
      # Dynamic columns with network in JSON
      $customConfig = @'
      {
        "network": "192.168.1.0/24",
        "subnets": [
          { "name": "DMZ", "vlan": 10, "zone": "External", "hosts": 30 },
          { "name": "LAN", "vlan": 20, "zone": "Internal", "hosts": 50 }
        ]
      }
      '@
      
      New-SubnetPlanFromConfig -JsonConfig $customConfig
      # Output will show: name, vlan, zone, Subnet, Prefix, Network, etc.

  .EXAMPLE
      # Export subnet plan to JSON for API integration or automation
      New-SubnetPlanFromConfig -ConfigPath "network-config.json" -AsJson | Out-File "subnet-plan.json"
      
      Generates the subnet plan and exports it to a JSON file for consumption by other tools.

  .EXAMPLE
      # Multiple primary networks - process independently and combine results
      $multiNetworkConfig = @'
      [
        {
          "network": "192.168.1.0/24",
          "subnets": [
            { "name": "csu-edge-transport-compute", "vlan": "203", "cidr": "28" },
            { "name": "csu-exchange-compute", "vlan": "102", "cidr": "27" }
          ]
        },
        {
          "network": "10.50.1.0/24",
          "subnets": [
            { "name": "csu-edge-transport-management", "vlan": "101", "cidr": "27" },
            { "name": "msu-management", "vlan": "201", "cidr": "27" }
          ]
        }
      ]
      '@
      
      New-SubnetPlanFromConfig -JsonConfig $multiNetworkConfig
      
      Processes multiple independent primary networks, each with their own subnets.
      Results are combined into a single output showing all networks and subnets.

  .OUTPUTS
      Array of PSCustomObject with dynamic properties based on JSON configuration.
      Columns automatically adjust to show all properties defined in the JSON.
      When -AsJson is specified, returns JSON string representation of the data.

  .NOTES
      JSON Configuration Format (Dynamic Properties):
      
      Single Network (Backward Compatible):
      {
        "network": "192.168.1.0/24",      // Optional - parent network in CIDR format
        "subnets": [
          {
            "name": "SubnetName",           // Required - subnet identifier
            "hosts": 15,                    // Option 1: Number of hosts needed (auto-calculates prefix)
            "cidr": "27",                   // Option 2: Direct CIDR prefix specification (use only one)
            "vlan": 101,                    // Optional - any custom property
            "zone": "DMZ",                  // Optional - any custom property
            "environment": "Production",    // Optional - any custom property
            "priority": "High",             // Optional - any custom property
            "description": "Purpose",       // Optional - subnet description
            "IPAssignments": [               // Optional - named host allocations within the subnet
              { "Name": "Gateway", "Position": 1 },
              { "Name": "VIP", "Position": 4 }
            ]
            // Add any custom properties you need!
          }
        ]
      }
      
      Multiple Networks (New):
      [
        {
          "network": "192.168.1.0/24",
          "subnets": [ /* subnets for first network */ ]
        },
        {
          "network": "10.50.1.0/24",
          "subnets": [ /* subnets for second network */ ]
        }
      ]
      
      Subnet Size Specification (choose one):
      - hosts: Specify number of hosts needed, function calculates optimal subnet size
      - cidr: Directly specify the CIDR prefix length (e.g., "28" for /28 subnet)
      
      Network Definition Priority:
      1. -Network parameter (if specified) - highest priority
      2. "network" field in JSON configuration
      3. Error if neither is provided
      
      The output table will automatically show all properties defined in your JSON.
      Only 'hosts', 'cidr', and 'network' are reserved - all other properties become display columns.
      
      IP Assignments:
      - Positions start at 1 and map to the nth usable host address (excluding network and broadcast).
      - Named assignments are rendered alongside unused address ranges for clear documentation.
      
      Multiple Networks:
      - Each network is processed independently with its own subnet allocations
      - Results from all networks are combined in the output
      - Custom properties are merged across all networks for consistent column display
  #>
  [CmdletBinding(DefaultParameterSetName = 'FilePath')]
  param(
    [Parameter(Mandatory = $false, HelpMessage = "IPv4 network in CIDR format (e.g., '192.168.1.0/24'). Optional if network is defined in JSON.")]
    [ValidatePattern('^(\d{1,3}\.){3}\d{1,3}/([1-2]?[0-9]|3[0-2])$')]
    [string]$Network,
    
    [Parameter(Mandatory = $true, ParameterSetName = 'FilePath', HelpMessage = "Path to JSON configuration file")]
    [ValidateScript({ Test-Path $_ })]
    [string]$ConfigPath,
    
    [Parameter(Mandatory = $true, ParameterSetName = 'JsonString', HelpMessage = "JSON configuration string")]
    [ValidateNotNullOrEmpty()]
    [string]$JsonConfig,
    
    [Parameter(Mandatory = $false, HelpMessage = "Output results in JSON format")]
    [switch]$AsJson,

    [Parameter(Mandatory = $false, HelpMessage = "Write detailed plan to a JSON file")]
    [ValidateNotNullOrEmpty()]
    [string]$ExportJsonPath,

    [Parameter(Mandatory = $false, HelpMessage = "Write detailed plan to a CSV file")]
    [ValidateNotNullOrEmpty()]
    [string]$ExportCsvPath,

    [Parameter(Mandatory = $false, HelpMessage = "Write detailed plan to a Markdown table file")]
    [ValidateNotNullOrEmpty()]
    [string]$ExportMarkdownPath
  )

  Write-Verbose "Starting subnet planning from JSON configuration"

  function Convert-IpToIntLocal {
    param([Parameter(Mandatory = $true)][string]$Ip)

    $octets = $Ip.Split('.') | ForEach-Object { [uint32]$_ }
    return ($octets[0] -shl 24) -bor ($octets[1] -shl 16) -bor ($octets[2] -shl 8) -bor $octets[3]
  }

  function Convert-IntToIpLocal {
    param([Parameter(Mandatory = $true)][uint32]$Value)

    $octet1 = ($Value -shr 24) -band 0xFF
    $octet2 = ($Value -shr 16) -band 0xFF
    $octet3 = ($Value -shr 8) -band 0xFF
    $octet4 = $Value -band 0xFF
    return "$octet1.$octet2.$octet3.$octet4"
  }

  function Get-NetmaskFromPrefix {
    param([Parameter(Mandatory = $true)][int]$Prefix)

    if ($Prefix -lt 0 -or $Prefix -gt 32) {
      throw "Invalid prefix length: $Prefix. Must be between 0 and 32."
    }

    if ($Prefix -eq 0) {
      return "0.0.0.0"
    }

    $maskValue = [uint32]([Math]::Pow(2, 32) - [Math]::Pow(2, 32 - $Prefix))
    return Convert-IntToIpLocal -Value $maskValue
  }

  function Get-PrefixForHostCountLocal {
    param([Parameter(Mandatory = $true)][int]$HostCount)

    if ($HostCount -le 0) {
      throw "Host count must be greater than 0"
    }

    $requiredAddresses = $HostCount + 2
    $bitsNeeded = [Math]::Ceiling([Math]::Log($requiredAddresses, 2))
    $prefix = 32 - $bitsNeeded

    if ($prefix -lt 1) { $prefix = 1 }
    if ($prefix -gt 30) { $prefix = 30 }

    return $prefix
  }

  function Get-StandardizedPropertyName {
    param([string]$PropertyName)

    if (-not $PropertyName -or $PropertyName.Length -eq 0) {
      return $PropertyName
    }

    if ($PropertyName.Length -eq 1) {
      return $PropertyName.ToUpper()
    }

    $first = $PropertyName.Substring(0, 1).ToUpper()
    $rest = $PropertyName.Substring(1)
    return "$first$rest"
  }

  function Set-DirectoryForPath {
    param([Parameter(Mandatory = $true)][string]$Path)

    $directory = Split-Path -Path $Path -Parent
    if ($directory -and -not (Test-Path -Path $directory)) {
      New-Item -ItemType Directory -Path $directory -Force | Out-Null
    }
  }

  function ConvertTo-MarkdownTable {
    param(
      [Parameter(Mandatory = $true)][System.Collections.IEnumerable]$Rows,
      [Parameter(Mandatory = $true)][string[]]$Columns
    )

    if (-not $Rows) {
      return ''
    }

    $header = ($Columns -join ' | ')
    $separator = ($Columns | ForEach-Object { '---' }) -join ' | '
    $lines = New-Object System.Collections.Generic.List[string]
    $lines.Add($header) | Out-Null
    $lines.Add($separator) | Out-Null

    foreach ($row in $Rows) {
      $values = foreach ($column in $Columns) {
        $property = $row.PSObject.Properties[$column]
        $value = if ($property) { [string]$property.Value } else { '' }
        #$value = $value -replace '\|', '|'
        $value = $value -replace "`r`n", '<br />'
        $value = $value -replace "`n", '<br />'
        $value
      }
      $lines.Add(($values -join ' | ')) | Out-Null
    }

    return [string]::Join([Environment]::NewLine, $lines)
  }

  try {
    if ($PSCmdlet.ParameterSetName -eq 'FilePath') {
      Write-Verbose "Loading configuration from file: $ConfigPath"
      $jsonContent = Get-Content -Path $ConfigPath -Raw
    }
    else {
      Write-Verbose "Using provided JSON configuration string"
      $jsonContent = $JsonConfig
    }

    $config = $jsonContent | ConvertFrom-Json

    # Detect if input is array of networks or single network object
    $networkConfigs = @()
    if ($config -is [System.Array]) {
      Write-Verbose "Detected array of network configurations"
      $networkConfigs = $config
    }
    else {
      Write-Verbose "Detected single network configuration"
      $networkConfigs = @($config)
    }

    # Validate we have at least one network configuration
    if ($networkConfigs.Count -eq 0) {
      throw "Invalid configuration: No network configurations found"
    }

    # Process each network configuration
    $allDetailedRows = @()
    $allCustomPropertyOrder = @()
    $allCustomPropertyLookup = @{}

    foreach ($networkConfig in $networkConfigs) {
      $networkToUse = $null
      if ($Network) {
        $networkToUse = $Network
        Write-Verbose "Using network from parameter: $Network"
      }
      elseif ($networkConfig.network) {
        $networkToUse = $networkConfig.network
        Write-Verbose "Using network from JSON configuration: $($networkConfig.network)"

        if ($networkToUse -notmatch '^(\d{1,3}\.){3}\d{1,3}/([1-2]?[0-9]|3[0-2])$') {
          throw "Invalid network format in JSON configuration: '$networkToUse'. Expected format: 'x.x.x.x/yy'"
        }
      }
      else {
        throw "Network must be specified either as parameter -Network or in JSON configuration as 'network' field"
      }

      Write-Verbose "Network to use for subnet planning: $networkToUse"

      if (-not $networkConfig.subnets) {
        throw "Invalid configuration: Missing 'subnets' array in JSON for network '$networkToUse'"
      }

      if ($networkConfig.subnets.Count -eq 0) {
        throw "Invalid configuration: 'subnets' array is empty for network '$networkToUse'"
      }

      $prefixCounts = @{}
      $metadataList = @()
      $reservedProperties = @('network', 'hosts', 'cidr', 'ipassignments', 'subnet', 'prefix', 'broadcast', 'firsthost', 'endhost', 'usablehosts', 'category', 'mask', 'label', 'ip', 'range', 'ipvalue')
      $customPropertyLookup = @{}
      $customPropertyOrder = @()
      $subnetIndex = 0

      Write-Verbose "Processing subnet configuration for network ${networkToUse}"

      foreach ($subnet in $networkConfig.subnets) {
        $subnetIndex++

        if (-not $subnet.name) {
          throw "Invalid configuration: Subnet missing 'name' field"
        }

        $hasHosts = $null -ne $subnet.hosts
        $hasCidr = $null -ne $subnet.cidr

        if (-not $hasHosts -and -not $hasCidr) {
          throw "Invalid configuration: Subnet '$($subnet.name)' missing both 'hosts' and 'cidr' fields. One of them is required."
        }

        if ($hasHosts -and $hasCidr) {
          throw "Invalid configuration: Subnet '$($subnet.name)' has both 'hosts' and 'cidr' fields. Only one should be specified."
        }

        $prefixLength = $null
        $hostCapacity = 0

        if ($hasHosts) {
          $hostCount = [int]$subnet.hosts
          $prefixLength = Get-PrefixForHostCountLocal -HostCount $hostCount
          $totalAddresses = [uint32][Math]::Pow(2, 32 - $prefixLength)
          $hostCapacity = if ($totalAddresses -gt 2) { [int]($totalAddresses - 2) } else { 0 }
          $descriptionText = "$hostCount hosts → /$prefixLength ($hostCapacity usable hosts)"
        }
        else {
          $prefixLength = [int]$subnet.cidr

          if ($prefixLength -lt 1 -or $prefixLength -gt 32) {
            throw "Invalid configuration: Subnet '$($subnet.name)' has invalid CIDR prefix '$prefixLength'. Must be between 1 and 32."
          }

          $totalAddresses = [uint32][Math]::Pow(2, 32 - $prefixLength)
          $hostCapacity = if ($totalAddresses -gt 2) { [int]($totalAddresses - 2) } else { 0 }
          $descriptionText = "/$prefixLength ($hostCapacity usable hosts)"
        }

        Write-Verbose "  $($subnet.name): $descriptionText"

        if ($prefixCounts.ContainsKey($prefixLength)) {
          $prefixCounts[$prefixLength] += 1
        }
        else {
          $prefixCounts[$prefixLength] = 1
        }

        $assignmentList = @()
        if ($subnet.IPAssignments) {
          if (-not ($subnet.IPAssignments -is [System.Collections.IEnumerable])) {
            throw "Invalid configuration: 'IPAssignments' for subnet '$($subnet.name)' must be an array"
          }

          $positionTracker = @{}
          foreach ($assignment in $subnet.IPAssignments) {
            if (-not $assignment.Name) {
              throw "Invalid configuration: IPAssignments entry for subnet '$($subnet.name)' missing 'Name'"
            }

            if ($null -eq $assignment.Position) {
              throw "Invalid configuration: IPAssignments entry '$($assignment.Name)' for subnet '$($subnet.name)' missing 'Position'"
            }

            $position = [int]$assignment.Position
            if ($position -lt 1) {
              throw "Invalid configuration: IPAssignments position for '$($assignment.Name)' in subnet '$($subnet.name)' must be 1 or greater"
            }

            if ($positionTracker.ContainsKey($position)) {
              throw "Invalid configuration: Duplicate IPAssignments position '$position' found in subnet '$($subnet.name)'"
            }

            $positionTracker[$position] = $true
            $assignmentList += [PSCustomObject]@{
              Name     = [string]$assignment.Name
              Position = $position
            }
          }
        }

        foreach ($property in $subnet.PSObject.Properties.Name) {
          $lower = $property.ToLower()
          if ($reservedProperties -contains $lower) {
            continue
          }

          if (-not $customPropertyLookup.ContainsKey($lower)) {
            $customPropertyLookup[$lower] = $property
            $customPropertyOrder += $property
          }
        }

        $metadataList += [PSCustomObject]@{
          OriginalIndex = $subnetIndex
          Config        = $subnet
          Prefix        = $prefixLength
          HostCapacity  = $hostCapacity
          IpAssignments = $assignmentList
          Assigned      = $false
        }
      }

      Write-Verbose "Calling New-SubnetPlan with processed requirements"

      $rawResults = New-SubnetPlan -Network $networkToUse -PrefixRequirements $prefixCounts

      $addRow = {
        param(
          $Metadata,
          [string]$Label,
          [uint32]$StartInt,
          $EndInt,
          [string]$Category,
          [string]$SubnetValue,
          [int]$Prefix,
          [string]$Mask
        )

        $properties = [ordered]@{}
        foreach ($property in $customPropertyOrder) {
          $standardizedName = Get-StandardizedPropertyName -PropertyName $property
          $value = ''

          if ($Metadata -and $Metadata.Config.PSObject.Properties[$property]) {
            $value = $Metadata.Config.PSObject.Properties[$property].Value
            if ($null -eq $value) { $value = '' }
          }
          elseif (-not $Metadata -and $standardizedName -eq 'Name') {
            $value = 'Available'
          }

          $properties[$standardizedName] = $value
        }

        $properties['Subnet'] = $SubnetValue
        $properties['Label'] = $Label

        $startIp = Convert-IntToIpLocal -Value $StartInt
        if ($null -ne $EndInt) {
          $endIp = Convert-IntToIpLocal -Value ([uint32]$EndInt)
          $properties['IP'] = "$startIp - $endIp"
        }
        else {
          $properties['IP'] = $startIp
        }

        $ipCount = if ($null -ne $EndInt) {
          ([uint64]$EndInt) - [uint64]$StartInt + 1
        }
        else {
          [uint64]1
        }

        $properties['Prefix'] = "/$Prefix"
        $properties['Mask'] = $Mask
        $properties['Category'] = $Category
        $properties['TotalIPs'] = $ipCount

        return [PSCustomObject]$properties
      }

      $detailedRows = @()

      foreach ($result in $rawResults) {
        if ($result.Name -eq 'Assigned') {
          $matching = $metadataList |
            Where-Object { -not $_.Assigned -and $_.Prefix -eq $result.Prefix } |
            Sort-Object OriginalIndex |
            Select-Object -First 1

          if (-not $matching) {
            throw "Unable to match allocated subnet '/$($result.Prefix)' to configuration entry"
          }

          $matching.Assigned = $true

          $mask = Get-NetmaskFromPrefix -Prefix $result.Prefix
          $networkInt = Convert-IpToIntLocal -Ip $result.Network
          $broadcastInt = Convert-IpToIntLocal -Ip $result.Broadcast
          $usableHostsCount = [int]$result.UsableHosts

          $sortedAssignments = $matching.IpAssignments | Sort-Object Position, Name

          if ($usableHostsCount -eq 0 -and $sortedAssignments.Count -gt 0) {
            throw "Subnet '$($matching.Config.name)' does not support host assignments but IPAssignments were provided"
          }

          foreach ($assignment in $sortedAssignments) {
            if ($assignment.Position -gt $usableHostsCount) {
              throw "IPAssignments position '$($assignment.Position)' for subnet '$($matching.Config.name)' exceeds available host addresses ($usableHostsCount)"
            }
          }

          $detailedRows += & $addRow $matching 'Network' $networkInt $null 'Network' $result.Subnet $result.Prefix $mask

          $previousPosition = 0
          foreach ($assignment in $sortedAssignments) {
            if ($usableHostsCount -gt 0) {
              $gapStart = $previousPosition + 1
              $gapEnd = $assignment.Position - 1
              if ($gapStart -le $gapEnd) {
                $startInt = $networkInt + [uint32]$gapStart
                $endInt = $networkInt + [uint32]$gapEnd
                $detailedRows += & $addRow $matching 'Unused Range' $startInt $endInt 'Unused' $result.Subnet $result.Prefix $mask
              }
            }

            $ipInt = $networkInt + [uint32]$assignment.Position
            $detailedRows += & $addRow $matching $assignment.Name $ipInt $null 'Assignment' $result.Subnet $result.Prefix $mask
            $previousPosition = $assignment.Position
          }

          if ($usableHostsCount -gt 0) {
            $finalGapStart = $previousPosition + 1
            if ($finalGapStart -le $usableHostsCount) {
              $startInt = $networkInt + [uint32]$finalGapStart
              $endInt = $networkInt + [uint32]$usableHostsCount
              $detailedRows += & $addRow $matching 'Unused Range' $startInt $endInt 'Unused' $result.Subnet $result.Prefix $mask
            }
          }

          $detailedRows += & $addRow $matching 'Broadcast' $broadcastInt $null 'Broadcast' $result.Subnet $result.Prefix $mask
        }
        else {
          $mask = Get-NetmaskFromPrefix -Prefix $result.Prefix
          $startInt = if ($result.FirstHost) { Convert-IpToIntLocal -Ip $result.FirstHost } else { Convert-IpToIntLocal -Ip $result.Network }
          $endInt = if ($result.EndHost) { Convert-IpToIntLocal -Ip $result.EndHost } else { Convert-IpToIntLocal -Ip $result.Broadcast }

          $detailedRows += & $addRow $null 'Available Range' $startInt $endInt $result.Name $result.Subnet $result.Prefix $mask
        }
      }

      $unmatched = $metadataList | Where-Object { -not $_.Assigned }
      if ($unmatched.Count -gt 0) {
        $unmatchedNames = $unmatched | ForEach-Object { $_.Config.name }
        throw "Configured subnets could not be matched to allocation results: $($unmatchedNames -join ', ')"
      }

      # Store custom properties for merging across networks
      foreach ($property in $customPropertyOrder) {
        $lower = $property.ToLower()
        if (-not $allCustomPropertyLookup.ContainsKey($lower)) {
          $allCustomPropertyLookup[$lower] = $property
          $allCustomPropertyOrder += $property
        }
      }

      # Add rows from this network to the combined results
      $allDetailedRows += $detailedRows
    }

    $displayProperties = @('Subnet')
    foreach ($property in $allCustomPropertyOrder) {
      $standardized = Get-StandardizedPropertyName -PropertyName $property
      if ($displayProperties -notcontains $standardized) {
        $displayProperties += $standardized
      }
    }

    foreach ($defaultProperty in @('Label', 'IP', 'TotalIPs', 'Prefix', 'Mask', 'Category')) {
      if ($displayProperties -notcontains $defaultProperty) {
        $displayProperties += $defaultProperty
      }
    }

    $summaryView = $allDetailedRows | Select-Object $displayProperties

    $jsonOutput = $null
    if ($AsJson -or $ExportJsonPath) {
      $jsonOutput = $summaryView | ConvertTo-Json -Depth 10
    }

    if ($ExportJsonPath) {
      Set-DirectoryForPath -Path $ExportJsonPath
      Set-Content -Path $ExportJsonPath -Value $jsonOutput -Encoding UTF8
    }

    if ($ExportCsvPath) {
      Set-DirectoryForPath -Path $ExportCsvPath
      $summaryView | Export-Csv -Path $ExportCsvPath -NoTypeInformation -Encoding UTF8
    }

    if ($ExportMarkdownPath) {
      Set-DirectoryForPath -Path $ExportMarkdownPath
      $markdownOutput = ConvertTo-MarkdownTable -Rows $summaryView -Columns $displayProperties
      Set-Content -Path $ExportMarkdownPath -Value $markdownOutput -Encoding UTF8
    }

    if ($AsJson) {
      return $jsonOutput
    }

    if ($MyInvocation.Line -match '\$\w+\s*=' -or $MyInvocation.Line -match '\|') {
      return $summaryView
    }
    else {
      $tableOutput = $summaryView | Format-Table -AutoSize | Out-String
      Write-Host $tableOutput.TrimEnd()
      Write-Host ""
    }
  }
  catch {
    Write-Error "Failed to process subnet configuration: $($_.Exception.Message)"
    throw
  }
}
