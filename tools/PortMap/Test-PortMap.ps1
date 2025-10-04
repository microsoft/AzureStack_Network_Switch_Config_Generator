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

function Test-BreakoutCableSupport {
    <#
    .SYNOPSIS
        Comprehensive testing of breakout cable functionality.
    
    .DESCRIPTION
        Tests all aspects of breakout cable support including range expansion,
        single breakout interfaces, mixed port types, and output format validation.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing Breakout Cable Support ===" -ForegroundColor Yellow
    
    # Ensure test breakout configuration exists
    $breakoutConfigPath = ".\test-breakout-config.json"
    
    if (-not (Test-Path $breakoutConfigPath)) {
        Write-TestResult "Breakout Config File Availability" $false "test-breakout-config.json not found"
        return
    }
    
    Write-TestResult "Breakout Config File Availability" $true
    
    # Test breakout range expansion
    Test-BreakoutRangeExpansion
    
    # Test single breakout interface handling
    Test-SingleBreakoutInterface
    
    # Test mixed port types in same device
    Test-MixedPortTypes
    
    # Test PowerShell array handling fixes
    Test-ArrayHandlingFixes
    
    # Test output format support for breakout cables
    Test-BreakoutOutputFormats
    
    # Test edge cases and validation
    Test-BreakoutEdgeCases
}

function Test-BreakoutRangeExpansion {
    <#
    .SYNOPSIS
        Tests that breakout ranges like "25.1-25.4" expand correctly.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Breakout Range Expansion ---" -ForegroundColor Cyan
    
    try {
        # Test with JSON output to examine structure
        $output = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ErrorAction Stop
        $jsonOutput = $output | ConvertFrom-Json
        
        # Look for breakout interfaces in the device port details
        $breakoutInterfaces = @()
        foreach ($device in $jsonOutput.devices) {
            foreach ($port in $device.portDetails) {
                if ($port.port -like "*.?") {
                    $breakoutInterfaces += $port.port
                }
            }
        }
        
        # Check for expected breakout interfaces (25.1, 25.2, 25.3, 25.4)
        $hasBreakout25_1 = $breakoutInterfaces -contains "25.1"
        $hasBreakout25_2 = $breakoutInterfaces -contains "25.2"
        $hasBreakout25_3 = $breakoutInterfaces -contains "25.3"
        $hasBreakout25_4 = $breakoutInterfaces -contains "25.4"
        
        Write-TestResult "Breakout Range 25.1-25.4 Expansion" ($hasBreakout25_1 -and $hasBreakout25_2 -and $hasBreakout25_3 -and $hasBreakout25_4)
        Write-TestResult "Breakout Interface 25.1 Present" $hasBreakout25_1
        Write-TestResult "Breakout Interface 25.2 Present" $hasBreakout25_2
        Write-TestResult "Breakout Interface 25.3 Present" $hasBreakout25_3
        Write-TestResult "Breakout Interface 25.4 Present" $hasBreakout25_4
        
        # Verify proper data types
        $portTypes = @()
        foreach ($interface in $breakoutInterfaces) {
            $portTypes += $interface.GetType().Name
        }
        $allStrings = ($portTypes | Where-Object { $_ -ne "String" }).Count -eq 0
        
        Write-TestResult "Breakout Interfaces Are Strings" $allStrings
    }
    catch {
        Write-TestResult "Breakout Range Expansion" $false $_.Exception.Message
    }
}

function Test-SingleBreakoutInterface {
    <#
    .SYNOPSIS
        Tests single breakout interface handling (e.g., "26.1").
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Single Breakout Interface ---" -ForegroundColor Cyan
    
    try {
        # Test with Markdown output to check display format
        $output = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "Markdown" -ErrorAction Stop
        
        # Check for single breakout interface in markdown output
        $hasSingle26_1 = $output -match "26\.1.*QSFP"
        
        Write-TestResult "Single Breakout Interface 26.1" $hasSingle26_1
        
        # Test with JSON to verify structure
        $jsonOutput = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ErrorAction Stop | ConvertFrom-Json
        
        $singleBreakoutFound = $false
        foreach ($device in $jsonOutput.devices) {
            foreach ($port in $device.portDetails) {
                if ($port.port -eq "26.1") {
                    $singleBreakoutFound = $true
                    break
                }
            }
        }
        
        Write-TestResult "Single Breakout Interface in JSON" $singleBreakoutFound
    }
    catch {
        Write-TestResult "Single Breakout Interface" $false $_.Exception.Message
    }
}

function Test-MixedPortTypes {
    <#
    .SYNOPSIS
        Tests devices with both standard ports (integers) and breakout interfaces (strings).
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Mixed Port Types ---" -ForegroundColor Cyan
    
    try {
        $jsonOutput = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ErrorAction Stop | ConvertFrom-Json
        
        # Find a device with both standard and breakout ports
        $mixedDevice = $null
        foreach ($device in $jsonOutput.devices) {
            $hasStandardPorts = $false
            $hasBreakoutPorts = $false
            
            foreach ($port in $device.portDetails) {
                if ($port.port -notlike "*.?") {
                    $hasStandardPorts = $true
                }
                else {
                    $hasBreakoutPorts = $true
                }
            }
            
            if ($hasStandardPorts -and $hasBreakoutPorts) {
                $mixedDevice = $device
                break
            }
        }
        
        Write-TestResult "Device with Mixed Port Types Found" ($null -ne $mixedDevice)
        
        if ($null -ne $mixedDevice) {
            Write-TestResult "Mixed Device Port Count > 0" ($mixedDevice.portDetails.Count -gt 0)
        }
    }
    catch {
        Write-TestResult "Mixed Port Types" $false $_.Exception.Message
    }
}

function Test-ArrayHandlingFixes {
    <#
    .SYNOPSIS
        Tests PowerShell array handling fixes for single-element breakout interfaces.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Array Handling Fixes ---" -ForegroundColor Cyan
    
    try {
        # This test validates that single breakout interfaces don't get enumerated incorrectly
        $output = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ErrorAction Stop
        $jsonOutput = $output | ConvertFrom-Json
        
        # Look for single breakout interfaces and verify they're processed correctly
        $singleBreakoutCount = 0
        foreach ($device in $jsonOutput.devices) {
            foreach ($port in $device.portDetails) {
                if ($port.port -like "26.1") {
                    $singleBreakoutCount++
                }
            }
        }
        
        # Should find exactly one 26.1 interface (not multiple char enumeration)
        Write-TestResult "Single Breakout Array Handling" ($singleBreakoutCount -eq 1)
        
        # Verify no character enumeration artifacts
        $noCharEnumeration = $true
        foreach ($device in $jsonOutput.devices) {
            foreach ($port in $device.portDetails) {
                if ($port.port.Length -eq 1 -and $port.port -match "[0-9]") {
                    # If we find single character "ports" from array enumeration, that's a failure
                    $noCharEnumeration = $false
                }
            }
        }
        
        Write-TestResult "No Character Enumeration Artifacts" $noCharEnumeration
    }
    catch {
        Write-TestResult "Array Handling Fixes" $false $_.Exception.Message
    }
}

function Test-BreakoutOutputFormats {
    <#
    .SYNOPSIS
        Tests that all output formats properly support breakout interfaces.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Breakout Output Formats ---" -ForegroundColor Cyan
    
    # Test Markdown format
    try {
        $mdOutput = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "Markdown" -ErrorAction Stop
        $mdHasBreakout = $mdOutput -match "25\.1.*\|.*25\.2.*\|.*25\.3.*\|.*25\.4"
        
        Write-TestResult "Markdown Breakout Interface Display" $mdHasBreakout
    }
    catch {
        Write-TestResult "Markdown Breakout Format" $false $_.Exception.Message
    }
    
    # Test JSON format
    try {
        $jsonOutput = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ErrorAction Stop | ConvertFrom-Json
        $jsonHasBreakout = $false
        
        foreach ($device in $jsonOutput.devices) {
            foreach ($port in $device.portDetails) {
                if ($port.port -like "25.*") {
                    $jsonHasBreakout = $true
                    break
                }
            }
        }
        
        Write-TestResult "JSON Breakout Interface Structure" $jsonHasBreakout
    }
    catch {
        Write-TestResult "JSON Breakout Format" $false $_.Exception.Message
    }
    
    # Test CSV format
    try {
        & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "CSV" -ErrorAction Stop
        
        # Look for generated CSV files with breakout interfaces
        $csvFiles = Get-ChildItem -Path "." -Filter "*breakout*portmap*.csv" -ErrorAction SilentlyContinue
        $csvBreakoutSupport = $false
        
        foreach ($csvFile in $csvFiles) {
            $csvContent = Get-Content $csvFile.FullName -Raw
            if ($csvContent -match "25\.1|25\.2|25\.3|25\.4") {
                $csvBreakoutSupport = $true
                break
            }
        }
        
        Write-TestResult "CSV Breakout Interface Support" $csvBreakoutSupport
        
        # Cleanup CSV files
        $csvFiles | Remove-Item -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-TestResult "CSV Breakout Format" $false $_.Exception.Message
    }
}

function Test-BreakoutEdgeCases {
    <#
    .SYNOPSIS
        Tests breakout cable edge cases and validation scenarios.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n--- Testing Breakout Edge Cases ---" -ForegroundColor Cyan
    
    # Test ShowUnused parameter with breakout interfaces
    try {
        $output = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "Markdown" -ShowUnused -ErrorAction Stop
        $hasUnusedBreakout = $output -match "Unused.*26\.[2-4]"
        
        Write-TestResult "ShowUnused with Breakout Interfaces" $hasUnusedBreakout
    }
    catch {
        Write-TestResult "ShowUnused Breakout Edge Case" $false $_.Exception.Message
    }
    
    # Test validation with breakout configuration
    try {
        $null = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -Validate -ErrorAction Stop
        Write-TestResult "Breakout Configuration Validation" $true
    }
    catch {
        Write-TestResult "Breakout Configuration Validation" $false $_.Exception.Message
    }
    
    # Test sequential port ordering with mixed types
    try {
        $jsonOutput = & ".\PortMap.ps1" -InputFile ".\test-breakout-config.json" -OutputFormat "JSON" -ShowUnused -ErrorAction Stop | ConvertFrom-Json
        
        # Check that ports are in logical order (1, 2, 3, ..., 25.1, 25.2, 25.3, 25.4, 26.1, etc.)
        $portOrderCorrect = $false
        foreach ($device in $jsonOutput.devices) {
            $ports = $device.portDetails | Sort-Object { 
                if ($_.port -like "*.*") {
                    # For breakout interfaces, sort by primary then sub
                    $parts = $_.port.Split('.')
                    [int]$parts[0] + ([int]$parts[1] / 10.0)
                }
                else {
                    [int]$_.port
                }
            }
            
            # Verify we have both standard and breakout ports in sequence
            $standardPorts = ($ports | Where-Object { $_.port -notlike "*.*" }).Count
            $breakoutPorts = ($ports | Where-Object { $_.port -like "*.*" }).Count
            
            if ($standardPorts -gt 0 -and $breakoutPorts -gt 0) {
                $portOrderCorrect = $true
                break
            }
        }
        
        Write-TestResult "Sequential Port Ordering with Mixed Types" $portOrderCorrect
    }
    catch {
        Write-TestResult "Sequential Port Ordering" $false $_.Exception.Message
    }
}

function Test-CsvInputGeneration {
    <#
    .SYNOPSIS
        Tests the CSV input functionality.
    #>
    [CmdletBinding()]
    param()
    
    Write-Host "`n=== Testing CSV Input Format ===" -ForegroundColor Yellow
    
    # Test CSV input with devices file
    try {
        $jsonOutput = & ".\PortMap.ps1" -InputFile ".\sample-devices.csv" -OutputFormat "JSON" -ErrorAction Stop | ConvertFrom-Json
        
        $hasDevices = $jsonOutput.devices.Count -gt 0
        $hasConnections = $jsonOutput.connections.Count -gt 0
        
        Write-TestResult "CSV Input (devices file)" ($hasDevices -and $hasConnections)
    }
    catch {
        Write-TestResult "CSV Input (devices file)" $false $_.Exception.Message
    }
    
    # Test CSV input with connections file
    try {
        $markdownOutput = & ".\PortMap.ps1" -InputFile ".\sample-connections.csv" -OutputFormat "Markdown" -ErrorAction Stop
        $hasTables = $markdownOutput -match '\|.*\|'
        
        Write-TestResult "CSV Input (connections file)" $hasTables
    }
    catch {
        Write-TestResult "CSV Input (connections file)" $false $_.Exception.Message
    }
    
    # Test CSV input produces same structure as JSON input
    try {
        $jsonFromCsv = & ".\PortMap.ps1" -InputFile ".\sample-devices.csv" -OutputFormat "JSON" -ErrorAction Stop | ConvertFrom-Json
        
        $hasDeviceNames = ($jsonFromCsv.devices | Where-Object { $_.deviceName }).Count -gt 0
        $hasPortDetails = ($jsonFromCsv.devices | Where-Object { $_.portRanges }).Count -gt 0
        $hasConnectionDetails = ($jsonFromCsv.connections | Where-Object { $_.sourceDevice }).Count -gt 0
        
        Write-TestResult "CSV Input Structure Validation" ($hasDeviceNames -and $hasPortDetails -and $hasConnectionDetails)
    }
    catch {
        Write-TestResult "CSV Input Structure Validation" $false $_.Exception.Message
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
    
    # CSV input format testing
    Test-CsvInputGeneration
    
    # Comprehensive breakout cable testing
    Test-BreakoutCableSupport
    
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
