#Requires -Version 5.1

<#
.SYNOPSIS
    EnvironmentDetailAssignment PowerShell Module for Azure Stack Network Switch Config Generator
    
.DESCRIPTION
    This module provides comprehensive environment configuration management functionality for 
    network device configurations. It enables seamless deployment across multiple environments
    (Development, Staging, Production) by providing environment-aware parameter assignment
    and template variable generation.

.NOTES
    File Name      : EnvironmentDetailAssignment.psm1
    Author         : Network Engineering Team
    Prerequisite   : PowerShell 5.1 or later
    Copyright      : (c) 2025 Azure Stack Framework. All rights reserved.
    Version        : 1.0.0
#>

# Set strict mode for better error handling
Set-StrictMode -Version Latest

# Module-level variables
$Script:LogLevel = "Info"
$Script:ProcessingStats = @{
    TotalEnvironments = 0
    ProcessedEnvironments = 0
    TotalDevices = 0
    ProcessedDevices = 0
    TotalParameters = 0
    ValidationErrors = 0
    StartTime = $null
    EndTime = $null
}

#region Helper Functions

function Write-LogMessage {
    <#
    .SYNOPSIS
        Writes formatted log messages with timestamps
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [ValidateSet("Info", "Warning", "Error", "Debug")]
        [string]$Level = "Info",
        
        [Parameter(Mandatory = $false)]
        [switch]$NoTimestamp
    )
    
    $timestamp = if (-not $NoTimestamp) { "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss'): " } else { "" }
    
    switch ($Level) {
        "Info"    { Write-Host "$timestamp$Message" -ForegroundColor Green }
        "Warning" { Write-Warning "$timestamp$Message" }
        "Error"   { Write-Error "$timestamp$Message" }
        "Debug"   { if ($Script:LogLevel -eq "Debug") { Write-Host "$timestamp[DEBUG] $Message" -ForegroundColor Cyan } }
    }
}

function Test-JsonConfiguration {
    <#
    .SYNOPSIS
        Validates JSON configuration structure and content
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [hashtable]$Configuration
    )
    
    Write-LogMessage "Validating configuration structure..." -Level Debug
    
    $validationErrors = @()
    
    # Check for required top-level properties
    if (-not $Configuration.ContainsKey("environments")) {
        $validationErrors += "Missing required 'environments' section"
    }
    
    # Validate environments section
    if ($Configuration.environments) {
        foreach ($envName in $Configuration.environments.Keys) {
            $env = $Configuration.environments[$envName]
            
            # Check for required environment properties
            $requiredProps = @("networkPrefix", "vlanBase", "deviceSuffix", "managementVlan")
            foreach ($prop in $requiredProps) {
                if (-not $env.ContainsKey($prop)) {
                    $validationErrors += "Environment '$envName' missing required property '$prop'"
                }
            }
            
            # Validate VLAN base is numeric
            if ($env.vlanBase -and $env.vlanBase -notmatch '^\d+$') {
                $validationErrors += "Environment '$envName' vlanBase must be numeric"
            }
            
            # Validate network prefix format
            if ($env.networkPrefix -and $env.networkPrefix -notmatch '^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$') {
                $validationErrors += "Environment '$envName' networkPrefix must be in CIDR format (e.g., 192.168.1.0/24)"
            }
        }
    }
    
    if (@($validationErrors).Count -gt 0) {
        $Script:ProcessingStats.ValidationErrors = @($validationErrors).Count
        throw "Configuration validation failed with $(@($validationErrors).Count) errors:`n" + ($validationErrors -join "`n")
    }
    
    Write-LogMessage "Configuration validation completed successfully" -Level Debug
    return $true
}

function Resolve-ParameterTemplate {
    <#
    .SYNOPSIS
        Resolves parameter templates with variable substitution
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$Template,
        
        [Parameter(Mandatory = $true)]
        [hashtable]$Variables
    )
    
    $result = $Template
    foreach ($key in $Variables.Keys) {
        $result = $result -replace "\{$key\}", $Variables[$key]
    }
    
    return $result
}

function Get-IntelligentFileName {
    <#
    .SYNOPSIS
        Generates intelligent output filenames based on configuration and parameters
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$ConfigPath,
        
        [Parameter(Mandatory = $true)]
        [string]$OutputFormat,
        
        [Parameter(Mandatory = $false)]
        [string]$Environment,
        
        [Parameter(Mandatory = $false)]
        [string[]]$Environments,
        
        [Parameter(Mandatory = $false)]
        [string]$BasePath = "."
    )
    
    $configBaseName = [System.IO.Path]::GetFileNameWithoutExtension($ConfigPath)
    
    # Determine environment suffix
    $envSuffix = if ($Environment) {
        $Environment
    } elseif ($Environments -and @($Environments).Count -gt 0) {
        if (@($Environments).Count -eq 1) {
            $Environments[0]
        } else {
            "multi-env"
        }
    } else {
        "all-environments"
    }
    
    # Determine file extension
    $extension = switch ($OutputFormat.ToLower()) {
        "json"       { ".json" }
        "csv"        { ".csv" }
        "powershell" { ".ps1" }
        default      { ".txt" }
    }
    
    # Create base filename
    $baseFileName = "$configBaseName-$envSuffix-environment-assignment"
    
    if ($OutputFormat.ToLower() -eq "csv" -and $envSuffix -eq "all-environments") {
        $baseFileName = "$configBaseName-$envSuffix-matrix"
    }
    
    $fileName = "$baseFileName$extension"
    $fullPath = Join-Path $BasePath $fileName
    
    # Handle file conflicts with incremental naming
    $counter = 1
    while (Test-Path $fullPath) {
        $fileName = "$baseFileName-$counter$extension"
        $fullPath = Join-Path $BasePath $fileName
        $counter++
    }
    
    return $fullPath
}

function Merge-EnvironmentParameters {
    <#
    .SYNOPSIS
        Merges global defaults with environment-specific parameters
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $false)]
        [hashtable]$GlobalDefaults = @{},
        
        [Parameter(Mandatory = $true)]
        [hashtable]$EnvironmentConfig
    )
    
    $mergedConfig = $GlobalDefaults.Clone()
    
    # Merge environment-specific parameters, overriding globals
    foreach ($key in $EnvironmentConfig.Keys) {
        $mergedConfig[$key] = $EnvironmentConfig[$key]
    }
    
    return $mergedConfig
}

#endregion

#region Main Functions

function New-EnvironmentAssignment {
    <#
    .SYNOPSIS
        Creates environment-specific parameter assignments from configuration
        
    .DESCRIPTION
        This function processes a JSON configuration file containing environment definitions
        and generates environment-specific parameter assignments. It supports multiple
        output formats including JSON, CSV, and PowerShell variables.
        
    .PARAMETER ConfigPath
        Path to the JSON configuration file containing environment definitions
        
    .PARAMETER OutputFormat
        Output format for the generated configuration. Valid values: JSON, CSV, PowerShell
        
    .PARAMETER Environment
        Target specific environment for processing. If not specified, all environments are processed
        
    .PARAMETER Environments
        Array of specific environments to process. Cannot be used with -Environment parameter
        
    .PARAMETER OutputFile
        Custom output file path. If not specified, intelligent naming is used
        
    .PARAMETER TemplateCompatible
        Generate template-compatible variable format for PowerShell output
        
    .PARAMETER IncludeDocumentation
        Include documentation and metadata in output (JSON and CSV formats)
        
    .PARAMETER Validate
        Validate configuration without generating output
        
    .PARAMETER Detailed
        Display detailed progress information during processing
        
    .EXAMPLE
        New-EnvironmentAssignment -ConfigPath "env-config.json" -OutputFormat JSON
        
        Generates JSON configuration for all environments defined in the configuration file
        
    .EXAMPLE
        New-EnvironmentAssignment -ConfigPath "env-config.json" -Environment "Production" -OutputFormat PowerShell -TemplateCompatible
        
        Generates PowerShell variables for the Production environment in template-compatible format
        
    .EXAMPLE
        New-EnvironmentAssignment -ConfigPath "env-config.json" -OutputFormat CSV -IncludeDocumentation
        
        Generates a CSV matrix showing parameter variations across all environments with documentation
        
    .OUTPUTS
        Generates output files in the specified format containing environment parameter assignments
        
    .NOTES
        This function is the main entry point for the EnvironmentDetailAssignment module
    #>
    [CmdletBinding(DefaultParameterSetName = "AllEnvironments")]
    param(
        [Parameter(Mandatory = $true, Position = 0)]
        [ValidateScript({
            if (-not (Test-Path $_)) {
                throw "Configuration file not found: $_"
            }
            if ([System.IO.Path]::GetExtension($_) -ne ".json") {
                throw "Configuration file must be a JSON file"
            }
            return $true
        })]
        [string]$ConfigPath,
        
        [Parameter(Mandatory = $true)]
        [ValidateSet("JSON", "CSV", "PowerShell")]
        [string]$OutputFormat,
        
        [Parameter(ParameterSetName = "SingleEnvironment")]
        [string]$Environment,
        
        [Parameter(ParameterSetName = "MultipleEnvironments")]
        [string[]]$Environments,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFile,
        
        [Parameter(Mandatory = $false)]
        [switch]$TemplateCompatible,
        
        [Parameter(Mandatory = $false)]
        [switch]$IncludeDocumentation,
        
        [Parameter(Mandatory = $false)]
        [switch]$Validate,
        
        [Parameter(Mandatory = $false)]
        [switch]$Detailed
    )
    
    try {
        # Initialize processing
        $Script:ProcessingStats.StartTime = Get-Date
        if ($Detailed) { $Script:LogLevel = "Debug" }
        
        Write-LogMessage "Starting EnvironmentDetailAssignment processing..." -Level Info
        Write-LogMessage "Configuration file: $ConfigPath" -Level Debug
        Write-LogMessage "Output format: $OutputFormat" -Level Debug
        
        # Load and parse configuration
        Write-LogMessage "Loading configuration file..." -Level Debug
        $configContent = Get-Content -Path $ConfigPath -Raw | ConvertFrom-Json -AsHashtable
        
        # Validate configuration
        Test-JsonConfiguration -Configuration $configContent
        
        if ($Validate) {
            Write-LogMessage "Configuration validation completed successfully" -Level Info
            return
        }
        
        # Determine target environments
        $targetEnvironments = if ($Environment) {
            @($Environment)
        } elseif ($Environments) {
            @($Environments)
        } else {
            @($configContent.environments.Keys)
        }
        
        Write-LogMessage "Target environments: $($targetEnvironments -join ', ')" -Level Debug
        $Script:ProcessingStats.TotalEnvironments = @($targetEnvironments).Count
        
        # Process environments
        $processedData = @{
            metadata = @{
                generatedBy = "EnvironmentDetailAssignment"
                generatedAt = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
                configFile = Split-Path -Leaf $ConfigPath
                environments = $targetEnvironments
                outputFormat = $OutputFormat
            }
            environments = @{}
        }
        
        foreach ($envName in $targetEnvironments) {
            if (-not $configContent.environments.ContainsKey($envName)) {
                Write-LogMessage "Environment '$envName' not found in configuration" -Level Warning
                continue
            }
            
            Write-LogMessage "Processing environment: $envName" -Level Debug
            
            $envConfig = $configContent.environments[$envName]
            $globalDefaults = if ($configContent.globalDefaults) { $configContent.globalDefaults } else { @{} }
            
            # Merge global defaults with environment config
            $mergedConfig = Merge-EnvironmentParameters -GlobalDefaults $globalDefaults -EnvironmentConfig $envConfig
            
            # Process devices in environment
            $processedDevices = @{}
            if ($envConfig.ContainsKey("devices") -and $envConfig.devices) {
                foreach ($deviceName in $envConfig.devices.Keys) {
                    $deviceConfig = $envConfig.devices[$deviceName].Clone()
                    
                    # Generate hostname if not specified
                    if (-not $deviceConfig.ContainsKey("hostname") -and $configContent.ContainsKey("parameterTemplates") -and $configContent.parameterTemplates.ContainsKey("hostnamePattern")) {
                        $variables = @{
                            deviceName = $deviceName
                            environment = $mergedConfig.deviceSuffix
                        }
                        $deviceConfig.hostname = Resolve-ParameterTemplate -Template $configContent.parameterTemplates.hostnamePattern -Variables $variables
                    } elseif (-not $deviceConfig.ContainsKey("hostname")) {
                        $deviceConfig.hostname = "$deviceName-$($mergedConfig.deviceSuffix)"
                    }
                    
                    # Add environment-specific parameters
                    $deviceConfig.managementVlan = $mergedConfig.managementVlan
                    
                    $processedDevices[$deviceName] = $deviceConfig
                    $Script:ProcessingStats.TotalDevices++
                }
            }
            
            # Build environment data structure
            $environmentData = @{
                globalParameters = @{
                    networkPrefix = $mergedConfig.networkPrefix
                    vlanBase = $mergedConfig.vlanBase
                    deviceSuffix = $mergedConfig.deviceSuffix
                    managementVlan = $mergedConfig.managementVlan
                }
                devices = $processedDevices
            }
            
            # Add service parameters if they exist
            if ($mergedConfig.ContainsKey("ntpServers")) {
                $environmentData.globalParameters.ntpServers = $mergedConfig.ntpServers
            }
            if ($mergedConfig.ContainsKey("dnsServers")) {
                $environmentData.globalParameters.dnsServers = $mergedConfig.dnsServers
            }
            
            $processedData.environments[$envName] = $environmentData
            $Script:ProcessingStats.ProcessedEnvironments++
        }
        
        # Generate output file path
        if (-not $OutputFile) {
            $OutputFile = Get-IntelligentFileName -ConfigPath $ConfigPath -OutputFormat $OutputFormat -Environment $Environment -Environments $Environments
        }
        
        # Generate output based on format
        switch ($OutputFormat.ToUpper()) {
            "JSON" {
                Write-LogMessage "Generating JSON output..." -Level Debug
                $jsonOutput = $processedData | ConvertTo-Json -Depth 10
                $jsonOutput | Out-File -FilePath $OutputFile -Encoding UTF8
            }
            
            "CSV" {
                Write-LogMessage "Generating CSV output..." -Level Debug
                $csvData = @()
                
                # Create parameter matrix
                $allParameters = @("NetworkPrefix", "VlanBase", "DeviceSuffix", "ManagementVlan")
                
                foreach ($param in $allParameters) {
                    $row = @{ Parameter = $param }
                    
                    foreach ($envName in $targetEnvironments) {
                        if ($processedData.environments.ContainsKey($envName)) {
                            $envData = $processedData.environments[$envName].globalParameters
                            switch ($param) {
                                "NetworkPrefix" { $row[$envName] = $envData.networkPrefix }
                                "VlanBase" { $row[$envName] = $envData.vlanBase }
                                "DeviceSuffix" { $row[$envName] = $envData.deviceSuffix }
                                "ManagementVlan" { $row[$envName] = $envData.managementVlan }
                            }
                        }
                    }
                    
                    if ($IncludeDocumentation) {
                        switch ($param) {
                            "NetworkPrefix" { $row.Notes = "Base network range" }
                            "VlanBase" { $row.Notes = "Starting VLAN ID" }
                            "DeviceSuffix" { $row.Notes = "Hostname suffix" }
                            "ManagementVlan" { $row.Notes = "Management VLAN ID" }
                        }
                    }
                    
                    $csvData += New-Object PSObject -Property $row
                }
                
                $csvData | Export-Csv -Path $OutputFile -NoTypeInformation
            }
            
            "POWERSHELL" {
                Write-LogMessage "Generating PowerShell output..." -Level Debug
                $psOutput = @()
                
                $psOutput += "# Generated by EnvironmentDetailAssignment Tool"
                
                if (@($targetEnvironments).Count -eq 1) {
                    $envName = @($targetEnvironments)[0]
                    Write-LogMessage "Looking for environment: '$envName' in processedData" -Level Debug
                    $envData = $processedData.environments[$envName]
                    
                    if ($envData -and $envData.globalParameters) {
                        $psOutput += "# Environment: $envName"
                        $psOutput += "# Generated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
                        $psOutput += ""
                        $psOutput += "# Global Environment Parameters"
                        $psOutput += "`$EnvironmentName = `"$envName`""
                        $psOutput += "`$NetworkPrefix = `"$($envData.globalParameters.networkPrefix)`""
                        $psOutput += "`$VlanBase = $($envData.globalParameters.vlanBase)"
                        $psOutput += "`$DeviceSuffix = `"$($envData.globalParameters.deviceSuffix)`""
                        $psOutput += "`$ManagementVlan = $($envData.globalParameters.managementVlan)"
                        
                        if ($envData.globalParameters.ContainsKey("ntpServers")) {
                            $ntpArray = ($envData.globalParameters.ntpServers | ForEach-Object { "`"$_`"" }) -join ", "
                            $psOutput += "`$NtpServers = @($ntpArray)"
                        }
                        
                        if ($envData.globalParameters.ContainsKey("dnsServers")) {
                            $dnsArray = ($envData.globalParameters.dnsServers | ForEach-Object { "`"$_`"" }) -join ", "
                            $psOutput += "`$DnsServers = @($dnsArray)"
                        }
                        
                        # Add device-specific parameters
                        $deviceKeys = if ($envData.devices) { @($envData.devices.Keys) } else { @() }
                        if (@($deviceKeys).Count -gt 0) {
                            $psOutput += ""
                            $psOutput += "# Device-Specific Parameters"
                            
                            foreach ($deviceName in $envData.devices.Keys) {
                                $device = $envData.devices[$deviceName]
                                $safeName = $deviceName -replace '[^a-zA-Z0-9]', ''
                                
                                if ($device.ContainsKey("hostname")) {
                                    $psOutput += "`$$($safeName)_Hostname = `"$($device.hostname)`""
                                }
                                if ($device.ContainsKey("managementIP")) {
                                    $psOutput += "`$$($safeName)_ManagementIP = `"$($device.managementIP)`""
                                }
                                if ($device.ContainsKey("location")) {
                                    $psOutput += "`$$($safeName)_Location = `"$($device.location)`""
                                }
                            }
                        }
                    } else {
                        $psOutput += "# Error: Environment data not found for $envName"
                    }
                } else {
                    $psOutput += "# Multiple Environments: $($targetEnvironments -join ', ')"
                    $psOutput += "# Generated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
                    $psOutput += ""
                    $psOutput += "# Environment Configurations"
                    
                    foreach ($envName in $targetEnvironments) {
                        if ($processedData.environments.ContainsKey($envName)) {
                            $envData = $processedData.environments[$envName]
                            $psOutput += ""
                            $psOutput += "# $envName Environment"
                            $psOutput += "`$${envName}_NetworkPrefix = `"$($envData.globalParameters.networkPrefix)`""
                            $psOutput += "`$${envName}_VlanBase = $($envData.globalParameters.vlanBase)"
                            $psOutput += "`$${envName}_DeviceSuffix = `"$($envData.globalParameters.deviceSuffix)`""
                            $psOutput += "`$${envName}_ManagementVlan = $($envData.globalParameters.managementVlan)"
                        }
                    }
                }
                
                $psOutput -join "`n" | Out-File -FilePath $OutputFile -Encoding UTF8
            }
        }
        
        # Complete processing
        $Script:ProcessingStats.EndTime = Get-Date
        $duration = $Script:ProcessingStats.EndTime - $Script:ProcessingStats.StartTime
        
        Write-LogMessage "Processing completed successfully" -Level Info
        Write-LogMessage "Output file: $OutputFile" -Level Info
        Write-LogMessage "Processed $($Script:ProcessingStats.ProcessedEnvironments) environments in $($duration.TotalSeconds.ToString('F2')) seconds" -Level Info
        
        if ($Detailed) {
            Write-LogMessage "Processing Statistics:" -Level Info
            Write-LogMessage "  Total Environments: $($Script:ProcessingStats.TotalEnvironments)" -Level Info
            Write-LogMessage "  Processed Environments: $($Script:ProcessingStats.ProcessedEnvironments)" -Level Info
            Write-LogMessage "  Total Devices: $($Script:ProcessingStats.TotalDevices)" -Level Info
            Write-LogMessage "  Validation Errors: $($Script:ProcessingStats.ValidationErrors)" -Level Info
        }
        
    } catch {
        Write-LogMessage "Error during processing: $($_.Exception.Message)" -Level Error
        Write-LogMessage "Stack trace: $($_.ScriptStackTrace)" -Level Debug
        throw
    }
}

#endregion

#region Export Functions

# Export module functions
Export-ModuleMember -Function New-EnvironmentAssignment

#endregion