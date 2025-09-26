<#
.SYNOPSIS
    Environment Detail Assignment Tool for Azure Stack Network Switch Config Generator

.DESCRIPTION
    This script provides comprehensive environment configuration management functionality
    for network device configurations. It enables seamless deployment across multiple
    environments (Development, Staging, Production) by providing environment-aware
    parameter assignment and template variable generation.

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
    .\EnvironmentDetailAssignment.ps1 -ConfigPath "env-config.json" -OutputFormat JSON
    
    Generates JSON configuration for all environments defined in the configuration file

.EXAMPLE
    .\EnvironmentDetailAssignment.ps1 -ConfigPath "env-config.json" -Environment "Production" -OutputFormat PowerShell -TemplateCompatible
    
    Generates PowerShell variables for the Production environment in template-compatible format

.EXAMPLE
    .\EnvironmentDetailAssignment.ps1 -ConfigPath "env-config.json" -OutputFormat CSV -IncludeDocumentation
    
    Generates a CSV matrix showing parameter variations across all environments with documentation

.NOTES
    File Name      : EnvironmentDetailAssignment.ps1
    Author         : Network Engineering Team
    Prerequisite   : PowerShell 5.1 or later
    Copyright      : (c) 2025 Azure Stack Framework. All rights reserved.
    Version        : 1.0.0

.LINK
    https://github.com/microsoft/AzureStack_Network_Switch_Config_Generator
#>

[CmdletBinding(DefaultParameterSetName = "AllEnvironments")]
param(
    [Parameter(Mandatory = $true, Position = 0, HelpMessage = "Path to the JSON configuration file")]
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
    
    [Parameter(Mandatory = $true, HelpMessage = "Output format: JSON, CSV, or PowerShell")]
    [ValidateSet("JSON", "CSV", "PowerShell")]
    [string]$OutputFormat,
    
    [Parameter(ParameterSetName = "SingleEnvironment", HelpMessage = "Target specific environment")]
    [string]$Environment,
    
    [Parameter(ParameterSetName = "MultipleEnvironments", HelpMessage = "Target multiple specific environments")]
    [string[]]$Environments,
    
    [Parameter(Mandatory = $false, HelpMessage = "Custom output file path")]
    [string]$OutputFile,
    
    [Parameter(Mandatory = $false, HelpMessage = "Generate template-compatible variable format")]
    [switch]$TemplateCompatible,
    
    [Parameter(Mandatory = $false, HelpMessage = "Include documentation and metadata in output")]
    [switch]$IncludeDocumentation,
    
    [Parameter(Mandatory = $false, HelpMessage = "Validate configuration without generating output")]
    [switch]$Validate,
    
    [Parameter(Mandatory = $false, HelpMessage = "Display detailed progress information")]
    [switch]$Detailed
)

# Set strict mode for better error handling
Set-StrictMode -Version Latest

# Get the script directory
$ScriptDirectory = Split-Path -Parent $MyInvocation.MyCommand.Path

try {
    # Import the EnvironmentDetailAssignment module
    $ModulePath = Join-Path $ScriptDirectory "EnvironmentDetailAssignment.psm1"
    
    if (-not (Test-Path $ModulePath)) {
        throw "EnvironmentDetailAssignment module not found at: $ModulePath"
    }
    
    Write-Host "🌍 Environment Detail Assignment Tool v1.0.0" -ForegroundColor Cyan
    Write-Host "Part of the Azure Stack Network Switch Config Generator" -ForegroundColor Gray
    Write-Host ""
    
    # Import module
    Import-Module $ModulePath -Force
    
    # Build parameter splat for function call
    $params = @{
        ConfigPath = $ConfigPath
        OutputFormat = $OutputFormat
    }
    
    # Add optional parameters
    if ($Environment) { $params.Environment = $Environment }
    if ($Environments) { $params.Environments = $Environments }
    if ($OutputFile) { $params.OutputFile = $OutputFile }
    if ($TemplateCompatible) { $params.TemplateCompatible = $true }
    if ($IncludeDocumentation) { $params.IncludeDocumentation = $true }
    if ($Validate) { $params.Validate = $true }
    if ($Detailed) { $params.Detailed = $true }
    
    # Call the main function
    New-EnvironmentAssignment @params
    
    if (-not $Validate) {
        Write-Host ""
        Write-Host "✅ Environment detail assignment completed successfully!" -ForegroundColor Green
        
        if ($OutputFile) {
            Write-Host "📄 Output file: $OutputFile" -ForegroundColor Cyan
        }
        
        Write-Host ""
        Write-Host "💡 Next steps:" -ForegroundColor Yellow
        Write-Host "   • Review the generated environment configuration" -ForegroundColor White
        Write-Host "   • Use the output with the main switch configuration generator" -ForegroundColor White
        Write-Host "   • Integrate variables with your Jinja2 templates" -ForegroundColor White
    }
    
} catch {
    Write-Host ""
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($Detailed) {
        Write-Host ""
        Write-Host "📋 Stack trace:" -ForegroundColor Yellow
        Write-Host $_.ScriptStackTrace -ForegroundColor Gray
    }
    
    Write-Host ""
    Write-Host "💡 Troubleshooting tips:" -ForegroundColor Yellow
    Write-Host "   • Verify the JSON configuration file syntax" -ForegroundColor White
    Write-Host "   • Check that all required parameters are present" -ForegroundColor White
    Write-Host "   • Use -Validate to check configuration without generating output" -ForegroundColor White
    Write-Host "   • Use -Detailed for more verbose error information" -ForegroundColor White
    
    exit 1
} finally {
    # Clean up module
    Remove-Module EnvironmentDetailAssignment -ErrorAction SilentlyContinue
}