#Requires -Version 5.1

<#
.SYNOPSIS
    Test script for PortMap.ps1 tool

.DESCRIPTION
    This script tests the functionality of the PortMap tool with various scenarios
    and validates the output formats.

.NOTES
    Version: 1.0
    Author: Network Engineering Team
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $false)]
    [string]$TestDataPath = ".\sample-network-config.json"
)

# Set strict mode for better error handling
Set-StrictMode -Version Latest

# Test configuration using proper collections
$Script:TestResults = [System.Collections.Generic.List[object]]::new()
$Script:TestsPassed = 0
$Script:TestsFailed = 0

function Write-TestResult {
    <#
    .SYNOPSIS
        Records and displays test results with proper formatting.
    
    .DESCRIPTION
        Maintains test statistics and provides colored console output for test results.
    
    .PARAMETER TestName
        The name of the test being reported.
    
    .PARAMETER Passed
        Boolean indicating if the test passed.
    
    .PARAMETER Message
        Optional message with additional test details.
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$TestName,
        
        [Parameter(Mandatory = $true)]
        [bool]$Passed,
        
        [Parameter(Mandatory = $false)]
        [string]$Message = ""
    )
    
    $result = [PSCustomObject]@{
        TestName  = $TestName
        Passed    = $Passed
        Message   = $Message
        Timestamp = Get-Date
    }
    
    $Script:TestResults.Add($result)
    
    if ($Passed) {
        Write-Host "✓ PASS: $TestName" -ForegroundColor Green
        $Script:TestsPassed++
    }
    else {
        Write-Host "✗ FAIL: $TestName - $Message" -ForegroundColor Red
        $Script:TestsFailed++
    }
}

function Test-ConfigurationValidation {
    <#
    .SYNOPSIS
        Tests the configuration validation functionality of PortMap.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing Configuration Validation ===" -ForegroundColor Yellow
    
    try {
        # Test valid configuration
        $null = & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "JSON" -Validate -ErrorAction Stop
        Write-TestResult "Valid Configuration Validation" $true
    }
    catch {
        Write-TestResult "Valid Configuration Validation" $false $_.Exception.Message
    }
    
    # Test invalid configuration (create temporary invalid file)
    $invalidConfig = @{
        devices = @(
            @{
                # Missing required fields
                deviceName = "TEST"
            }
        )
    }
    
    $tempFile = "temp-invalid-config.json"
    $invalidConfig | ConvertTo-Json -Depth 10 | Out-File -FilePath $tempFile
    
    try {
        $null = & ".\PortMap.ps1" -InputFile $tempFile -OutputFormat "JSON" -Validate -ErrorAction Stop 2>&1
        Write-TestResult "Invalid Configuration Detection" $false "Should have failed validation"
    }
    catch {
        Write-TestResult "Invalid Configuration Detection" $true
    }
    finally {
        if (Test-Path $tempFile) {
            Remove-Item $tempFile -Force
        }
    }
}

function Test-MarkdownOutputGeneration {
    <#
    .SYNOPSIS
        Tests the Markdown output generation functionality.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing Markdown Output ===" -ForegroundColor Yellow
    
    try {
        $output = & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "Markdown" -ErrorAction Stop
        
        # Validate markdown structure
        $hasDeviceSummary = $output -match "## Device Summary"
        $hasConnectionMapping = $output -match "## Connection Mapping"
        $hasTables = $output -match "\|.*\|.*\|"
        
        Write-TestResult "Markdown Generation" ($output.Length -gt 0)
        Write-TestResult "Markdown Device Summary Section" $hasDeviceSummary
        Write-TestResult "Markdown Connection Mapping Section" $hasConnectionMapping
        Write-TestResult "Markdown Table Format" $hasTables
    }
    catch {
        Write-TestResult "Markdown Generation" $false $_.Exception.Message
    }
}

function Test-CsvOutputGeneration {
    <#
    .SYNOPSIS
        Tests the CSV output generation functionality.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing CSV Output ===" -ForegroundColor Yellow
    
    $csvFile = "test-output.csv"
    
    try {
        & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "CSV" -OutputFile $csvFile -ErrorAction Stop
        
        if (Test-Path $csvFile) {
            $csvContent = Import-Csv $csvFile
            $hasDeviceRows = ($csvContent | Where-Object { $_.Type -eq "Device" }).Count -gt 0
            $hasConnectionRows = ($csvContent | Where-Object { $_.Type -eq "Connection" }).Count -gt 0
            
            Write-TestResult "CSV File Generation" $true
            Write-TestResult "CSV Device Rows" $hasDeviceRows
            Write-TestResult "CSV Connection Rows" $hasConnectionRows
        }
        else {
            Write-TestResult "CSV File Generation" $false "Output file not created"
        }
    }
    catch {
        Write-TestResult "CSV Generation" $false $_.Exception.Message
    }
    finally {
        if (Test-Path $csvFile) {
            Remove-Item $csvFile -Force
        }
    }
}

function Test-JsonOutputGeneration {
    <#
    .SYNOPSIS
        Tests the JSON output generation functionality.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing JSON Output ===" -ForegroundColor Yellow
    
    $jsonFile = "test-output.json"
    
    try {
        & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "JSON" -OutputFile $jsonFile -ErrorAction Stop
        
        if (Test-Path $jsonFile) {
            $jsonContent = Get-Content $jsonFile -Raw | ConvertFrom-Json
            
            $hasMetadata = $null -ne $jsonContent.metadata
            $hasDevices = $null -ne $jsonContent.devices -and $jsonContent.devices.Count -gt 0
            $hasConnections = $null -ne $jsonContent.connections
            $hasSummary = $null -ne $jsonContent.summary
            
            Write-TestResult "JSON File Generation" $true
            Write-TestResult "JSON Metadata Section" $hasMetadata
            Write-TestResult "JSON Devices Section" $hasDevices
            Write-TestResult "JSON Connections Section" $hasConnections
            Write-TestResult "JSON Summary Section" $hasSummary
        }
        else {
            Write-TestResult "JSON File Generation" $false "Output file not created"
        }
    }
    catch {
        Write-TestResult "JSON Generation" $false $_.Exception.Message
    }
    finally {
        if (Test-Path $jsonFile) {
            Remove-Item $jsonFile -Force
        }
    }
}

function Test-UnusedPortsFeature {
    <#
    .SYNOPSIS
        Tests the unused ports reporting feature.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing Show Unused Ports ===" -ForegroundColor Yellow
    
    try {
        $output = & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "Markdown" -ShowUnused -ErrorAction Stop
        $hasUnusedSection = $output -match "## Unused Ports"
        
        Write-TestResult "Show Unused Ports Feature" $hasUnusedSection
    }
    catch {
        Write-TestResult "Show Unused Ports Feature" $false $_.Exception.Message
    }
}

function Test-DeviceFilteringFeature {
    <#
    .SYNOPSIS
        Tests the device filtering functionality.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing Device Filter ===" -ForegroundColor Yellow
    
    try {
        $output = & ".\PortMap.ps1" -InputFile $TestDataPath -OutputFormat "JSON" -DeviceFilter @("TOR-1") -ErrorAction Stop
        $jsonOutput = $output | ConvertFrom-Json
        
        $filteredCorrectly = $jsonOutput.devices.Count -eq 1 -and $jsonOutput.devices[0].deviceName -eq "TOR-1"
        
        Write-TestResult "Device Filter Functionality" $filteredCorrectly
    }
    catch {
        Write-TestResult "Device Filter Functionality" $false $_.Exception.Message
    }
}

function Show-TestSummary {
    <#
    .SYNOPSIS
        Displays a comprehensive summary of all test results.
    
    .DESCRIPTION
        Shows test statistics, failed tests, and overall success rate with color coding.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Test Summary ===" -ForegroundColor Cyan
    Write-Host "Total Tests: $($Script:TestsPassed + $Script:TestsFailed)" -ForegroundColor White
    Write-Host "Passed: $($Script:TestsPassed)" -ForegroundColor Green
    Write-Host "Failed: $($Script:TestsFailed)" -ForegroundColor Red
    
    if ($Script:TestsFailed -gt 0) {
        Write-Host "`nFailed Tests:" -ForegroundColor Red
        $failedTests = $Script:TestResults | Where-Object { -not $_.Passed }
        foreach ($test in $failedTests) {
            Write-Host "  - $($test.TestName): $($test.Message)" -ForegroundColor Red
        }
    }
    
    $totalTests = $Script:TestsPassed + $Script:TestsFailed
    if ($totalTests -gt 0) {
        $successRate = [math]::Round(($Script:TestsPassed / $totalTests) * 100, 2)
        $successColor = if ($successRate -eq 100) { "Green" } elseif ($successRate -ge 80) { "Yellow" } else { "Red" }
        Write-Host "`nSuccess Rate: $successRate%" -ForegroundColor $successColor
    }
}

# Main execution function
function Start-PortMapTests {
    <#
    .SYNOPSIS
        Main entry point for the PortMap test suite.
    
    .DESCRIPTION
        Executes all tests and returns appropriate exit codes for CI/CD integration.
    #>
    Write-Host "PortMap Tool Test Suite" -ForegroundColor Cyan
    Write-Host "======================" -ForegroundColor Cyan
    Write-Host "Test Data: $TestDataPath" -ForegroundColor White
    
    # Verify test dependencies
    if (-not (Test-Path ".\PortMap.ps1")) {
        Write-Host "ERROR: PortMap.ps1 not found in current directory" -ForegroundColor Red
        exit 1
    }
    
    if (-not (Test-Path $TestDataPath)) {
        Write-Host "ERROR: Test data file not found: $TestDataPath" -ForegroundColor Red
        exit 1
    }
    
    # Run tests
    Test-ConfigurationValidation
    Test-MarkdownOutputGeneration
    Test-CsvOutputGeneration
    Test-JsonOutputGeneration
    Test-UnusedPortsFeature
    Test-DeviceFilteringFeature
    
    # Show results
    Show-TestSummary
    
    # Exit with appropriate code
    if ($Script:TestsFailed -gt 0) {
        exit 1
    }
    else {
        exit 0
    }
}

# Execute tests
Start-PortMapTests
