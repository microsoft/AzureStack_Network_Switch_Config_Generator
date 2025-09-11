#Requires -Version 5.1

<#
.SYNOPSIS
    PortMap Tool - Network Device Physical Port Assignment Documentation Generator

.DESCRIPTION
    The PortMap tool is a network documentation utility designed to detail physical port assignments 
    and cable connections for network devices. It creates comprehensive documentation artifacts that 
    describe the physical cabling configuration of network switches.

.PARAMETER InputFile
    Path to the input JSON configuration file describing network devices and connections.

.PARAMETER OutputFormat
    Output format for the port mapping data. Valid values: Markdown, CSV, JSON

.PARAMETER OutputFile
    Optional output file path. If not specified, Markdown output goes to console.

.PARAMETER ShowUnused
    Include unused ports in the output analysis.

.PARAMETER DeviceFilter
    Filter output by specific device names (comma-separated list).

.PARAMETER Validate
    Validate the input JSON configuration without generating output.

.EXAMPLE
    .\PortMap.ps1 -InputFile "network-config.json" -OutputFormat Markdown

.EXAMPLE
    .\PortMap.ps1 -InputFile "datacenter-design.json" -OutputFormat CSV -OutputFile "port-mapping.csv"

.EXAMPLE
    .\PortMap.ps1 -InputFile "rack-layout.json" -OutputFormat JSON -OutputFile "port-data.json" -ShowUnused

.NOTES
    Version: 1.0
    Author: Network Engineering Team
    Purpose: Part of AzureStack Network Switch Config Generator project
#>

[CmdletBinding(DefaultParameterSetName = 'Process')]
param(
    [Parameter(
        Mandatory = $true,
        Position = 0,
        HelpMessage = "Path to the input JSON configuration file describing network devices and connections"
    )]
    [ValidateScript({
            if (-not (Test-Path -Path $_ -PathType Leaf)) {
                throw "File does not exist: $_"
            }
            if (-not ($_ -match '\.json$')) {
                throw "File must have .json extension: $_"
            }
            return $true
        })]
    [Alias('Config', 'Input')]
    [string]$InputFile,

    [Parameter(
        Mandatory = $true,
        Position = 1,
        HelpMessage = "Output format for the port mapping data"
    )]
    [ValidateSet("Markdown", "CSV", "JSON", IgnoreCase = $true)]
    [Alias('Format')]
    [string]$OutputFormat,

    [Parameter(
        Mandatory = $false,
        HelpMessage = "Optional output file path. If not specified, output goes to console"
    )]
    [ValidateScript({
            $directory = Split-Path -Path $_ -Parent
            if ($directory -and -not (Test-Path -Path $directory)) {
                try {
                    New-Item -Path $directory -ItemType Directory -Force -WhatIf | Out-Null
                    return $true
                }
                catch {
                    throw "Cannot create output directory: $directory"
                }
            }
            return $true
        })]
    [Alias('Output', 'File')]
    [string]$OutputFile,

    [Parameter(
        Mandatory = $false,
        HelpMessage = "Include unused ports in the output analysis"
    )]
    [Alias('Unused')]
    [switch]$ShowUnused,

    [Parameter(
        Mandatory = $false,
        HelpMessage = "Filter output by specific device names (comma-separated list)"
    )]
    [ValidateNotNullOrEmpty()]
    [Alias('Filter', 'Devices')]
    [string[]]$DeviceFilter,

    [Parameter(
        Mandatory = $false,
        ParameterSetName = 'Validate',
        HelpMessage = "Validate the input JSON configuration without generating output"
    )]
    [switch]$Validate,

    [Parameter(
        Mandatory = $false,
        HelpMessage = "Display detailed progress information"
    )]
    [switch]$Detailed
)

# Script-level variables - Using proper scoping
$Script:ErrorCount = 0
$Script:ValidationErrors = [System.Collections.Generic.List[string]]::new()
$Script:PortMapVersion = "1.0.0"

# Set strict mode for better error handling
Set-StrictMode -Version Latest

# Initialize error action preference
$ErrorActionPreference = 'Stop'

#region Helper Functions

function Write-PortMapLog {
    <#
    .SYNOPSIS
        Writes formatted log messages for the PortMap tool.
    
    .DESCRIPTION
        Provides consistent logging with timestamps and color coding based on message level.
    
    .PARAMETER Message
        The message to log.
    
    .PARAMETER Level
        The logging level (Info, Warning, Error).
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [ValidateSet("Info", "Warning", "Error")]
        [string]$Level = "Info"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    
    switch ($Level) {
        "Info" { 
            Write-Host $logMessage -ForegroundColor Green 
        }
        "Warning" { 
            Write-Warning $logMessage 
        }
        "Error" { 
            Write-Error $logMessage
            $Script:ErrorCount++
        }
    }
}

function Test-NetworkConfiguration {
    <#
    .SYNOPSIS
        Validates the JSON network configuration structure.
    
    .DESCRIPTION
        Performs comprehensive validation of the input JSON configuration to ensure
        all required fields are present and properly formatted.
    
    .PARAMETER Configuration
        The configuration object to validate.
    
    .OUTPUTS
        System.Boolean
        Returns $true if configuration is valid, $false otherwise.
    #>
    [CmdletBinding()]
    [OutputType([bool])]
    param(
        [Parameter(Mandatory = $true)]
        [object]$Configuration
    )
    
    Write-PortMapLog "Validating JSON configuration..." -Level Info
    
    # Validate devices section
    if (-not $Configuration.devices) {
        $Script:ValidationErrors.Add("Missing 'devices' section in configuration")
        return $false
    }
    
    foreach ($device in $Configuration.devices) {
        if (-not $device.deviceName) {
            $Script:ValidationErrors.Add("Device missing 'deviceName' property")
        }
        if (-not $device.deviceMake) {
            $Script:ValidationErrors.Add("Device '$($device.deviceName)' missing 'deviceMake' property")
        }
        if (-not $device.deviceModel) {
            $Script:ValidationErrors.Add("Device '$($device.deviceName)' missing 'deviceModel' property")
        }
        if (-not $device.portRanges) {
            $Script:ValidationErrors.Add("Device '$($device.deviceName)' missing 'portRanges' property")
        }
        
        # Validate port ranges structure
        if ($device.portRanges) {
            foreach ($portRange in $device.portRanges) {
                if (-not $portRange.range) {
                    $Script:ValidationErrors.Add("Device '$($device.deviceName)' has port range missing 'range' property")
                }
                if (-not $portRange.mediaType) {
                    $Script:ValidationErrors.Add("Device '$($device.deviceName)' has port range missing 'mediaType' property")
                }
            }
        }
    }
    
    # Validate connections section
    if ($Configuration.connections) {
        foreach ($connection in $Configuration.connections) {
            if (-not $connection.sourceDevice) {
                $Script:ValidationErrors.Add("Connection missing 'sourceDevice' property")
            }
            if (-not $connection.destinationDevice) {
                $Script:ValidationErrors.Add("Connection missing 'destinationDevice' property")
            }
            if (-not $connection.sourcePorts) {
                $Script:ValidationErrors.Add("Connection missing 'sourcePorts' property")
            }
        }
    }
    
    if ($Script:ValidationErrors.Count -gt 0) {
        foreach ($validationError in $Script:ValidationErrors) {
            Write-PortMapLog $validationError -Level Error
        }
        return $false
    }
    
    Write-PortMapLog "Configuration validation passed" -Level Info
    return $true
}

function Expand-NetworkPortRange {
    <#
    .SYNOPSIS
        Expands port range notation into individual port numbers.
    
    .DESCRIPTION
        Converts port range strings like "1-48" into arrays of individual port numbers.
        Also handles single port numbers.
    
    .PARAMETER Range
        The port range string to expand (e.g., "1-48" or "25").
    
    .OUTPUTS
        System.Int32[]
        Array of individual port numbers.
    #>
    [CmdletBinding()]
    [OutputType([int[]])]
    param(
        [Parameter(Mandatory = $true)]
        [string]$Range
    )
    
    if ($Range -match '^(\d+)-(\d+)$') {
        $start = [int]$matches[1]
        $end = [int]$matches[2]
        
        if ($start -gt $end) {
            Write-PortMapLog "Invalid port range: start ($start) is greater than end ($end)" -Level Warning
            return @()
        }
        
        return $start..$end
    }
    elseif ($Range -match '^\d+$') {
        return @([int]$Range)
    }
    else {
        Write-PortMapLog "Invalid port range format: '$Range'. Expected format: '0-48' or '25'. Decimal values like '3.0' are not supported." -Level Warning
        return @()
    }
}

function New-UniqueOutputFileName {
    <#
    .SYNOPSIS
        Generates a unique output filename with device information and incremental numbering.
    
    .DESCRIPTION
        Creates a filename that includes input parameters, device information, and ensures
        uniqueness by adding incremental numbers if the file already exists.
    
    .PARAMETER InputFile
        The original input file path.
    
    .PARAMETER OutputFormat
        The desired output format (JSON, CSV, Markdown).
    
    .PARAMETER Devices
        Array of device objects containing deviceName, deviceMake, and deviceModel.
    
    .PARAMETER OutputFile
        Optional explicit output file path.
    
    .OUTPUTS
        System.String
        Returns the unique output file path.
    #>
    [CmdletBinding()]
    [OutputType([string])]
    param(
        [Parameter(Mandatory = $true)]
        [string]$InputFile,
        
        [Parameter(Mandatory = $true)]
        [string]$OutputFormat,
        
        [Parameter(Mandatory = $true)]
        [array]$Devices,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFile
    )
    
    if ($OutputFile) {
        # If explicit output file is provided, ensure it has a full path
        if ([System.IO.Path]::IsPathRooted($OutputFile)) {
            $baseFileName = $OutputFile
        }
        else {
            $baseFileName = Join-Path -Path (Get-Location) -ChildPath $OutputFile
        }
    }
    else {
        # Generate filename with device information
        $inputBaseName = [System.IO.Path]::GetFileNameWithoutExtension($InputFile)
        $extension = switch ($OutputFormat.ToLower()) {
            "csv" { "csv" }
            "json" { "json" }
            default { "md" }
        }
        
        # Extract unique device makes and models
        $deviceMakes = $Devices | Select-Object -ExpandProperty deviceMake -Unique | Sort-Object
        $deviceModels = $Devices | Select-Object -ExpandProperty deviceModel -Unique | Sort-Object
        $deviceCount = $Devices.Count
        
        # Create a compact device identifier
        $makeString = ($deviceMakes -join "-").Replace(" ", "")
        $modelString = ($deviceModels -join "-").Replace(" ", "").Replace("/", "-")
        
        # For CSV format with single device, include device name in filename
        if ($OutputFormat.ToLower() -eq "csv" -and $deviceCount -eq 1) {
            $deviceName = $Devices[0].deviceName
            $sanitizedDeviceName = $deviceName -replace '[^\w\-.]', '_'
            $deviceInfo = "$sanitizedDeviceName-$makeString-$modelString"
            $baseFileName = Join-Path -Path (Get-Location) -ChildPath "$inputBaseName-$deviceInfo-portmap.$extension"
        }
        else {
            # Generate base filename with device info (for multi-device or non-CSV formats)
            $deviceInfo = "$makeString-$modelString-$($deviceCount)dev"
            $baseFileName = Join-Path -Path (Get-Location) -ChildPath "$inputBaseName-$deviceInfo-portmap.$extension"
        }
    }
    
    # Check if file exists and create unique name if needed
    $finalFileName = $baseFileName
    $counter = 1
    
    while (Test-Path -Path $finalFileName) {
        $directory = Split-Path -Path $baseFileName -Parent
        if (-not $directory) {
            $directory = Get-Location
        }
        $nameWithoutExt = [System.IO.Path]::GetFileNameWithoutExtension($baseFileName)
        $extension = [System.IO.Path]::GetExtension($baseFileName)
        
        # Remove previous counter if it exists
        if ($nameWithoutExt -match '(.+)-(\d+)$') {
            $nameWithoutExt = $matches[1]
        }
        
        $finalFileName = Join-Path -Path $directory -ChildPath "$nameWithoutExt-$counter$extension"
        $counter++
    }
    
    return $finalFileName
}

function Get-NetworkDevicePortInfo {
    <#
    .SYNOPSIS
        Extracts port information from a device configuration.
    
    .DESCRIPTION
        Processes device configuration to create a comprehensive port information structure
        including total ports, port details, and usage tracking.
    
    .PARAMETER Device
        The device configuration object.
    
    .OUTPUTS
        System.Collections.Hashtable
        Hashtable containing port information and details.
    #>
    [CmdletBinding()]
    [OutputType([hashtable])]
    param(
        [Parameter(Mandatory = $true)]
        [object]$Device
    )
    
    $portInfo = @{
        TotalPorts  = 0
        UsedPorts   = [System.Collections.Generic.List[int]]::new()
        UnusedPorts = [System.Collections.Generic.List[int]]::new()
        PortDetails = @{}
    }
    
    # Calculate total ports and create port details
    foreach ($portRange in $Device.portRanges) {
        $ports = Expand-NetworkPortRange -Range $portRange.range
        $portInfo.TotalPorts += @($ports).Count
        
        foreach ($port in $ports) {
            $portInfo.PortDetails[$port] = @{
                MediaType   = $portRange.mediaType
                Speed       = $portRange.speed ?? "Unknown"
                Description = $portRange.description ?? "No description"
                IsUsed      = $false
                Connection  = $null
            }
        }
    }
    
    return $portInfo
}

function New-ConnectionMappings {
    <#
    .SYNOPSIS
        Creates connection mappings from configuration data.
    
    .DESCRIPTION
        Processes connection configurations to create detailed port-to-port mappings
        and updates device port usage information.
    
    .PARAMETER Connections
        Array of connection configuration objects.
    
    .PARAMETER DevicePortInfo
        Hashtable of device port information to update.
    
    .OUTPUTS
        System.Object[]
        Array of connection mapping objects.
    #>
    [CmdletBinding()]
    [OutputType([object[]])]
    param(
        [Parameter(Mandatory = $true)]
        [AllowEmptyCollection()]
        [object[]]$Connections,
        
        [Parameter(Mandatory = $true)]
        [hashtable]$DevicePortInfo
    )
    
    $connectionMappings = [System.Collections.Generic.List[object]]::new()
    
    foreach ($connection in $Connections) {
        try {
            # Parse source ports with error context
            $sourcePorts = Expand-NetworkPortRange -Range $connection.sourcePorts
            if ($null -eq $sourcePorts -or @($sourcePorts).Count -eq 0) {
                throw "Failed to parse source ports: '$($connection.sourcePorts)'"
            }
            
            # Parse destination ports with error context
            $destPorts = if ($connection.destinationPorts) { 
                $destPortsResult = Expand-NetworkPortRange -Range $connection.destinationPorts
                # Allow empty results for destination ports (port 0 or invalid formats are handled gracefully)
                if ($null -eq $destPortsResult) {
                    Write-PortMapLog "Could not parse destination ports '$($connection.destinationPorts)' for connection from '$($connection.sourceDevice)' to '$($connection.destinationDevice)', using port 0 as placeholder" -Level Warning
                    @(0) * @($sourcePorts).Count
                }
                elseif (@($destPortsResult).Count -eq 0) {
                    Write-PortMapLog "Destination ports '$($connection.destinationPorts)' for connection from '$($connection.sourceDevice)' to '$($connection.destinationDevice)' resulted in empty range, using port 0 as placeholder" -Level Warning
                    @(0) * @($sourcePorts).Count
                }
                else {
                    $destPortsResult
                }
            }
            else { 
                @(0) * @($sourcePorts).Count 
            }
        }
        catch {
            $errorMsg = "Connection processing failed - SourceDevice: '$($connection.sourceDevice)', SourcePorts: '$($connection.sourcePorts)', DestinationDevice: '$($connection.destinationDevice)', DestinationPorts: '$($connection.destinationPorts)'. Error: $($_.Exception.Message)"
            Write-PortMapLog $errorMsg -Level Error
            throw $errorMsg
        }
        
        for ($i = 0; $i -lt @($sourcePorts).Count; $i++) {
            $sourcePort = $sourcePorts[$i]
            $destPort = if ($i -lt @($destPorts).Count) { $destPorts[$i] } else { $destPorts[0] }
            
            # Mark source port as used
            if ($DevicePortInfo.ContainsKey($connection.sourceDevice)) {
                if ($DevicePortInfo[$connection.sourceDevice].PortDetails.ContainsKey($sourcePort)) {
                    $DevicePortInfo[$connection.sourceDevice].PortDetails[$sourcePort].IsUsed = $true
                    $DevicePortInfo[$connection.sourceDevice].PortDetails[$sourcePort].Connection = $connection.destinationDevice
                    $DevicePortInfo[$connection.sourceDevice].UsedPorts.Add($sourcePort)
                }
            }
            
            # Create connection mapping
            $mapping = [PSCustomObject]@{
                SourceDevice      = $connection.sourceDevice
                SourcePort        = $sourcePort
                SourceMedia       = $connection.sourceMedia ?? "Unknown"
                DestinationDevice = $connection.destinationDevice
                DestinationPort   = $destPort
                DestinationMedia  = $connection.destinationMedia ?? "Unknown"
                Status            = "Active"
                ConnectionType    = $connection.connectionType ?? "Unknown"
                Notes             = $connection.notes ?? ""
            }
            
            $connectionMappings.Add($mapping)
        }
    }
    
    # Update unused ports
    foreach ($deviceName in $DevicePortInfo.Keys) {
        $device = $DevicePortInfo[$deviceName]
        $unusedPorts = $device.PortDetails.Keys | Where-Object { -not $device.PortDetails[$_].IsUsed }
        $device.UnusedPorts.Clear()
        foreach ($port in $unusedPorts) {
            $device.UnusedPorts.Add($port)
        }
    }
    
    return $connectionMappings.ToArray()
}

#endregion

#region Output Generators

function ConvertTo-MarkdownOutput {
    <#
    .SYNOPSIS
        Converts port mapping data to Markdown format.
    
    .DESCRIPTION
        Creates human-readable Markdown tables for network documentation including
        device summaries, connection mappings, and optionally unused ports.
    
    .PARAMETER Devices
        Array of device configuration objects.
    
    .PARAMETER DevicePortInfo
        Hashtable containing device port information.
    
    .PARAMETER ConnectionMappings
        Array of connection mapping objects.
    
    .PARAMETER InputFile
        Path to the input JSON configuration file.
    
    .PARAMETER OutputFormat
        The output format being used.
    
    .OUTPUTS
        System.String
        Formatted Markdown documentation.
    #>
    [CmdletBinding()]
    [OutputType([string])]
    param(
        [Parameter(Mandatory = $true)]
        [object[]]$Devices,
        
        [Parameter(Mandatory = $true)]
        [hashtable]$DevicePortInfo,
        
        [Parameter(Mandatory = $true)]
        [AllowEmptyCollection()]
        [object[]]$ConnectionMappings,
        
        [Parameter(Mandatory = $false)]
        [string]$InputFile,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFormat
    )
    
    $markdownBuilder = [System.Text.StringBuilder]::new()
    
    # Header with metadata
    [void]$markdownBuilder.AppendLine("# Network Port Mapping Documentation")
    [void]$markdownBuilder.AppendLine("")
    [void]$markdownBuilder.AppendLine("**Generated on:** $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')")
    [void]$markdownBuilder.AppendLine("**Tool:** PortMap v1.0")
    if ($InputFile) {
        [void]$markdownBuilder.AppendLine("**Input File:** $(Split-Path -Leaf $InputFile)")
    }
    [void]$markdownBuilder.AppendLine("**Output Format:** $($OutputFormat ?? 'Markdown')")
    [void]$markdownBuilder.AppendLine("**Generated by:** $env:USERNAME on $env:COMPUTERNAME")
    [void]$markdownBuilder.AppendLine("")
    
    # Device Information Summary
    $deviceMakes = $Devices | Select-Object -ExpandProperty deviceMake -Unique | Sort-Object
    $deviceModels = $Devices | Select-Object -ExpandProperty deviceModel -Unique | Sort-Object
    [void]$markdownBuilder.AppendLine("**Device Makes:** $($deviceMakes -join ', ')")
    [void]$markdownBuilder.AppendLine("**Device Models:** $($deviceModels -join ', ')")
    [void]$markdownBuilder.AppendLine("**Total Devices:** $($Devices.Count)")
    [void]$markdownBuilder.AppendLine("")
    
    # Device Summary Table
    [void]$markdownBuilder.AppendLine("## Device Summary")
    [void]$markdownBuilder.AppendLine("")
    [void]$markdownBuilder.AppendLine("| Device Name | Make | Model | Total Ports | Used Ports | Unused Ports |")
    [void]$markdownBuilder.AppendLine("|-------------|------|-------|-------------|------------|--------------|")
    
    foreach ($device in $Devices) {
        $portInfo = $DevicePortInfo[$device.deviceName]
        $usedCount = if ($portInfo.UsedPorts) { @($portInfo.UsedPorts).Count } else { 0 }
        $unusedCount = if ($portInfo.UnusedPorts) { @($portInfo.UnusedPorts).Count } else { 0 }
        [void]$markdownBuilder.AppendLine("| $($device.deviceName) | $($device.deviceMake) | $($device.deviceModel) | $($portInfo.TotalPorts) | $usedCount | $unusedCount |")
    }
    
    [void]$markdownBuilder.AppendLine("")
    
    # Connection and Port Tables - Individual tables per device
    [void]$markdownBuilder.AppendLine("## Port Mapping by Device")
    [void]$markdownBuilder.AppendLine("")
    
    foreach ($device in $Devices) {
        $deviceName = $device.deviceName
        $deviceMake = $device.deviceMake
        $deviceModel = $device.deviceModel
        $portInfo = $DevicePortInfo[$deviceName]
        
        # Get connections for this device
        $deviceConnections = if ($ConnectionMappings) { 
            $ConnectionMappings | Where-Object { $_.SourceDevice -eq $deviceName } | Sort-Object { [int]$_.SourcePort }
        }
        else { 
            @() 
        }
        
        # Calculate counts
        $connectionCount = if ($deviceConnections) { @($deviceConnections).Count } else { 0 }
        $unusedPortCount = if ($ShowUnused -and $portInfo.UnusedPorts) { @($portInfo.UnusedPorts).Count } else { 0 }
        $totalRowCount = $connectionCount + $unusedPortCount
        
        [void]$markdownBuilder.AppendLine("### $deviceName")
        [void]$markdownBuilder.AppendLine("")
        [void]$markdownBuilder.AppendLine("**Device:** $deviceName ($deviceMake $deviceModel)")
        [void]$markdownBuilder.AppendLine("**Total Ports:** $($portInfo.TotalPorts)")
        [void]$markdownBuilder.AppendLine("**Connected Ports:** $connectionCount")
        if ($ShowUnused) {
            [void]$markdownBuilder.AppendLine("**Unused Ports:** $unusedPortCount")
        }
        [void]$markdownBuilder.AppendLine("")
        
        if ($totalRowCount -gt 0) {
            [void]$markdownBuilder.AppendLine("| Port | Media | Status | Destination Device | Destination Port | Destination Media | Type | Notes |")
            [void]$markdownBuilder.AppendLine("|------|-------|--------|-------------------|------------------|-------------------|------|-------|")
            
            # Create a unified list of all ports to display, sorted by port number
            $allPortEntries = [System.Collections.Generic.List[object]]::new()
            
            # Add connected ports
            foreach ($connection in $deviceConnections) {
                $portEntry = [PSCustomObject]@{
                    PortNumber        = [int]$connection.SourcePort
                    Port              = $connection.SourcePort
                    Media             = $connection.SourceMedia
                    Status            = $connection.Status
                    DestinationDevice = $connection.DestinationDevice
                    DestinationPort   = $connection.DestinationPort
                    DestinationMedia  = $connection.DestinationMedia
                    Type              = $connection.ConnectionType
                    Notes             = if ($connection.Notes) { $connection.Notes } else { "" }
                    IsConnected       = $true
                }
                $allPortEntries.Add($portEntry)
            }
            
            # Add unused ports if requested
            if ($ShowUnused -and $unusedPortCount -gt 0) {
                foreach ($port in $portInfo.UnusedPorts) {
                    $details = $portInfo.PortDetails[$port]
                    $portEntry = [PSCustomObject]@{
                        PortNumber        = [int]$port
                        Port              = $port
                        Media             = $details.MediaType
                        Status            = "**Unused**"
                        DestinationDevice = "-"
                        DestinationPort   = "-"
                        DestinationMedia  = "-"
                        Type              = "-"
                        Notes             = "Available"
                        IsConnected       = $false
                    }
                    $allPortEntries.Add($portEntry)
                }
            }
            
            # Sort all entries by port number and output
            $sortedPortEntries = $allPortEntries | Sort-Object PortNumber
            foreach ($portEntry in $sortedPortEntries) {
                [void]$markdownBuilder.AppendLine("| $($portEntry.Port) | $($portEntry.Media) | $($portEntry.Status) | $($portEntry.DestinationDevice) | $($portEntry.DestinationPort) | $($portEntry.DestinationMedia) | $($portEntry.Type) | $($portEntry.Notes) |")
            }
        }
        else {
            [void]$markdownBuilder.AppendLine("*No port connections or unused ports to display for this device.*")
        }
        
        [void]$markdownBuilder.AppendLine("")
    }
    
    return $markdownBuilder.ToString()
}

function ConvertTo-CsvOutput {
    <#
    .SYNOPSIS
        Converts port mapping data to individual CSV files per device.
    
    .DESCRIPTION
        Creates individual CSV files for each device with structured data suitable for 
        spreadsheet applications and data processing. Each device gets its own CSV file
        with complete port information, connections, and metadata.
    
    .PARAMETER Devices
        Array of device configuration objects.
    
    .PARAMETER DevicePortInfo
        Hashtable containing device port information.
    
    .PARAMETER ConnectionMappings
        Array of connection mapping objects.
    
    .PARAMETER InputFile
        Path to the input JSON configuration file.
    
    .PARAMETER OutputFormat
        The output format being used.
    
    .PARAMETER ShowUnused
        Include unused ports in the CSV output.
    
    .OUTPUTS
        System.Collections.Hashtable
        Hashtable containing device names as keys and CSV data arrays as values.
    #>
    [CmdletBinding()]
    [OutputType([hashtable])]
    param(
        [Parameter(Mandatory = $true)]
        [object[]]$Devices,
        
        [Parameter(Mandatory = $true)]
        [hashtable]$DevicePortInfo,
        
        [Parameter(Mandatory = $true)]
        [AllowEmptyCollection()]
        [object[]]$ConnectionMappings,
        
        [Parameter(Mandatory = $false)]
        [string]$InputFile,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFormat,
        
        [Parameter(Mandatory = $false)]
        [switch]$ShowUnused
    )
    
    # Create a hashtable to store CSV data for each device
    $deviceCsvFiles = @{}
    
    foreach ($device in $Devices) {
        $deviceName = $device.deviceName
        $portInfo = $DevicePortInfo[$deviceName]
        $csvData = [System.Collections.Generic.List[object]]::new()
        
        # CSV output excludes metadata - only port data
        
        # Create unified port list for this device (similar to Markdown logic)
        $allPortEntries = [System.Collections.Generic.List[object]]::new()
        
        # Get connections for this device
        $deviceConnections = if ($ConnectionMappings) { 
            $ConnectionMappings | Where-Object { $_.SourceDevice -eq $deviceName } | Sort-Object { [int]$_.SourcePort }
        }
        else { 
            @() 
        }
        
        # Add connected ports
        foreach ($connection in $deviceConnections) {
            $portEntry = [PSCustomObject]@{
                PortNumber        = [int]$connection.SourcePort
                Port              = $connection.SourcePort
                Media             = $connection.SourceMedia
                Status            = $connection.Status
                DestinationDevice = $connection.DestinationDevice
                DestinationPort   = $connection.DestinationPort
                DestinationMedia  = $connection.DestinationMedia
                Type              = $connection.ConnectionType
                Notes             = if ($connection.Notes) { $connection.Notes } else { "" }
                IsConnected       = $true
            }
            $allPortEntries.Add($portEntry)
        }
        
        # Add unused ports if requested
        if ($ShowUnused -and $portInfo.UnusedPorts) {
            foreach ($port in $portInfo.UnusedPorts) {
                $details = $portInfo.PortDetails[$port]
                $portEntry = [PSCustomObject]@{
                    PortNumber        = [int]$port
                    Port              = $port
                    Media             = $details.MediaType
                    Status            = "Unused"
                    DestinationDevice = ""
                    DestinationPort   = ""
                    DestinationMedia  = ""
                    Type              = ""
                    Notes             = "Available"
                    IsConnected       = $false
                }
                $allPortEntries.Add($portEntry)
            }
        }
        
        # Sort all entries by port number and add to CSV data
        $sortedPortEntries = $allPortEntries | Sort-Object PortNumber
        foreach ($portEntry in $sortedPortEntries) {
            $csvRow = [PSCustomObject]@{
                DeviceName        = $deviceName
                Make              = $device.deviceMake
                Model             = $device.deviceModel
                Location          = $device.location ?? ""
                Port              = $portEntry.Port
                Media             = $portEntry.Media
                Status            = $portEntry.Status
                DestinationDevice = $portEntry.DestinationDevice
                DestinationPort   = $portEntry.DestinationPort
                DestinationMedia  = $portEntry.DestinationMedia
                Type              = $portEntry.Type
                Notes             = $portEntry.Notes
            }
            $csvData.Add($csvRow)
        }
        
        # Store CSV data for this device
        $deviceCsvFiles[$deviceName] = $csvData.ToArray()
    }
    
    return $deviceCsvFiles
}

function ConvertTo-JsonOutput {
    <#
    .SYNOPSIS
        Converts port mapping data to JSON format.
    
    .DESCRIPTION
        Creates machine-readable JSON output suitable for integration with other tools
        and programmatic consumption.
    
    .PARAMETER Devices
        Array of device configuration objects.
    
    .PARAMETER DevicePortInfo
        Hashtable containing device port information.
    
    .PARAMETER ConnectionMappings
        Array of connection mapping objects.
    
    .PARAMETER InputFile
        Path to the input JSON configuration file.
    
    .PARAMETER OutputFormat
        The output format being used.
    
    .OUTPUTS
        System.String
        JSON-formatted string containing all port mapping data.
    #>
    [CmdletBinding()]
    [OutputType([string])]
    param(
        [Parameter(Mandatory = $true)]
        [object[]]$Devices,
        
        [Parameter(Mandatory = $true)]
        [hashtable]$DevicePortInfo,
        
        [Parameter(Mandatory = $true)]
        [AllowEmptyCollection()]
        [object[]]$ConnectionMappings,
        
        [Parameter(Mandatory = $false)]
        [string]$InputFile,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFormat
    )
    
    # Collect device information for metadata
    $deviceMakes = $Devices | Select-Object -ExpandProperty deviceMake -Unique | Sort-Object
    $deviceModels = $Devices | Select-Object -ExpandProperty deviceModel -Unique | Sort-Object
    $deviceNames = $Devices | Select-Object -ExpandProperty deviceName | Sort-Object
    
    $jsonOutput = [ordered]@{
        metadata    = [ordered]@{
            generatedOn   = (Get-Date -Format "o")
            tool          = "PortMap"
            version       = "1.0"
            generatedBy   = $env:USERNAME
            computerName  = $env:COMPUTERNAME
            inputFile     = if ($InputFile) { Split-Path -Leaf $InputFile } else { "Unknown" }
            outputFormat  = $OutputFormat ?? "JSON"
            deviceSummary = [ordered]@{
                deviceNames  = $deviceNames
                deviceMakes  = $deviceMakes
                deviceModels = $deviceModels
                deviceCount  = $Devices.Count
            }
        }
        summary     = [ordered]@{
            totalDevices     = if ($Devices) { @($Devices).Count } else { 0 }
            totalConnections = if ($ConnectionMappings) { @($ConnectionMappings).Count } else { 0 }
            totalPorts       = ($DevicePortInfo.Values | Measure-Object -Property TotalPorts -Sum).Sum
            totalUsedPorts   = ($DevicePortInfo.Values | ForEach-Object { if ($_.UsedPorts) { @($_.UsedPorts).Count } else { 0 } } | Measure-Object -Sum).Sum
        }
        devices     = [System.Collections.Generic.List[object]]::new()
        connections = $ConnectionMappings
    }
    
    foreach ($device in $Devices) {
        $portInfo = $DevicePortInfo[$device.deviceName]
        
        $deviceData = [ordered]@{
            deviceName  = $device.deviceName
            deviceMake  = $device.deviceMake
            deviceModel = $device.deviceModel
            location    = $device.location ?? ""
            portSummary = [ordered]@{
                totalPorts         = $portInfo.TotalPorts
                usedPorts          = if ($portInfo.UsedPorts) { @($portInfo.UsedPorts).Count } else { 0 }
                unusedPorts        = if ($portInfo.UnusedPorts) { @($portInfo.UnusedPorts).Count } else { 0 }
                utilizationPercent = if ($portInfo.TotalPorts -gt 0) { 
                    $usedPortCount = if ($portInfo.UsedPorts) { @($portInfo.UsedPorts).Count } else { 0 }
                    [math]::Round(($usedPortCount / $portInfo.TotalPorts) * 100, 2) 
                }
                else { 
                    0 
                }
            }
            portDetails = [System.Collections.Generic.List[object]]::new()
        }
        
        $sortedPorts = $portInfo.PortDetails.Keys | Sort-Object
        foreach ($portNum in $sortedPorts) {
            $portDetail = $portInfo.PortDetails[$portNum]
            $portDetailObj = [ordered]@{
                port        = $portNum
                mediaType   = $portDetail.MediaType
                speed       = $portDetail.Speed
                description = $portDetail.Description
                isUsed      = $portDetail.IsUsed
                connection  = $portDetail.Connection
            }
            $deviceData.portDetails.Add($portDetailObj)
        }
        
        $unusedPortCount = if ($portInfo.UnusedPorts) { @($portInfo.UnusedPorts).Count } else { 0 }
        if ($ShowUnused -and $unusedPortCount -gt 0) {
            $deviceData.unusedPorts = ($portInfo.UnusedPorts | Sort-Object)
        }
        
        $jsonOutput.devices.Add($deviceData)
    }
    
    try {
        return ($jsonOutput | ConvertTo-Json -Depth 15 -Compress:$false)
    }
    catch {
        Write-PortMapLog "Failed to convert to JSON: $($_.Exception.Message)" -Level Error
        throw
    }
}

#endregion

#region Main Execution

function Start-PortMappingProcess {
    <#
    .SYNOPSIS
        Main execution function for the PortMap tool.
    
    .DESCRIPTION
        Orchestrates the entire port mapping process including configuration loading,
        validation, processing, and output generation.
    
    .NOTES
        This function handles all error conditions and provides comprehensive logging.
    #>
    [CmdletBinding()]
    param()
    
    try {
        Write-PortMapLog "Starting PortMap Tool v1.0" -Level Info
        Write-PortMapLog "Input file: $InputFile" -Level Info
        Write-PortMapLog "Output format: $OutputFormat" -Level Info
        
        # Load and validate configuration
        Write-PortMapLog "Loading configuration from $InputFile..." -Level Info
        
        if (-not (Test-Path -Path $InputFile -PathType Leaf)) {
            throw "Input file does not exist: $InputFile"
        }
        
        $configContent = $null
        try {
            $configContent = Get-Content -Path $InputFile -Raw -ErrorAction Stop | ConvertFrom-Json -ErrorAction Stop
        }
        catch {
            throw "Failed to parse JSON from input file: $($_.Exception.Message)"
        }
        
        if (-not (Test-NetworkConfiguration -Configuration $configContent)) {
            throw "Configuration validation failed. Please check the error messages above."
        }
        
        if ($Validate) {
            Write-PortMapLog "Configuration validation completed successfully" -Level Info
            return 0
        }
        
        # Filter devices if specified
        $devices = [array]$configContent.devices
        if ($DeviceFilter -and @($DeviceFilter).Count -gt 0) {
            $devices = $devices | Where-Object { $_.deviceName -in $DeviceFilter }
            Write-PortMapLog "Filtered to devices: $($DeviceFilter -join ', ')" -Level Info
            
            if (@($devices).Count -eq 0) {
                Write-PortMapLog "No devices matched the filter criteria" -Level Warning
                return 1
            }
        }
        
        # Build device port information
        Write-PortMapLog "Processing device port information..." -Level Info
        $devicePortInfo = @{}
        
        foreach ($device in $devices) {
            try {
                $devicePortInfo[$device.deviceName] = Get-NetworkDevicePortInfo -Device $device
                Write-PortMapLog "Processed device: $($device.deviceName) ($($devicePortInfo[$device.deviceName].TotalPorts) ports)" -Level Info
            }
            catch {
                Write-PortMapLog "Failed to process device $($device.deviceName): $($_.Exception.Message)" -Level Error
                throw
            }
        }
        
        # Process connections
        Write-PortMapLog "Processing connections..." -Level Info
        $connectionMappings = @()
        
        if ($configContent.connections -and @($configContent.connections).Count -gt 0) {
            try {
                $connectionMappings = New-ConnectionMappings -Connections $configContent.connections -DevicePortInfo $devicePortInfo
                $connectionCount = if ($connectionMappings) { @($connectionMappings).Count } else { 0 }
                Write-PortMapLog "Processed $connectionCount connections" -Level Info
            }
            catch {
                Write-PortMapLog "Failed to process connections: $($_.Exception.Message)" -Level Error
                throw
            }
        }
        else {
            Write-PortMapLog "No connections defined in configuration" -Level Info
        }
        
        # Generate output
        Write-PortMapLog "Generating $OutputFormat output..." -Level Info
        $output = $null
        
        try {
            $output = switch ($OutputFormat) {
                "Markdown" { 
                    ConvertTo-MarkdownOutput -Devices $devices -DevicePortInfo $devicePortInfo -ConnectionMappings $connectionMappings -InputFile $InputFile -OutputFormat $OutputFormat
                }
                "CSV" { 
                    ConvertTo-CsvOutput -Devices $devices -DevicePortInfo $devicePortInfo -ConnectionMappings $connectionMappings -InputFile $InputFile -OutputFormat $OutputFormat -ShowUnused:$ShowUnused
                }
                "JSON" { 
                    ConvertTo-JsonOutput -Devices $devices -DevicePortInfo $devicePortInfo -ConnectionMappings $connectionMappings -InputFile $InputFile -OutputFormat $OutputFormat
                }
                default {
                    throw "Unsupported output format: $OutputFormat"
                }
            }
        }
        catch {
            Write-PortMapLog "Failed to generate $OutputFormat output: $($_.Exception.Message)" -Level Error
            throw
        }
        
        # Handle output generation
        if ($OutputFormat -eq "CSV") {
            # CSV output creates individual files per device
            Write-PortMapLog "Creating individual CSV files per device..." -Level Info
            $outputFiles = @()
            
            try {
                foreach ($deviceName in $output.Keys) {
                    # Generate unique filename for each device
                    $deviceOutputFile = New-UniqueOutputFileName -InputFile $InputFile -OutputFormat $OutputFormat -Devices @($devices | Where-Object { $_.deviceName -eq $deviceName }) -OutputFile $OutputFile
                    
                    # Create output directory if needed
                    $outputDir = Split-Path -Path $deviceOutputFile -Parent
                    if ($outputDir -and -not (Test-Path -Path $outputDir)) {
                        New-Item -Path $outputDir -ItemType Directory -Force | Out-Null
                        Write-PortMapLog "Created output directory: $outputDir" -Level Info
                    }
                    
                    # Export CSV for this device
                    $output[$deviceName] | Export-Csv -Path $deviceOutputFile -NoTypeInformation -Encoding UTF8
                    
                    # Track output file information
                    $fileInfo = Get-Item -Path $deviceOutputFile
                    $outputFiles += $fileInfo
                    
                    Write-PortMapLog "Created CSV file for $deviceName" -Level Info
                    Write-PortMapLog "  File: $($fileInfo.Name)" -Level Info
                    Write-PortMapLog "  Size: $($fileInfo.Length) bytes" -Level Info
                }
                
                # Display summary of all files created
                Write-PortMapLog "CSV output files created:" -Level Info
                foreach ($fileInfo in $outputFiles) {
                    Write-PortMapLog "  $($fileInfo.FullName)" -Level Info
                }
                
                # Console output for first device if no explicit output file specified
                if (-not $OutputFile -and $outputFiles.Count -gt 0) {
                    Write-Host ""
                    Write-PortMapLog "Console Output (Sample - First Device):" -Level Info
                    $firstDeviceData = $output[($output.Keys | Select-Object -First 1)]
                    $firstDeviceData | Format-Table -AutoSize
                }
            }
            catch {
                Write-PortMapLog "Failed to write CSV output files: $($_.Exception.Message)" -Level Error
                throw
            }
        }
        else {
            # Single file output for Markdown and JSON
            $finalOutputFile = New-UniqueOutputFileName -InputFile $InputFile -OutputFormat $OutputFormat -Devices $devices -OutputFile $OutputFile
            
            try {
                $outputDir = Split-Path -Path $finalOutputFile -Parent
                if ($outputDir -and -not (Test-Path -Path $outputDir)) {
                    New-Item -Path $outputDir -ItemType Directory -Force | Out-Null
                    Write-PortMapLog "Created output directory: $outputDir" -Level Info
                }
                
                $output | Out-File -FilePath $finalOutputFile -Encoding UTF8 -Force
                
                # Display file information
                $fileInfo = Get-Item -Path $finalOutputFile
                Write-PortMapLog "Output written to file:" -Level Info
                Write-PortMapLog "  File: $($fileInfo.Name)" -Level Info
                Write-PortMapLog "  Location: $($fileInfo.DirectoryName)" -Level Info
                Write-PortMapLog "  Full Path: $($fileInfo.FullName)" -Level Info
                Write-PortMapLog "  Format: $OutputFormat" -Level Info
                Write-PortMapLog "  Size: $($fileInfo.Length) bytes" -Level Info
                
                # Also output to console if not explicitly specifying an output file
                if (-not $OutputFile) {
                    Write-Host ""
                    Write-PortMapLog "Console Output:" -Level Info
                    Write-Output $output
                }
            }
            catch {
                Write-PortMapLog "Failed to write output file: $($_.Exception.Message)" -Level Error
                throw
            }
        }
        
        # Display summary statistics
        $totalPorts = ($devicePortInfo.Values | Measure-Object -Property TotalPorts -Sum).Sum
        $usedPorts = ($devicePortInfo.Values | ForEach-Object { if ($_.UsedPorts) { @($_.UsedPorts).Count } else { 0 } } | Measure-Object -Sum).Sum
        $utilizationRate = if ($totalPorts -gt 0) { [math]::Round(($usedPorts / $totalPorts) * 100, 2) } else { 0 }
        $connectionCount = if ($connectionMappings) { @($connectionMappings).Count } else { 0 }
        
        Write-PortMapLog "Processing Summary:" -Level Info
        Write-PortMapLog "  Devices processed: $(if ($devices) { @($devices).Count } else { 0 })" -Level Info
        Write-PortMapLog "  Total ports: $totalPorts" -Level Info
        Write-PortMapLog "  Used ports: $usedPorts" -Level Info
        Write-PortMapLog "  Port utilization: $utilizationRate%" -Level Info
        Write-PortMapLog "  Connections mapped: $connectionCount" -Level Info
        
        if ($OutputFormat -eq "CSV") {
            Write-PortMapLog "  Output files: $($outputFiles.Count) CSV files (one per device)" -Level Info
            Write-PortMapLog "  Output location: $(Split-Path -Parent $outputFiles[0].FullName)" -Level Info
        }
        else {
            Write-PortMapLog "  Output file: $(Split-Path -Leaf $finalOutputFile)" -Level Info
            Write-PortMapLog "  Output location: $(Split-Path -Parent $finalOutputFile)" -Level Info
        }
        
        Write-PortMapLog "PortMap processing completed successfully" -Level Info
        return 0
        
    }
    catch {
        Write-PortMapLog "Critical error in PortMap processing: $($_.Exception.Message)" -Level Error
        
        if ($_.Exception.InnerException) {
            Write-PortMapLog "Inner exception: $($_.Exception.InnerException.Message)" -Level Error
        }
        
        # Write stack trace for debugging (only if in verbose mode)
        if ($VerbosePreference -eq 'Continue') {
            Write-PortMapLog "Stack trace: $($_.ScriptStackTrace)" -Level Error
        }
        
        return 1
    }
}

# Execute main function and exit with appropriate code
$exitCode = Start-PortMappingProcess
exit $exitCode

#endregion
